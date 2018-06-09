package top.iotplatform.server.dao;

import top.iotplatform.server.bean.User;

import java.util.List;

public interface UserDAO {
    User getUserById(int id);
    User getUserByName(String name);
    public List<User> getUserList();
    public int add(User user);
    public int update(User user);
    public int delete(int id);
}
