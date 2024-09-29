package com.jtlapp.jvmvsjs.springasync.controllers;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.*;

import java.util.concurrent.CompletableFuture;
import java.util.concurrent.ScheduledExecutorService;
import java.util.concurrent.TimeUnit;

@RestController
@RequestMapping("/api")
public class ApiController {

    @Autowired
    private ScheduledExecutorService scheduler;

    @GetMapping("/sleep/{millis}")
    public CompletableFuture<String> sleep(@PathVariable(name = "millis") int millis) {
        var future = new CompletableFuture<String>();

        scheduler.schedule(() -> {
            future.complete("");
        }, millis, TimeUnit.MILLISECONDS);

        return future;
    }
}
