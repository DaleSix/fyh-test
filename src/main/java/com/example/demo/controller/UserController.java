package com.example.demo.controller;

import com.example.demo.model.dto.UserDTO;
import com.example.demo.service.UserService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

@RestController
@RequestMapping("/users")
public class UserController {

    @Autowired
    private UserService userService;

    // 前端发请求：GET /users/1
    @GetMapping("/{id}")
    public UserDTO getUser(@PathVariable Long id) {
        // 服务员只管把单子给后厨，然后把做好的菜(DTO)端给客人
        return userService.getUserDetail(id);
    }
}
