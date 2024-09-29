package com.jtlapp.jvmvsjs.joobynomvc;

import com.fasterxml.jackson.databind.ObjectMapper;
import com.google.gson.JsonParser;
import io.jooby.ExecutionMode;
import io.jooby.Jooby;
import io.jooby.ReactiveSupport;
import io.jooby.ServerOptions;
import io.jooby.netty.NettyServer;

import java.util.concurrent.CompletableFuture;
import java.util.concurrent.Executors;
import java.util.concurrent.TimeUnit;

public class JoobyNoMVCApp extends Jooby {

    {
        var scheduler = Executors.newScheduledThreadPool(1);
        var objectMapper = new ObjectMapper();

        install(new NettyServer().setOptions(
                new ServerOptions()
                        .setIoThreads(Runtime.getRuntime().availableProcessors() + 1)
                        .setWorkerThreads(Runtime.getRuntime().availableProcessors() + 1)
        ));

        use(ReactiveSupport.concurrent());

        get("/", ctx -> "Running Jooby without MVC\n");

        post("/api/echoText", ctx -> {
            var body = ctx.body(String.class);
            return CompletableFuture.completedFuture(body);
        });

        post("/api/echoGson", ctx -> {
            var body = ctx.body(String.class);
            var gson = JsonParser.parseString(body).getAsJsonObject();
            return CompletableFuture.completedFuture(gson.toString());
        });

        post("/api/echoJackson", ctx -> {
            var body = ctx.body(String.class);
            var jackson = objectMapper.readTree(body);
            return CompletableFuture.completedFuture(jackson.toString());
        });

        get("/api/sleep/{millis}", ctx -> {
            int millis = ctx.path("millis").intValue();
            var future = new CompletableFuture<String>();

            scheduler.schedule(() -> {
                future.complete("");
            }, millis, TimeUnit.MILLISECONDS);

            return future;
        });

        onStop(scheduler::shutdown);
    }

    public static void main(final String[] args) {
        runApp(args, ExecutionMode.EVENT_LOOP, JoobyNoMVCApp::new);
    }
}
