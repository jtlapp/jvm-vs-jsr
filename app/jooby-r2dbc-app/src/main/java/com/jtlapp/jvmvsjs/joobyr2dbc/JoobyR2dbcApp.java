package com.jtlapp.jvmvsjs.joobyr2dbc;

import com.fasterxml.jackson.databind.ObjectMapper;
import com.jtlapp.jvmvsjs.javalib.AppProperties;
import com.jtlapp.jvmvsjs.joobyr2dbc.config.AppConfig;
import com.jtlapp.jvmvsjs.r2dbclib.Database;
import io.jooby.*;
import io.jooby.exception.StatusCodeException;
import io.jooby.netty.NettyServer;
import io.jooby.reactor.Reactor;
import reactor.core.publisher.Mono;

import java.time.Duration;
import java.util.concurrent.CompletableFuture;

public class JoobyR2dbcApp extends Jooby {

    public final String APP_NAME = System.getenv("APP_NAME");
    public final String APP_VERSION = "0.1.0";

    private final AppConfig appConfig;
    private final Database db;

    {
        var objectMapper = new ObjectMapper();
        var server = new NettyServer();

        AppProperties.init(JoobyR2dbcApp.class.getClassLoader());
        appConfig = new AppConfig();
        appConfig.server.setOptions(server);
        db = appConfig.createDatabase();
        install(server);

        use(ReactiveSupport.concurrent());
        use(Reactor.reactor());

        get("/", ctx -> "Running Jooby with Netty and R2DBC\n");

        get("/api/info", ctx -> {
            var jsonObj = objectMapper.createObjectNode()
                    .put("appName", APP_NAME)
                    .put("appVersion", APP_VERSION)
                    .set("appConfig", appConfig.toJsonNode(objectMapper));
            return CompletableFuture.completedFuture(jsonObj.toString());
        });

        get("/api/app-sleep", ctx -> {
            int millis = ctx.query("millis").intValue(0);
            return Mono.delay(Duration.ofMillis(millis))
                    .thenReturn("{}")
                    .onErrorMap(e ->
                            new StatusCodeException(StatusCode.SERVER_ERROR,
                                    toErrorJson("app-sleep", e))
                    );
        });

        get("/api/pg-sleep", ctx -> {
            int millis = ctx.query("millis").intValue(0);
            return db.issueSleepQuery(millis)
                    .thenReturn("{}")
                    .onErrorMap(e ->
                            new StatusCodeException(StatusCode.SERVER_ERROR,
                                    toErrorJson("pg-sleep", e))
                    );
        });
    }

    private String toErrorJson(String endpoint, Throwable e) {
        return String.format("{\"endpoint\": \"%s\", \"error\": \"%s: %s\"}",
                endpoint, e.getClass().getSimpleName(), e.getMessage());
    }

    public static void main(final String[] args) {
        runApp(args, ExecutionMode.EVENT_LOOP, JoobyR2dbcApp::new);
    }
}
