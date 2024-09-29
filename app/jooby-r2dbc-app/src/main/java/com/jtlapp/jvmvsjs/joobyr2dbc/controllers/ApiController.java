package com.jtlapp.jvmvsjs.joobyr2dbc.controllers;

import com.jtlapp.jvmvsjs.r2dbcquery.SharedQueryRepo;
import io.jooby.Context;
import io.jooby.StatusCode;
import io.jooby.annotation.GET;
import io.jooby.annotation.POST;
import io.jooby.annotation.Path;
import io.jooby.annotation.PathParam;
import jakarta.inject.Inject;
import jakarta.inject.Singleton;
import org.springframework.r2dbc.core.DatabaseClient;
import reactor.core.publisher.Mono;

import java.util.concurrent.CompletableFuture;
import java.util.concurrent.ScheduledExecutorService;
import java.util.concurrent.TimeUnit;

@Singleton
@Path("/api")
public class ApiController {

    @Inject
    ScheduledExecutorService scheduler;

    @Inject
    SharedQueryRepo sharedQueryRepo;

    @Inject
    DatabaseClient db;

    @POST("/query/{queryName}")
    public Mono<String> query(
            @PathParam String queryName, String jsonBody, Context ctx
    ) {
        return sharedQueryRepo.get(queryName)
                .flatMap(sq -> sq.executeUsingGson(db, jsonBody))
                .onErrorResume(e -> {
                    ctx.setResponseCode(StatusCode.SERVER_ERROR);
                    return Mono.just(toErrorJson(queryName, e));
                });
    }

//    public Mono<ResponseEntity<String>> query(
//            String queryName, String jsonBody, Context ctx
//    ) {
//        return sharedQueryRepo.get(queryName)
//                .flatMap(sq -> sq.executeUsingGson(db, jsonBody))
//                .map(json -> ResponseEntity.ok().body(json));
//                .onErrorResume(e -> Mono.just(ResponseEntity
//                        .status(HttpStatus.INTERNAL_SERVER_ERROR)
//                        .body(toErrorJson(queryName, e))));
//    }

    @GET("/sleep/{millis}")
    public CompletableFuture<String> sleep(@PathParam int millis) {
        var future = new CompletableFuture<String>();

        scheduler.schedule(() -> {
            future.complete("");
        }, millis, TimeUnit.MILLISECONDS);

        return future;
    }

    private String toErrorJson(String queryName, Throwable e) {
        return String.format("{\"query\": \"%s\", \"error\": \"%s\"}",
                queryName, e.getMessage());
    }
}
