package com.jtlapp.jvmvsjs.springwebflux.controllers;

import com.fasterxml.jackson.databind.ObjectMapper;
import com.jtlapp.jvmvsjs.r2dbclib.Database;
import com.jtlapp.jvmvsjs.springwebflux.config.AppConfig;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;
import reactor.core.publisher.Mono;

import java.util.concurrent.ScheduledExecutorService;
import java.util.concurrent.TimeUnit;

@RestController
@RequestMapping("/api")
public class ApiController {

    static final String APP_NAME = System.getenv("APP_NAME");;
    static final String appVersion = "0.1.0";;
    static final ObjectMapper objectMapper = new ObjectMapper();

    @Autowired
    private ScheduledExecutorService scheduler;

    @Autowired
    private AppConfig appConfig;

    @Autowired
    private Database db;

    @GetMapping("/info")
    public Mono<String> info() {
        var jsonObj = objectMapper.createObjectNode()
                .put("appName", APP_NAME)
                .put("appVersion", appVersion)
                .set("appConfig", appConfig.toJsonNode(objectMapper));
        return Mono.just(jsonObj.toString());
    }

    @GetMapping("/app-sleep")
    public Mono<String> appSleep(@RequestParam("millis") int millis) {
        return Mono.create(sink ->
                scheduler.schedule(() -> sink.success(""), millis, TimeUnit.MILLISECONDS)
        );
    }

    @GetMapping("/pg-sleep")
    public Mono<ResponseEntity<String>> pgSleep(@RequestParam("millis") int millis) {
        return db.issueSleepQuery(millis)
                .thenReturn(ResponseEntity.ok().body("{}"))
                .onErrorResume(e -> Mono.just(ResponseEntity
                        .status(HttpStatus.INTERNAL_SERVER_ERROR)
                        .body(toErrorJson("pg-sleep", e))));
    }

    private String toErrorJson(String endpoint, Throwable e) {
        return String.format("{\"endpoint\": \"%s\", \"error\": \"%s: %s\"}",
                endpoint, e.getClass().getSimpleName(), e.getMessage());
    }
}
