package top.iotplatform.server.controller;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;
import top.iotplatform.server.bean.JsonResult;
import top.iotplatform.server.bean.User;
import top.iotplatform.server.service.UserService;

import java.util.List;

import static org.springframework.web.bind.annotation.RequestMethod.*;

@RequestMapping("/api")
@RestController
public class UserController {
    @Autowired
    private UserService userService;

    @GetMapping("/users/{id}")
    public ResponseEntity<JsonResult> getUserById (@PathVariable(value = "id") int id) {
        JsonResult r = new JsonResult();
        try {
            User user = userService.getUserById(id);
            r.setResult(user);
            r.setStatus("ok");
        } catch (Exception e) {
            r.setResult(e.getClass().getName() + ":" + e.getMessage());
            r.setStatus("error");
        }
        return ResponseEntity.ok(r);
    }

    @GetMapping("/users")
    public ResponseEntity<JsonResult> getUserList () {
        JsonResult r = new JsonResult();
        try {
            List<User> userList = userService.getUserList();
            r.setResult(userList);
            r.setStatus("ok");
        } catch (Exception e) {
            r.setResult(e.getClass().getName() + ":" + e.getMessage());
            r.setStatus("error");
        }
        return ResponseEntity.ok(r);
    }

    @PostMapping("/users")
    public ResponseEntity<JsonResult> add (@RequestBody User user) {
        JsonResult r = new JsonResult();
        try {
            if (userService.isUsernameExist(user.getUserName())) {
                r.setResult("用户名已被占用，请尝试其他用户名!");
                r.setStatus("fail");
            } else {
                int orderId = userService.add(user);
                if (orderId < 0) {
                    r.setResult(orderId);
                    r.setStatus("fail");
                } else {
                    r.setResult(orderId);
                    r.setStatus("ok");
                }
            }
        } catch (Exception e) {
            r.setResult(e.getClass().getName() + ":" + e.getMessage());
            r.setStatus("error");
        }
        return ResponseEntity.ok(r);
    }
    @PutMapping("/users/{id}")
    public ResponseEntity<JsonResult> update (@PathVariable("id") int id, @RequestBody User user) {
        JsonResult r = new JsonResult();
        try {
            int ret = userService.update(id, user);
            if (ret < 0) {
                r.setResult(ret);
                r.setStatus("fail");
            } else {
                r.setResult(ret);
                r.setStatus("ok");
            }
        } catch (Exception e) {
            r.setResult(e.getClass().getName() + ":" + e.getMessage());
            r.setStatus("error");
        }
        return ResponseEntity.ok(r);
    }

    @DeleteMapping("/users/{id}")
    public ResponseEntity<JsonResult> delete (@PathVariable("id") int id) {
        JsonResult r = new JsonResult();
        try {
            int ret = userService.delete(id);
            if (ret < 0) {
                r.setResult(ret);
                r.setStatus("fail");
            } else {
                r.setResult(ret);
                r.setStatus("ok");
            }
        } catch (Exception e) {
            r.setResult(e.getClass().getName() + ":" + e.getMessage());
            r.setStatus("error");
        }
        return ResponseEntity.ok(r);
    }

}
