package com.jtlapp.jvmvsjs.springjdbc.controllers;

import com.google.gson.JsonObject;
import com.jtlapp.jvmvsjs.jdbcquery.SharedQuery;
import com.jtlapp.jvmvsjs.jdbcquery.SharedQueryDB;
import com.jtlapp.jvmvsjs.jdbcquery.SharedQueryException;
import com.jtlapp.jvmvsjs.jdbcquery.SharedQueryRepo;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

@RestController
@RequestMapping("/api")
public class ApiController {

    @Value("${application.name")
    public String appName;
    @Value("${application.version")
    public String appVersion;

    @Autowired
    private SharedQueryDB sharedQueryDB;
    @Autowired
    private SharedQueryRepo sharedQueryRepo;

    @GetMapping("/info")
    public ResponseEntity<String> info() {
        var gson = new JsonObject();
        gson.addProperty("appName", appName);
        gson.addProperty("appVersion", appVersion);
        gson.add("appConfig", new JsonObject());
        return ResponseEntity.ok(gson.toString());
    }

    @PostMapping("/query/{queryName}")
    public ResponseEntity<String> query(
            @PathVariable(name = "queryName") String queryName,
            @RequestBody String jsonBody
    ) {
        try {
            SharedQuery query = sharedQueryRepo.get(sharedQueryDB, queryName);
            String jsonResponse = query.executeUsingGson(sharedQueryDB, jsonBody);
            return ResponseEntity.ok(jsonResponse);
        }
        catch (SharedQueryException e) {
            return ResponseEntity
                    .status(HttpStatus.BAD_REQUEST)
                    .body(toErrorJson(queryName, e));
        }
    }

    @GetMapping("/sleep/{millis}")
    public ResponseEntity<Void> sleep(@PathVariable(name = "millis") int millis) {
        try {
            Thread.sleep(millis);
            return ResponseEntity.ok().build();
        } catch (InterruptedException e) {
            return ResponseEntity
                    .status(HttpStatus.INTERNAL_SERVER_ERROR)
                    .build();
        }
    }

    private String toErrorJson(String queryName, Throwable e) {
        return String.format("{\"query\": \"%s\", \"error\": \"%s: %s\"}",
                queryName, e.getClass().getSimpleName(), e.getMessage());
    }
}
