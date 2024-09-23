package com.jtlapp.jvmvsjs.controllers;

import com.jtlapp.jvmvsjs.SharedQuery;
import com.jtlapp.jvmvsjs.SharedQueryDB;
import com.jtlapp.jvmvsjs.SharedQueryException;
import com.jtlapp.jvmvsjs.SharedQueryRepo;
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
}
