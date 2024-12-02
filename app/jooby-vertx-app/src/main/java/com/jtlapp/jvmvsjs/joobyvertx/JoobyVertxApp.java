package com.jtlapp.jvmvsjs.joobyvertx;

import com.fasterxml.jackson.databind.ObjectMapper;
import com.jtlapp.jvmvsjs.joobyvertx.config.AppConfig;
import com.jtlapp.jvmvsjs.vertxlib.Database;
import io.jooby.ExecutionMode;
import io.jooby.Jooby;
import io.jooby.ReactiveSupport;
import io.jooby.StatusCode;
import io.jooby.netty.NettyServer;
import io.vertx.pgclient.PgConnectOptions;
import io.vertx.sqlclient.Pool;
import io.vertx.sqlclient.PoolOptions;

import java.util.concurrent.CompletableFuture;
import java.util.concurrent.Executors;
import java.util.concurrent.TimeUnit;

public class JoobyVertxApp extends Jooby {
    public final String appName = System.getenv("APP_NAME");
    public final String appVersion = System.getenv("APP_VERSION");

    private final AppConfig appConfig = new AppConfig();
    private final Database db = createDatabase();

    {
        var scheduler = Executors.newScheduledThreadPool(1);
        var objectMapper = new ObjectMapper();
        var server = new NettyServer();

        appConfig.server.setOptions(server);
        install(server);

        use(ReactiveSupport.concurrent());

        get("/", ctx -> "Running Jooby with Netty and Vert.x Postgres\n");

        get("/api/info", ctx -> {
            var jsonObj = objectMapper.createObjectNode()
                    .put("appName", appName)
                    .put("appVersion", appVersion)
                    .set("appConfig", appConfig.toJsonNode(objectMapper));
            return CompletableFuture.completedFuture(jsonObj.toString());
        });

        get("/api/app-sleep", ctx -> {
            int millis = ctx.query("millis").intValue(0);

            scheduler.schedule(() -> {
                ctx.send("{}");
            }, millis, TimeUnit.MILLISECONDS);

            return ctx;
        }).setNonBlocking(true);

        get("/api/pg-sleep", ctx -> {
            int millis = ctx.query("millis").intValue(0);
            db.issueSleepQuery(millis)
                    .andThen(result -> ctx.send("{}"))
                    .onFailure(e -> {
                        ctx.setResponseCode(StatusCode.SERVER_ERROR);
                        ctx.send(toErrorJson("pg-sleep", e));
                    });
            return ctx;
        }).setNonBlocking(true);

        onStop(scheduler::shutdown);
    }

    private Database createDatabase() {
        PgConnectOptions connectOptions = new PgConnectOptions()
                .setHost(System.getenv("DATABASE_HOST_NAME"))
                .setPort(Integer.parseInt(System.getenv("DATABASE_PORT")))
                .setDatabase(System.getenv("DATABASE_NAME"))
                .setUser(System.getenv("DATABASE_USERNAME"))
                .setPassword(System.getenv("DATABASE_PASSWORD"))
                .setConnectTimeout(appConfig.dbclient.connectionTimeout)
                // leave prepared statement caching to pgbouncer
                .setPreparedStatementCacheMaxSize(0)
                // pgbouncer is a layer 7 proxy
                .setUseLayer7Proxy(true);

        PoolOptions poolOptions = new PoolOptions()
                .setMaxSize(appConfig.dbclient.maxPoolSize)
                .setMaxWaitQueueSize(appConfig.dbclient.maxWaitQueueSize);

        return new Database(Pool.pool(connectOptions, poolOptions));
    }

    private String toErrorJson(String endpoint, Throwable e) {
        return String.format("{\"endpoint\": \"%s\", \"error\": \"%s: %s\"}",
                endpoint, e.getClass().getSimpleName(), e.getMessage());
    }

    public static void main(final String[] args) {
        runApp(args, ExecutionMode.EVENT_LOOP, JoobyVertxApp::new);
    }
}
