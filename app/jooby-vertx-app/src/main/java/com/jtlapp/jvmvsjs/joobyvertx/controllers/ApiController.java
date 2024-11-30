package com.jtlapp.jvmvsjs.joobyvertx.controllers;

import com.fasterxml.jackson.databind.ObjectMapper;
import com.jtlapp.jvmvsjs.joobyvertx.config.AppConfig;
import com.jtlapp.jvmvsjs.vertxlib.Database;
import com.jtlapp.jvmvsjs.vertxlib.VertxUtil;
import io.jooby.Context;
import io.jooby.StatusCode;
import io.jooby.annotation.*;
import io.vertx.core.Future;
import jakarta.inject.Inject;
import jakarta.inject.Singleton;

import java.util.concurrent.CompletableFuture;
import java.util.concurrent.ScheduledExecutorService;
import java.util.concurrent.TimeUnit;

@Singleton
@Path("/api")
public class ApiController {

    static final String appName = System.getenv("APP_NAME");;
    static final String appVersion = System.getenv("APP_VERSION");;
    static final ObjectMapper objectMapper = new ObjectMapper();

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

        scheduler.schedule(() ->
                future.complete("{}"), millis, TimeUnit.MILLISECONDS);

        return future;
    }

    @GET("/pg-sleep")
    public CompletableFuture<String> pgSleep(@QueryParam int millis, Context ctx
    ) {
        var vertxFuture = db.issueSleepQuery(millis)
                .map(result -> "{}")
                .recover(e -> {
                    ctx.setResponseCode(StatusCode.SERVER_ERROR);
                    return Future.succeededFuture(toErrorJson("pg-sleep", e));
                });
        return VertxUtil.toCompletableFuture(vertxFuture);
    }

    private String toErrorJson(String endpoint, Throwable e) {
        return String.format("{\"endpoint\": \"%s\", \"error\": \"%s: %s\"}",
        endpoint, e.getClass().getSimpleName(), e.getMessage());
    }
}
