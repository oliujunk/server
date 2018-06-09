package top.iotplatform.server.controller;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;
import top.iotplatform.server.bean.JsonResult;
import top.iotplatform.server.bean.User;
import top.iotplatform.server.service.UserService;

import java.util.List;

@RestController
public class UserController {
    @Autowired
    private UserService userService;

    @RequestMapping("/users/{id:[0-9]]}")
    public ResponseEntity<JsonResult> getUserById (@PathVariable(value = "id") int id) {
        JsonResult r = new JsonResult();
        try {
            User user = userService.getUserById(id);
            r.setResult(user);
            r.setStatus("ok");
        } catch (Exception e) {
            r.setResult(e.getClass().getName() + ":" + e.getMessage());
            r.setStatus("error");
            e.printStackTrace();
        }
        return ResponseEntity.ok(r);
    }

    @RequestMapping("/users/{name}")
    public ResponseEntity<JsonResult> getUserByName (@PathVariable(value = "name") String name) {
        JsonResult r = new JsonResult();
        try {
            User user = userService.getUserByName(name);
            r.setResult(user);
            r.setStatus("ok");
        } catch (Exception e) {
            r.setResult(e.getClass().getName() + ":" + e.getMessage());
            r.setStatus("error");
            e.printStackTrace();
        }
        return ResponseEntity.ok(r);
    }

    @RequestMapping("/users")
    public ResponseEntity<JsonResult> getUserList () {
        JsonResult r = new JsonResult();
        try {
            List<User> userList = userService.getUserList();
            r.setResult(userList);
            r.setStatus("ok");
        } catch (Exception e) {
            r.setResult(e.getClass().getName() + ":" + e.getMessage());
            r.setStatus("error");
            e.printStackTrace();
        }
        return ResponseEntity.ok(r);
    }
}
