package com.jtlapp.jvmvsjs.springjdbc.controllers;

import com.google.gson.JsonObject;
import com.jtlapp.jvmvsjs.jdbcquery.Database;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.sql.SQLException;

@RestController
@RequestMapping("/api")
public class ApiController {

    static final String appName = System.getenv("APP_NAME");;
    static final String appVersion = System.getenv("APP_VERSION");;

    @Autowired
    private Database db;

    @GetMapping("/info")
    public ResponseEntity<String> info() {
        var gson = new JsonObject();
        gson.addProperty("appName", appName);
        gson.addProperty("appVersion", appVersion);
        gson.add("appConfig", new JsonObject());
        return ResponseEntity.ok(gson.toString());
    }

    @GetMapping("/app-sleep")
    public ResponseEntity<Void> appSleep(@RequestParam int millis) {
        try {
            Thread.sleep(millis);
            return ResponseEntity.ok().build();
        } catch (InterruptedException e) {
            return ResponseEntity
                    .status(HttpStatus.INTERNAL_SERVER_ERROR)
                    .build();
        }
    }

    @GetMapping("/pg-sleep")
    public ResponseEntity<String> pgSleep(@RequestParam int millis) {
        try {
            var conn = db.openConnection();
            Database.issueSleepQuery(conn, millis);
            return ResponseEntity.ok("{}");
        }
        catch (SQLException e) {
            return ResponseEntity
                    .status(HttpStatus.BAD_REQUEST)
                    .body(toErrorJson("pg-sleep", e));
        }
    }

    private String toErrorJson(String endpoint, Throwable e) {
        return String.format("{\"endpoint\": \"%s\", \"error\": \"%s: %s\"}",
                endpoint, e.getClass().getSimpleName(), e.getMessage());
    }
}
