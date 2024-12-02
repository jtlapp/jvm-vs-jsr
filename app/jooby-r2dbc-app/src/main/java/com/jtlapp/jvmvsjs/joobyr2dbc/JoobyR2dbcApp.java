package com.jtlapp.jvmvsjs.joobyr2dbc;

import com.fasterxml.jackson.databind.ObjectMapper;
import com.jtlapp.jvmvsjs.joobyr2dbc.config.AppConfig;
import com.jtlapp.jvmvsjs.r2dbclib.Database;
import io.jooby.*;
import io.jooby.exception.StatusCodeException;
import io.jooby.netty.NettyServer;
import io.jooby.reactor.Reactor;
import io.r2dbc.spi.ConnectionFactories;
import io.r2dbc.spi.ConnectionFactoryOptions;
import org.springframework.r2dbc.core.DatabaseClient;

import java.time.Duration;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.Executors;
import java.util.concurrent.TimeUnit;

import static io.r2dbc.spi.ConnectionFactoryOptions.*;

public class JoobyR2dbcApp extends Jooby {

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
        use(Reactor.reactor());

        get("/", ctx -> "Running Jooby with Netty and R2DBC\n");

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

        get("/api/pg-sleep", ctx -> {
            int millis = ctx.query("millis").intValue(0);
            return db.issueSleepQuery(millis)
                    .thenReturn("{}")
                    .onErrorMap(e ->
                            new StatusCodeException(StatusCode.SERVER_ERROR,
                                    toErrorJson("pg-sleep", e))
                    );
        });

        onStop(scheduler::shutdown);
    }

    private Database createDatabase() {
        var connFactory = ConnectionFactories.get(ConnectionFactoryOptions.builder()
                .option(DRIVER, "postgresql")
                .option(HOST, System.getenv("DATABASE_HOST_NAME"))
                .option(PORT, Integer.parseInt(System.getenv("DATABASE_PORT")))
                .option(DATABASE, System.getenv("DATABASE_NAME"))
                .option(USER, System.getenv("DATABASE_USERNAME"))
                .option(PASSWORD, System.getenv("DATABASE_PASSWORD"))
                .option(CONNECT_TIMEOUT,
                        Duration.ofSeconds(appConfig.dbclient.connectionTimeout))
                .build());

        var client = DatabaseClient.create(connFactory);
        return new Database(client);
    }

    private String toErrorJson(String endpoint, Throwable e) {
        return String.format("{\"endpoint\": \"%s\", \"error\": \"%s: %s\"}",
                endpoint, e.getClass().getSimpleName(), e.getMessage());
    }

    public static void main(final String[] args) {
        runApp(args, ExecutionMode.EVENT_LOOP, JoobyR2dbcApp::new);
    }
}
