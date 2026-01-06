package com.example.demo.dao;

import com.example.demo.model.po.UserPO;
import org.apache.ibatis.annotations.Mapper;
import org.apache.ibatis.annotations.Param;

/**
 * Dao / Mapper: 唯一能直接碰数据库的地方
 * 返回的一定是 PO 对象。
 */
@Mapper
public interface UserMapper {
    // 根据ID查用户，返回 PO
    UserPO selectById(@Param("id") Long id);
}
