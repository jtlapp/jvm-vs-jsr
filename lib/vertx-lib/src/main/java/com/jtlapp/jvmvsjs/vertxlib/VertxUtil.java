package com.jtlapp.jvmvsjs.vertxlib;

import io.vertx.core.Future;
        import java.util.concurrent.CompletableFuture;

public class VertxUtil {

    public static <T> CompletableFuture<T> toCompletableFuture(Future<T> vertxFuture) {
        CompletableFuture<T> completableFuture = new CompletableFuture<>();

        vertxFuture.onComplete(result -> {
            if (result.succeeded()) {
                completableFuture.complete(result.result());
            } else {
                completableFuture.completeExceptionally(result.cause());
            }
        });

        return completableFuture;
    }
}
