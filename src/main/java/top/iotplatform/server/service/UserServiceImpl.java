package top.iotplatform.server.service;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;
import top.iotplatform.server.bean.User;
import top.iotplatform.server.dao.UserDAO;

import java.util.List;

@Service
public class UserServiceImpl implements UserService {

    @Autowired
    private UserDAO userDAO;

    @Override
    public User getUserById(int id) {
        return userDAO.getUserById(id);
    }

    @Override
    public User getUserByName(String name) {
        return userDAO.getUserByName(name);
    }

    @Override
    public List<User> getUserList() {
        return userDAO.getUserList();
    }

    @Override
    public int add(User user) {
        return userDAO.add(user);
    }

    @Override
    public int update(int id, User user) {
        return userDAO.update(id, user);
    }

    @Override
    public int delete(int id) {
        return userDAO.delete(id);
    }

    @Override
    public boolean isUsernameExist(String name) {
        return userDAO.isUsernameExist(name);
    }

}
