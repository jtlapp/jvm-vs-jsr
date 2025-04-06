package com.jtlapp.jvmvsjs.joobyvertx.config;

import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.jtlapp.jvmvsjs.joobylib.ServerConfig;
import com.jtlapp.jvmvsjs.vertxlib.Database;
import io.vertx.pgclient.PgConnectOptions;
import io.vertx.sqlclient.Pool;
import io.vertx.sqlclient.PoolOptions;

/**
 * Configurations that may affect load performance.
 */
public class AppConfig {

    private static final int threadCount =
            Math.max(Runtime.getRuntime().availableProcessors(), 2) * 8;

    public final ServerConfig server =
            new ServerConfig("Netty", threadCount, threadCount);
    public final VertxConfig dbclient = new VertxConfig();

    public Database createDatabase() {
        PgConnectOptions connectOptions = new PgConnectOptions()
                .setHost(System.getenv("DATABASE_HOST_NAME"))
                .setPort(Integer.parseInt(System.getenv("DATABASE_PORT")))
                .setDatabase(System.getenv("DATABASE_NAME"))
                .setUser(System.getenv("DATABASE_USERNAME"))
                .setPassword(System.getenv("DATABASE_PASSWORD"))
                .setConnectTimeout(dbclient.connectionTimeout)
                // leave prepared statement caching to pgbouncer
                .setPreparedStatementCacheMaxSize(0)
                // pgbouncer is a layer 7 proxy
                .setUseLayer7Proxy(true);

        PoolOptions poolOptions = new PoolOptions()
                .setMaxSize(dbclient.maxPoolSize)
                .setMaxWaitQueueSize(dbclient.maxWaitQueueSize);

        return new Database(Pool.pool(connectOptions, poolOptions));
    }

    public JsonNode toJsonNode(ObjectMapper mapper) {
        return mapper.valueToTree(this);
    }
}
