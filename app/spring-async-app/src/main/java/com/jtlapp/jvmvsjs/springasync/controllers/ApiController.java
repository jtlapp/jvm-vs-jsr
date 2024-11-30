package com.jtlapp.jvmvsjs.springasync.controllers;

import com.google.gson.JsonObject;
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

    @Autowired
    private ScheduledExecutorService scheduler;

    @GetMapping("/info")
    public CompletableFuture<String> info() {
        var gson = new JsonObject();
        gson.addProperty("appName", appName);
        gson.addProperty("appVersion", appVersion);
        gson.add("appConfig", new JsonObject());
        return CompletableFuture.completedFuture(gson.toString());
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
