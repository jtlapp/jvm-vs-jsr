package com.jtlapp.jvmvsjs.springwebflux.controllers;

import com.jtlapp.jvmvsjs.r2dbcquery.SharedQueryRepo;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.HttpStatus;
import org.springframework.http.MediaType;
import org.springframework.http.ResponseEntity;
import org.springframework.r2dbc.core.DatabaseClient;
import org.springframework.web.bind.annotation.*;
import org.springframework.web.reactive.function.server.ServerResponse;
import org.springframework.web.server.ResponseStatusException;
import reactor.core.publisher.Mono;

import java.util.concurrent.ScheduledExecutorService;
import java.util.concurrent.TimeUnit;

@RestController
@RequestMapping("/api")
public class ApiController {

    @Autowired
    private ScheduledExecutorService scheduler;

    @Autowired
    private DatabaseClient db;

    @Autowired
    private SharedQueryRepo sharedQueryRepo;

    @PostMapping("/query/{queryName}")
    public Mono<ResponseEntity<String>> query(
            @PathVariable(name = "queryName") String queryName,
            @RequestBody String jsonBody
    ) {
        return sharedQueryRepo.get(queryName)
                .flatMap(sq -> sq.executeUsingGson(db, jsonBody))
                .map(json -> ResponseEntity.ok().body(json))
                .onErrorResume(e -> Mono.just(ResponseEntity
                        .status(HttpStatus.INTERNAL_SERVER_ERROR)
                        .body(toErrorJson(queryName, e))));
    }

    @GetMapping("/sleep/{millis}")
    public Mono<String> sleep(@PathVariable(name = "millis") int millis) {
        return Mono.create(sink ->
                scheduler.schedule(() -> sink.success(""), millis, TimeUnit.MILLISECONDS)
        );
    }

    private String toErrorJson(String queryName, Throwable e) {
        return String.format("{\"query\": \"%s\", \"error\": \"%s: %s\"}",
                queryName, e.getClass().getSimpleName(), e.getMessage());
    }
}
