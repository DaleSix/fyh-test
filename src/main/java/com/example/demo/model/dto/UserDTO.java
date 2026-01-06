package com.example.demo.model.dto;

import lombok.Data;

/**
 * DTO (Data Transfer Object): 给前端/APP看的数据
 * 经过了脱敏、格式化。
 */
@Data
public class UserDTO {
    private String userId;   // 前端可能喜欢叫 userId 而不是 id
    private String name;     // 前端可能喜欢叫 name 而不是 username
    private Integer age;
    private String ageGroup; // "成年人" 或 "未成年"
    
    // 绝对没有 password
    // 绝对没有 isDeleted
}
