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
        String sql = "select * from users where id = ?";
        List<User> list = jdbcTemplate.query(sql, new Object[]{id}, new BeanPropertyRowMapper(User.class));
        if (list != null && list.size() > 0) {
            return list.get(0);
        } else {
            return null;
        }
    }

    @Override
    public User getUserByName(String name) {
        String sql = "select * from users where user_name = ?";
        List<User> list = jdbcTemplate.query(sql, new Object[]{name}, new BeanPropertyRowMapper(User.class));
        if (list != null && list.size() > 0) {
            return list.get(0);
        } else {
            return null;
        }
    }

    @Override
    public List<User> getUserList() {
        String sql = "select * from users";
        List<User> list = jdbcTemplate.query(sql, new BeanPropertyRowMapper(User.class));
        return list;
    }

    @Override
    public int add(User user) {
        return 0;
    }

    @Override
    public int update(User user) {
        return 0;
    }

    @Override
    public int delete(int id) {
        return 0;
    }
}
