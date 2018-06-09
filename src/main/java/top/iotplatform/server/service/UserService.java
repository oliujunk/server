package top.iotplatform.server.service;

import top.iotplatform.server.bean.User;

import java.util.List;

public interface UserService {
    User getUserById(int id);
    User getUserByName(String name);
    List<User> getUserList();
}
