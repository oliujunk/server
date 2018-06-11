package top.iotplatform.server.dao;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.jdbc.core.BeanPropertyRowMapper;
import org.springframework.jdbc.core.JdbcTemplate;
import org.springframework.stereotype.Repository;
import top.iotplatform.server.bean.User;

import java.util.List;

@Repository
public class UserDAOImpl implements UserDAO {

    @Autowired
    private JdbcTemplate jdbcTemplate;

    @Override
    public User getUserById(int id) {
        String sql = "SELECT * FROM users WHERE id = ?";
        List<User> list = jdbcTemplate.query(sql, new Object[]{id}, new BeanPropertyRowMapper(User.class));
        if (list != null && list.size() > 0) {
            return list.get(0);
        } else {
            return null;
        }
    }

    @Override
    public User getUserByName(String name) {
        String sql = "SELECT * FROM users WHERE user_name = ?";
        List<User> list = jdbcTemplate.query(sql, new Object[]{name}, new BeanPropertyRowMapper(User.class));
        if (list != null && list.size() > 0) {
            return list.get(0);
        } else {
            return null;
        }
    }

    @Override
    public List<User> getUserList() {
        String sql = "SELECT * FROM users";
        List<User> list = jdbcTemplate.query(sql, new BeanPropertyRowMapper(User.class));
        return list;
    }

    @Override
    public int add(User user) {
        String sql = "INSERT INTO users (user_name, user_passwd, user_group) VALUES (?, ?, ?)";
        return jdbcTemplate.update(sql, user.getUserName(), user.getUserPasswd(), user.getUserGroup());
    }

    @Override
    public int update(int id, User user) {
        String sql = "UPDATE users SET user_name=?, user_passwd=?, user_group=? WHERE id=?";
        return jdbcTemplate.update(sql, user.getUserName(), user.getUserPasswd(), user.getUserGroup(), id);
    }

    @Override
    public int delete(int id) {
        String sql = "DELETE FROM users WHERE id=?";
        return jdbcTemplate.update(sql, id);
    }

    @Override
    public boolean isUsernameExist(String name) {
        String sql = "SELECT * FROM users WHERE user_name=?";
        List<User> list = jdbcTemplate.query(sql, new Object[]{name}, new BeanPropertyRowMapper(User.class));
        if (list != null && list.size() > 0) {
            return true;
        } else {
            return false;
        }
    }
}
