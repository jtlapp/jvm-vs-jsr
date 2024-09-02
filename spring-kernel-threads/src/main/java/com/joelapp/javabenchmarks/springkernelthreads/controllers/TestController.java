package com.joelapp.javabenchmarks.springkernelthreads.controllers;

import com.joelapp.javabenchmarks.hikariquerylibrary.HikariQueries;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

@RestController
@RequestMapping("/api")
public class TestController {

    @Autowired
    private HikariQueries queries;

    @GetMapping("/setup")
    public ResponseEntity<String> setup() {
        try {
            queries.createTables();
            queries.populateDatabase();
            return ResponseEntity.ok("Completed setup.");
        }
        catch (Exception e) {
            return ResponseEntity.status(HttpStatus.INTERNAL_SERVER_ERROR).body(e.getMessage());
        }
    }

    @GetMapping("/read")
    public ResponseEntity<String> read(
            @RequestParam("user") Integer userNumber,
            @RequestParam("order") Integer orderNumber
    ) {
        try {
            String json = queries.readQuery(userNumber, orderNumber);
            return ResponseEntity.ok(json);
        }
        catch (Exception e) {
            return ResponseEntity.status(HttpStatus.INTERNAL_SERVER_ERROR).body(e.getMessage());
        }
    }
}
