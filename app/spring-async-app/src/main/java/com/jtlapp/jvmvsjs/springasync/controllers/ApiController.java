package com.jtlapp.jvmvsjs.springasync.controllers;

import com.fasterxml.jackson.databind.ObjectMapper;
import com.jtlapp.jvmvsjs.springasync.config.AppConfig;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.*;

import java.util.concurrent.CompletableFuture;
import java.util.concurrent.ScheduledExecutorService;
import java.util.concurrent.TimeUnit;

@RestController
@RequestMapping("/api")
public class ApiController {

    static final String appName = System.getenv("APP_NAME");
    static final String appVersion = System.getenv("APP_VERSION");
    static final ObjectMapper objectMapper = new ObjectMapper();

    @Autowired
    private ScheduledExecutorService scheduler;

    @Autowired
    private AppConfig appConfig;

    @GetMapping("/info")
    public CompletableFuture<String> info() {
        var jsonObj = objectMapper.createObjectNode()
                .put("appName", appName)
                .put("appVersion", appVersion)
                .set("appConfig", appConfig.toJsonNode(objectMapper));
        return CompletableFuture.completedFuture(jsonObj.toString());
    }

    @GetMapping("/app-sleep")
    public CompletableFuture<String> appSleep(@RequestParam("millis") int millis) {
        var future = new CompletableFuture<String>();

        scheduler.schedule(() -> {
            future.complete("{}");
        }, millis, TimeUnit.MILLISECONDS);

        return future;
    }
}
