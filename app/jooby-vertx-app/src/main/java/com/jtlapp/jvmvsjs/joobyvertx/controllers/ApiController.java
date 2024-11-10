package com.jtlapp.jvmvsjs.joobyvertx.controllers;

import com.google.gson.JsonObject;
import com.jtlapp.jvmvsjs.vertxquery.SharedQueryRepo;
import com.jtlapp.jvmvsjs.vertxquery.VertxUtil;
import io.jooby.Context;
import io.jooby.StatusCode;
import io.jooby.annotation.GET;
import io.jooby.annotation.POST;
import io.jooby.annotation.Path;
import io.jooby.annotation.PathParam;
import io.vertx.core.Future;
import io.vertx.sqlclient.Pool;
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
    SharedQueryRepo sharedQueryRepo;

    @Inject
    Pool pgPool;

    @GET("/info")
    public CompletableFuture<String> info() {
        var gson = new JsonObject();
        gson.addProperty("appName", appName);
        gson.addProperty("appVersion", appVersion);
        gson.add("appConfig", new JsonObject());
        return CompletableFuture.completedFuture(gson.toString());
    }

    @POST("/query/{queryName}")
    public CompletableFuture<String> query(
            @PathParam String queryName, String jsonBody, Context ctx
    ) {
        var vertxFuture = sharedQueryRepo.get(queryName)
                .flatMap(sq -> sq.executeUsingGson(pgPool, jsonBody))
                .recover(e -> {
                    ctx.setResponseCode(StatusCode.SERVER_ERROR);
                    return Future.succeededFuture(toErrorJson(queryName, e));
                });
        return VertxUtil.toCompletableFuture(vertxFuture);
    }

    @GET("/sleep/{millis}")
    public CompletableFuture<String> sleep(@PathParam int millis) {
        var future = new CompletableFuture<String>();

        scheduler.schedule(() ->
                future.complete(""), millis, TimeUnit.MILLISECONDS);

        return future;
    }

    private String toErrorJson(String queryName, Throwable e) {
        return String.format("{\"query\": \"%s\", \"error\": \"%s: %s\"}",
                queryName, e.getClass().getSimpleName(), e.getMessage());
    }
}
