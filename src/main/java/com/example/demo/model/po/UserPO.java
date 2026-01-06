package com.example.demo.model.po;

import lombok.Data;

/**
 * PO (Persistent Object): 数据库里的原始数据
 * 绝对不包含业务逻辑，字段和数据库表一一对应。
 */
@Data
public class UserPO {
    private Long id;
    private String username;
    private String password; // 敏感数据！数据库里有，但不能给前端看
    private Integer age;     // 这是我们新加的字段
    private Integer isDeleted; // 逻辑删除标记，业务层不需要关心，只在SQL里用
}
