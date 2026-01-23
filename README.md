# Maven 命令解释：mvn dependency:tree -Dverbose -Dincludes=

这条命令主要用于 **排查 Maven 项目中特定依赖的冲突和路径来源**。

## 详细参数解析

### 1. `mvn dependency:tree`
这是 Maven 依赖插件的一个 goal。它的作用是以树状结构打印出项目中所有的依赖关系（包括直接依赖和传递性依赖）。

### 2. `-Dverbose`
开启详细模式。
*   **作用**：默认情况下，Maven 的依赖树只显示最终被采用的依赖版本。加上 `-Dverbose` 后，它会显示 **所有** 被解析过的依赖，包括那些因为版本冲突而被 Maven 仲裁机制（"nearest definition" 策略）忽略掉的依赖。
*   **输出特征**：你会看到类似 `(omitted for conflict with 1.2.3)` 的标记，这对解决 `Jar 包冲突` 或 `NoSuchMethodError` 非常关键。

### 3. `-Dincludes=[pattern]`
过滤条件。
*   **作用**：在一个庞大的项目中，打印完整的依赖树会有成百上千行，很难阅读。`-Dincludes` 允许你指定只显示包含特定关键字的依赖路径。
*   **格式**：`groupId:artifactId:type:classifier:version`（支持通配符 `*`）。

---

## 组合使用的场景

当你运行 `mvn dependency:tree -Dverbose -Dincludes=groupId:artifactId` 时，你是在告诉 Maven：

> "请把所有与 `groupId:artifactId` 相关的依赖树都列出来，包括那些被冲突掉的版本，并且只显示这一部分，其他的不要显示。"

### 举例

假设你想查看项目中 `jackson-databind` 的依赖情况，以及是否有版本冲突：

```bash
mvn dependency:tree -Dverbose -Dincludes=com.fasterxml.jackson.core:jackson-databind
```

**可能的输出示例：**

```text
[INFO] com.example:my-project:jar:1.0.0
[INFO] \- org.springframework.boot:spring-boot-starter-web:jar:2.5.0:compile
[INFO]    \- org.springframework.boot:spring-boot-starter-json:jar:2.5.0:compile
[INFO]       \- com.fasterxml.jackson.core:jackson-databind:jar:2.12.3:compile
[INFO]          \- (com.fasterxml.jackson.core:jackson-databind:jar:2.9.8:compile - omitted for conflict with 2.12.3)
```

在这个例子中：
1. 你可以看到 `jackson-databind:2.12.3` 是最终生效的版本。
2. 你还可以看到曾经有一个 `2.9.8` 的版本尝试进入项目，但因为冲突被忽略了（这通常是 `-Dverbose` 带来的额外信息）。
