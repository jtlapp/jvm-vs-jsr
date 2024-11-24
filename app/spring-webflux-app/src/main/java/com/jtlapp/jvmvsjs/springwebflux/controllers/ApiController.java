package com.jtlapp.jvmvsjs.springwebflux.controllers;

import com.google.gson.JsonObject;
import com.jtlapp.jvmvsjs.r2dbcquery.Database;
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

    static final String appName = System.getenv("APP_NAME");;
    static final String appVersion = System.getenv("APP_VERSION");;

    @Autowired
    private ScheduledExecutorService scheduler;

    @Autowired
    private Database db;

    @GetMapping("/info")
    public Mono<String> info() {
        var gson = new JsonObject();
        gson.addProperty("appName", appName);
        gson.addProperty("appVersion", appVersion);
        gson.add("appConfig", new JsonObject());
        return Mono.just(gson.toString());
    }

    @GetMapping("/app-sleep")
    public Mono<String> appSleep(@RequestParam int millis) {
        return Mono.create(sink ->
                scheduler.schedule(() -> sink.success(""), millis, TimeUnit.MILLISECONDS)
        );
    }

    @GetMapping("/pg-sleep")
    public Mono<ResponseEntity<String>> pgSleep(@RequestParam int millis) {
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
