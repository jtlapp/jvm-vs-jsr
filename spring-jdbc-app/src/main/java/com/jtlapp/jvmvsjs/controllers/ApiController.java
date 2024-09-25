package com.jtlapp.jvmvsjs.controllers;

import com.jtlapp.jvmvsjs.jdbcsharedquery.SharedQuery;
import com.jtlapp.jvmvsjs.jdbcsharedquery.SharedQueryDB;
import com.jtlapp.jvmvsjs.jdbcsharedquery.SharedQueryException;
import com.jtlapp.jvmvsjs.jdbcsharedquery.SharedQueryRepo;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

@RestController
@RequestMapping("/api")
public class ApiController {

    @Autowired
    private SharedQueryDB sharedQueryDB;
    @Autowired
    private SharedQueryRepo sharedQueryRepo;

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
            String jsonResponse = String.format("{\"error\": \"%s\"}", e.getMessage());
            return ResponseEntity
                    .status(HttpStatus.BAD_REQUEST)
                    .body(jsonResponse);
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
}
