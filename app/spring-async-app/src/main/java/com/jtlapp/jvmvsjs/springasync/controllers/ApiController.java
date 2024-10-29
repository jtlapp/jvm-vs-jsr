package com.jtlapp.jvmvsjs.springasync.controllers;

import com.google.gson.JsonObject;
import com.jtlapp.jvmvsjs.springasync.SpringAsyncApp;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.web.bind.annotation.*;

import java.util.concurrent.CompletableFuture;
import java.util.concurrent.ScheduledExecutorService;
import java.util.concurrent.TimeUnit;

@RestController
@RequestMapping("/api")
public class ApiController {

    @Value("${application.name")
    public String appName;
    @Value("${application.version")
    public String appVersion;

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

    @GetMapping("/sleep/{millis}")
    public CompletableFuture<String> sleep(@PathVariable(name = "millis") int millis) {
        var future = new CompletableFuture<String>();

        scheduler.schedule(() -> {
            future.complete("");
        }, millis, TimeUnit.MILLISECONDS);

        return future;
    }
}
