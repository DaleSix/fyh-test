package main

import (
	"errors"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/xuri/excelize/v2"
)

type Sample struct {
	T time.Time
	V float64
}

type Series struct {
	Name    string
	Samples []Sample
	Values  []float64
	ByTime  map[time.Time]float64
}

func main() {
	if len(os.Args) < 3 {
		os.Args = append(os.Args, "solar_bw.xlsx", "solar_ab_bw_aggregation.xlsx")
	}

	s1, err := readExcelSeries(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "读取失败 %s: %v\n", os.Args[1], err)
		os.Exit(1)
	}
	s2, err := readExcelSeries(os.Args[2])
	if err != nil {
		fmt.Fprintf(os.Stderr, "读取失败 %s: %v\n", os.Args[2], err)
		os.Exit(1)
	}

	p1, t1, ok1 := percentile95(s1.Samples)
	p2, t2, ok2 := percentile95(s2.Samples)

	if ok1 {
		fmt.Printf("%s 95带宽点: %.6f (示例时间: %s)\n", s1.Name, p1, formatTime(t1))
	} else {
		fmt.Printf("%s 无可用数据，无法计算95带宽点\n", s1.Name)
	}
	if ok2 {
		fmt.Printf("%s 95带宽点: %.6f (示例时间: %s)\n", s2.Name, p2, formatTime(t2))
	} else {
		fmt.Printf("%s 无可用数据，无法计算95带宽点\n", s2.Name)
	}

	unionTimes := unionSortedTimes(s1.ByTime, s2.ByTime)
	x := make([]string, 0, len(unionTimes))
	for _, tt := range unionTimes {
		x = append(x, formatTime(tt))
	}

	line := charts.NewLine()
	line.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{
			PageTitle: "Bandwidth",
			Width:     "1200px",
			Height:    "600px",
		}),
		charts.WithTitleOpts(opts.Title{
			Title:    "Bandwidth Curve",
			Subtitle: "From two Excel files (Time + Bandwidth)",
		}),
		charts.WithTooltipOpts(opts.Tooltip{
			Show:        opts.Bool(true),
			Trigger:     "axis",
			AxisPointer: &opts.AxisPointer{Type: "cross"},
		}),
		charts.WithLegendOpts(opts.Legend{
			Show: opts.Bool(true),
		}),
		charts.WithXAxisOpts(opts.XAxis{
			Type: "category",
		}),
		charts.WithYAxisOpts(opts.YAxis{
			Name: "Bandwidth",
			Type: "value",
		}),
	)

	line.SetXAxis(x)

	s1SeriesOpts := []charts.SeriesOpts{
		charts.WithLineChartOpts(opts.LineChart{Smooth: opts.Bool(true)}),
	}
	if ok1 {
		s1SeriesOpts = append(
			s1SeriesOpts,
			charts.WithMarkLineStyleOpts(opts.MarkLineStyle{
				Label: &opts.Label{Show: opts.Bool(true), Formatter: "{b}: {c}"},
				LineStyle: &opts.LineStyle{Type: "dashed"},
			}),
			charts.WithMarkLineNameYAxisItemOpts(opts.MarkLineNameYAxisItem{
				Name:  fmt.Sprintf("%s 95计费点", s1.Name),
				YAxis: p1,
			}),
		)
	}
	line.AddSeries(s1.Name, alignSeries(unionTimes, s1.ByTime), s1SeriesOpts...)

	s2SeriesOpts := []charts.SeriesOpts{
		charts.WithLineChartOpts(opts.LineChart{Smooth: opts.Bool(true)}),
	}
	if ok2 {
		s2SeriesOpts = append(
			s2SeriesOpts,
			charts.WithMarkLineStyleOpts(opts.MarkLineStyle{
				Label: &opts.Label{Show: opts.Bool(true), Formatter: "{b}: {c}"},
				LineStyle: &opts.LineStyle{Type: "dashed"},
			}),
			charts.WithMarkLineNameYAxisItemOpts(opts.MarkLineNameYAxisItem{
				Name:  fmt.Sprintf("%s 95计费点", s2.Name),
				YAxis: p2,
			}),
		)
	}
	line.AddSeries(s2.Name, alignSeries(unionTimes, s2.ByTime), s2SeriesOpts...)

	outName := "bandwidth.html"
	f, err := os.Create(outName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "创建输出文件失败: %v\n", err)
		os.Exit(1)
	}
	defer func() {
		_ = f.Close()
	}()

	if err := line.Render(f); err != nil {
		fmt.Fprintf(os.Stderr, "渲染图表失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("已生成图表: %s\n", outName)
}

func readExcelSeries(path string) (Series, error) {
	f, err := excelize.OpenFile(path)
	if err != nil {
		return Series{}, err
	}
	defer func() {
		_ = f.Close()
	}()

	sheets := f.GetSheetList()
	if len(sheets) == 0 {
		return Series{}, errors.New("xlsx没有工作表")
	}
	var sheet string
	for _, sheetName := range sheets {
		rows, err := f.GetRows(sheetName)
		if err != nil {
			continue
		}
		if len(rows) > 1 && len(rows[0]) > 1 {
			sheet = sheetName
			break
		}
	}
	if sheet == "" {
		sheet = sheets[0]
	}

	rows, err := f.GetRows(sheet)
	if err != nil {
		return Series{}, err
	}
	name := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
	out := Series{
		Name:   name,
		ByTime: make(map[time.Time]float64, len(rows)),
	}

	for i, row := range rows {
		if len(row) < 2 {
			fmt.Printf("跳过无效行 %d (列数不足2)\n", i+1)
			continue
		}

		t, errT := parseExcelTime(row[0])
		v, errV := parseFloat(row[1])
		if errT != nil || errV != nil {
			if i == 0 {
				continue
			}
			fmt.Printf("警告: 跳过无效行 %d (时间=%q, 带宽=%q) errT:%v errV: %v\n", i+1, row[0], row[1], errT, errV)
			continue
		}

		out.Samples = append(out.Samples, Sample{T: t, V: v})
		out.Values = append(out.Values, v)
		out.ByTime[t] = v
	}

	if len(out.Samples) == 0 {
		return out, errors.New("未读到有效数据行(需要两列: 时间, 带宽)")
	}

	sort.Slice(out.Samples, func(i, j int) bool {
		return out.Samples[i].T.Before(out.Samples[j].T)
	})

	return out, nil
}

func parseExcelTime(s string) (time.Time, error) {
	ss := strings.TrimSpace(s)
	if ss == "" {
		return time.Time{}, errors.New("empty time")
	}

	loc := time.Local
	layouts := []string{
		"2006/1/2 15:04",
		"2006/01/02 15:04",
		"2006-1-2 15:04",
		"2006-01-02 15:04",
		"2006-01-02 15:04:05",
		time.RFC3339,
		time.RFC3339Nano,
	}

	for _, layout := range layouts {
		if t, err := time.ParseInLocation(layout, ss, loc); err == nil {
			return t, nil
		}
	}

	if fv, err := strconv.ParseFloat(ss, 64); err == nil {
		t, err2 := excelize.ExcelDateToTime(fv, false)
		if err2 != nil {
			return time.Time{}, err2
		}
		return t.In(loc), nil
	}

	return time.Time{}, fmt.Errorf("无法解析时间: %q", s)
}

func parseFloat(s string) (float64, error) {
	ss := strings.TrimSpace(s)
	ss = strings.ReplaceAll(ss, ",", "")
	if ss == "" {
		return 0, errors.New("empty number")
	}
	return strconv.ParseFloat(ss, 64)
}

func formatTime(t time.Time) string {
	return t.Format("2006-01-02 15:04")
}

func unionSortedTimes(m1, m2 map[time.Time]float64) []time.Time {
	set := make(map[time.Time]struct{}, len(m1)+len(m2))
	for t := range m1 {
		set[t] = struct{}{}
	}
	for t := range m2 {
		set[t] = struct{}{}
	}

	out := make([]time.Time, 0, len(set))
	for t := range set {
		out = append(out, t)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].Before(out[j])
	})
	return out
}

func alignSeries(times []time.Time, m map[time.Time]float64) []opts.LineData {
	out := make([]opts.LineData, 0, len(times))
	for _, t := range times {
		if v, ok := m[t]; ok {
			out = append(out, opts.LineData{Value: v})
		} else {
			out = append(out, opts.LineData{Value: nil})
		}
	}
	return out
}

func percentile95(samples []Sample) (p95 float64, sampleTime time.Time, ok bool) {
	n := len(samples)
	if n == 0 {
		return 0, time.Time{}, false
	}

	vals := make([]float64, 0, n)
	for _, s := range samples {
		vals = append(vals, s.V)
	}
	sort.Float64s(vals)

	rank := int(math.Ceil(0.95*float64(n))) - 1
	if rank < 0 {
		rank = 0
	}
	if rank >= n {
		rank = n - 1
	}
	p95 = vals[rank]

	best := samples[0]
	bestDist := math.Abs(best.V - p95)
	for _, s := range samples[1:] {
		d := math.Abs(s.V - p95)
		if d < bestDist {
			best = s
			bestDist = d
		}
	}
	return p95, best.T, true
}

