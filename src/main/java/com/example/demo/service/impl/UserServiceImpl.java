package com.example.demo.service.impl;

import com.example.demo.dao.UserMapper;
import com.example.demo.model.bo.UserBO;
import com.example.demo.model.dto.UserDTO;
import com.example.demo.model.po.UserPO;
import com.example.demo.service.UserService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

@Service
public class UserServiceImpl implements UserService {

    @Autowired
    private UserMapper userMapper; // 注入库房管理员

    @Override
    public UserDTO getUserDetail(Long id) {
        // 1. 调用 Dao 拿到“带泥土豆” (PO)
        UserPO po = userMapper.selectById(id);
        if (po == null) {
            return null;
        }

        // 2. 把 PO 转成 BO (清洗、去皮)
        // 在这一步，我们把数据库的字段转成了业务对象
        UserBO bo = new UserBO();
        bo.setId(po.getId());
        bo.setUsername(po.getUsername());
        bo.setAge(po.getAge());
        // 密码在这里就被丢弃了，不会进入业务层逻辑

        // 3. 执行业务逻辑 (切丝、烹饪)
        // 比如：判断是否成年，这属于业务逻辑
        boolean isAdult = bo.isAdult();

        // 4. 把 BO 转成 DTO (摆盘)
        UserDTO dto = new UserDTO();
        dto.setUserId(String.valueOf(bo.getId())); // ID转字符串，防止前端精度丢失
        dto.setName(bo.getUsername());
        dto.setAge(bo.getAge());
        dto.setAgeGroup(isAdult ? "成年人" : "未成年");

        return dto;
    }
}
