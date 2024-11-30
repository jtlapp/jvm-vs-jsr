package com.jtlapp.jvmvsjs.joobynomvc;

import com.fasterxml.jackson.databind.ObjectMapper;
import com.jtlapp.jvmvsjs.joobynomvc.config.AppConfig;
import io.jooby.ExecutionMode;
import io.jooby.Jooby;
import io.jooby.ReactiveSupport;
import io.jooby.netty.NettyServer;

import java.util.concurrent.CompletableFuture;
import java.util.concurrent.Executors;
import java.util.concurrent.TimeUnit;

public class JoobyNoMVCApp extends Jooby {
    public final String appName = System.getenv("APP_NAME");
    public final String appVersion = System.getenv("APP_VERSION");
    public final AppConfig appConfig = new AppConfig();

    {
        var scheduler = Executors.newScheduledThreadPool(1);
        var objectMapper = new ObjectMapper();

        var server = new NettyServer();
        appConfig.server.setOptions(server);
        install(server);

        use(ReactiveSupport.concurrent());

        get("/", ctx -> "Running Jooby without MVC\n");

        get("/api/info", ctx -> {
            var jsonObj = objectMapper.createObjectNode()
                    .put("appName", appName)
                    .put("appVersion", appVersion)
                    .set("appConfig", appConfig.toJsonNode(objectMapper));
            return CompletableFuture.completedFuture(jsonObj.toString());
        });

        get("/api/app-sleep", ctx -> {
            int millis = ctx.query("millis").intValue(0);
            var future = new CompletableFuture<String>();

            scheduler.schedule(() -> {
                future.complete("{}");
            }, millis, TimeUnit.MILLISECONDS);

            return future;
        });

        post("/api/echoText", ctx -> {
            var body = ctx.body(String.class);
            return CompletableFuture.completedFuture(body);
        });

        post("/api/echoJackson", ctx -> {
            var body = ctx.body(String.class);
            var jackson = objectMapper.readTree(body);
            return CompletableFuture.completedFuture(jackson.toString());
        });

        onStop(scheduler::shutdown);
    }

    public static void main(final String[] args) {
        runApp(args, ExecutionMode.EVENT_LOOP, JoobyNoMVCApp::new);
    }
}
