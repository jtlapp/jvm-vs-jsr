package com.joelapp.javabenchmarks.springkernelthreads.controllers;

import com.joelapp.javabenchmarks.jdbcquerylibrary.JdbcQueries;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

@RestController
@RequestMapping("/api")
public class ApiController {

    @Autowired
    private JdbcQueries queries;

    @GetMapping("/setup")
    public ResponseEntity<String> setup() {
        try {
            queries.dropTables();
            queries.createTables();
            queries.populateDatabase();
            return ResponseEntity.ok("Completed setup");
        }
        catch (Exception e) {
            return ResponseEntity.status(HttpStatus.INTERNAL_SERVER_ERROR)
                    .body(e.getMessage());
        }
    }

    @GetMapping("/select")
    public ResponseEntity<String> select(
            @RequestParam("user") Integer userNumber,
            @RequestParam("order") Integer orderNumber
    ) {
        try {
            String json = queries.selectQuery(userNumber, orderNumber);
            return ResponseEntity.ok(json);
        }
        catch (Exception e) {
            return ResponseEntity.status(HttpStatus.INTERNAL_SERVER_ERROR)
                    .body(e.getMessage());
        }
    }

    @GetMapping("/update")
    public ResponseEntity<String> update(
            @RequestParam("user") Integer userNumber,
            @RequestParam("order") Integer orderNumber
    ) {
        try {
            int rowsUpdated = queries.updateQuery(userNumber, orderNumber);
            return ResponseEntity.ok(rowsUpdated == 0 ? "Nothing updated" : "Updated");
        }
        catch (Exception e) {
            return ResponseEntity.status(HttpStatus.INTERNAL_SERVER_ERROR)
                    .body(e.getMessage());
        }
    }
}
