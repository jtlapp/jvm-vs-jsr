package com.jtlapp.jvmvsjs.joobyr2dbc.controllers;

import com.fasterxml.jackson.databind.ObjectMapper;
import com.google.gson.JsonObject;
import com.jtlapp.jvmvsjs.joobyr2dbc.config.AppConfig;
import com.jtlapp.jvmvsjs.r2dbclib.Database;
import io.jooby.Context;
import io.jooby.annotation.*;
import jakarta.inject.Inject;
import jakarta.inject.Singleton;
import reactor.core.publisher.Mono;

import java.util.concurrent.CompletableFuture;
import java.util.concurrent.ScheduledExecutorService;
import java.util.concurrent.TimeUnit;

@Singleton
@Path("/api")
public class ApiController {

    static final String appName = System.getenv("APP_NAME");
    static final String appVersion = System.getenv("APP_VERSION");

    ObjectMapper objectMapper = new ObjectMapper();

    @Inject
    ScheduledExecutorService scheduler;
    @Inject
    AppConfig appConfig;
    @Inject
    Database db;

    @GET("/info")
    public CompletableFuture<String> info() {
        var jsonObj = objectMapper.createObjectNode()
                .put("appName", appName)
                .put("appVersion", appVersion)
                .set("appConfig", appConfig.toJsonNode(objectMapper));
        return CompletableFuture.completedFuture(jsonObj.toString());
    }

    @GET("/app-sleep")
    public CompletableFuture<String> appSleep(@QueryParam int millis) {
        var future = new CompletableFuture<String>();

        scheduler.schedule(() -> {
            future.complete("{}");
        }, millis, TimeUnit.MILLISECONDS);

        return future;
    }

    @GET("/pg-sleep")
    public Mono<String> pgSleep(@QueryParam int millis, Context ctx) {
        return db.issueSleepQuery(millis)
                .map(result -> "{}")
                .onErrorResume(e -> Mono.just(toErrorJson("pg-sleep", e)));
    }

    private String toErrorJson(String endpoint, Throwable e) {
        return String.format("{\"endpoint\": \"%s\", \"error\": \"%s: %s\"}",
                endpoint, e.getClass().getSimpleName(), e.getMessage());
    }
}
