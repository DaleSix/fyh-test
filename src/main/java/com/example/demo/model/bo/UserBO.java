package com.example.demo.model.bo;

import lombok.Data;

/**
 * BO (Business Object): 业务层处理的数据
 * 可能比 PO 多一些计算属性，或者包含其他 PO 的组合。
 */
@Data
public class UserBO {
    private Long id;
    private String username;
    private Integer age;
    
    // 注意：这里没有 password，因为业务层处理用户信息通常不需要明文密码
    // 注意：这里没有 isDeleted，因为业务层只处理“活着”的用户
    
    // BO 可以包含简单的业务逻辑方法
    public boolean isAdult() {
        return this.age != null && this.age >= 18;
    }
    
    public String getDisplayName() {
        return "User: " + username;
    }
}
