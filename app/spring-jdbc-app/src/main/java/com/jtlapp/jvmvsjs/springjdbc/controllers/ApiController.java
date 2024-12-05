package com.jtlapp.jvmvsjs.springjdbc.controllers;

import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.jtlapp.jvmvsjs.jdbclib.Database;
import com.jtlapp.jvmvsjs.springjdbc.config.AppConfig;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.sql.SQLException;

@RestController
@RequestMapping("/api")
public class ApiController {

    static final String appName = System.getenv("APP_NAME");
    static final String appVersion = System.getenv("APP_VERSION");
    static final ObjectMapper objectMapper = new ObjectMapper();

    @Autowired
    public AppConfig appConfig;

    @Autowired
    private Database db;

    @GetMapping("/info")
    public ResponseEntity<String> info() {
        try {
            var jsonObj = objectMapper.createObjectNode()
                    .put("appName", appName)
                    .put("appVersion", appVersion)
                    .set("appConfig", appConfig.toJsonNode(objectMapper));
            return ResponseEntity.ok(objectMapper.writeValueAsString(jsonObj));
        } catch (JsonProcessingException e) {
            return ResponseEntity
                    .status(HttpStatus.INTERNAL_SERVER_ERROR)
                    .body(toErrorJson("info", e));
        }
    }

    @GetMapping("/app-sleep")
    public ResponseEntity<Void> appSleep(@RequestParam("millis") int millis) {
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
    public ResponseEntity<String> pgSleep(@RequestParam("millis") int millis) {
        try(var conn = db.openConnection()) {
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
