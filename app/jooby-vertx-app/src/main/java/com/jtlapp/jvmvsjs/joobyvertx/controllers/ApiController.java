package com.jtlapp.jvmvsjs.joobyvertx.controllers;

import com.google.gson.JsonObject;
import com.jtlapp.jvmvsjs.vertxquery.Database;
import com.jtlapp.jvmvsjs.vertxquery.VertxUtil;
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

    @Inject
    ScheduledExecutorService scheduler;

    @Inject
    Database db;

    @GET("/info")
    public CompletableFuture<String> info() {
        var gson = new JsonObject();
        gson.addProperty("appName", appName);
        gson.addProperty("appVersion", appVersion);
        gson.add("appConfig", new JsonObject());
        return CompletableFuture.completedFuture(gson.toString());
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
