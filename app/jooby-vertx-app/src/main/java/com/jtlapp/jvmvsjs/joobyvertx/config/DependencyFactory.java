package com.jtlapp.jvmvsjs.joobyvertx.config;

import com.jtlapp.jvmvsjs.vertxquery.SharedQueryRepo;
import io.avaje.inject.Bean;
import io.avaje.inject.External;
import io.avaje.inject.Factory;
import io.vertx.pgclient.PgConnectOptions;
import io.vertx.sqlclient.Pool;
import io.vertx.sqlclient.PoolOptions;
import jakarta.inject.Named;

import java.util.concurrent.Executors;
import java.util.concurrent.ScheduledExecutorService;

@Factory
public class DependencyFactory {

    @Bean
    public ScheduledExecutorService scheduledExecutorService() {
        return Executors.newScheduledThreadPool(1);
    }

    @Bean
    public Pool pgPool(
            @External @Named("DATABASE_HOST_NAME") String hostName,
            @External @Named("DATABASE_PORT") String port,
            @External @Named("DATABASE_NAME") String databaseName,
            @External @Named("DATABASE_USERNAME") String username,
            @External @Named("DATABASE_PASSWORD") String password
    ) {
        PgConnectOptions connectOptions = new PgConnectOptions()
                .setHost(hostName)
                .setPort(Integer.parseInt(port))
                .setDatabase(databaseName)
                .setUser(username)
                .setPassword(password)
                // leave prepared statement caching to pgbouncer
                .setPreparedStatementCacheMaxSize(0)
                // pgbouncer is a layer 7 proxy
                .setUseLayer7Proxy(true);

        // TODO: What pool size?
        PoolOptions poolOptions = new PoolOptions().setMaxSize(5);

        // TODO: Pool vs PgConnection?
        return Pool.pool(connectOptions, poolOptions);
    }

    @Bean
    public SharedQueryRepo sharedQueryRepo(Pool pgPool) {
        return new SharedQueryRepo(pgPool);
    }
}
