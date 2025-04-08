package com.jtlapp.jvmvsjs.joobyr2dbc.config;

import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.jtlapp.jvmvsjs.joobylib.ServerConfig;
import com.jtlapp.jvmvsjs.r2dbclib.Database;
import io.r2dbc.spi.ConnectionFactories;
import io.r2dbc.spi.ConnectionFactoryOptions;
import org.springframework.r2dbc.core.DatabaseClient;

import java.time.Duration;

import static io.r2dbc.spi.ConnectionFactoryOptions.*;

/**
 * Configurations that may affect load performance.
 */
public class AppConfig {

    private static final int threadCount =
            Math.max(Runtime.getRuntime().availableProcessors(), 2) * 8;

    public final ServerConfig server =
            new ServerConfig("Netty", threadCount, threadCount);
    public final R2dbcConfig dbclient = new R2dbcConfig();

    public Database createDatabase() {
        var connFactory = ConnectionFactories.get(ConnectionFactoryOptions.builder()
                .option(DRIVER, "postgresql")
                .option(HOST, System.getenv("DATABASE_HOST_NAME"))
                .option(PORT, Integer.parseInt(System.getenv("DATABASE_PORT")))
                .option(DATABASE, System.getenv("DATABASE_NAME"))
                .option(USER, System.getenv("DATABASE_USERNAME"))
                .option(PASSWORD, System.getenv("DATABASE_PASSWORD"))
                .option(CONNECT_TIMEOUT,
                        Duration.ofSeconds(dbclient.connectionTimeout))
                .build());

        var client = DatabaseClient.create(connFactory);
        return new Database(client);
    }

    public JsonNode toJsonNode(ObjectMapper mapper) {
        return mapper.valueToTree(this);
    }
}
