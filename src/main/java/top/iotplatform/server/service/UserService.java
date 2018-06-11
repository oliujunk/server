package top.iotplatform.server.service;

import top.iotplatform.server.bean.User;

import java.util.List;

public interface UserService {
    public User getUserById(int id);
    public User getUserByName(String name);
    public List<User> getUserList();
    public int add(User user);
    public int update(int id, User user);
    public int delete(int id);
    public boolean isUsernameExist(String name);
}
