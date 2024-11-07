package com.jtlapp.jvmvsjs.joobyr2dbc.config;

import com.jtlapp.jvmvsjs.r2dbcquery.SharedQueryRepo;
import io.avaje.inject.Bean;
import io.avaje.inject.External;
import io.avaje.inject.Factory;
import io.r2dbc.spi.ConnectionFactories;
import io.r2dbc.spi.ConnectionFactory;
import io.r2dbc.spi.ConnectionFactoryOptions;
import jakarta.inject.Named;
import org.springframework.r2dbc.core.DatabaseClient;

import java.util.concurrent.Executors;
import java.util.concurrent.ScheduledExecutorService;

import static io.r2dbc.spi.ConnectionFactoryOptions.*;

@Factory
public class AppConfig {
    public static final String version = "0.1.0";

    @Bean
    @Named("application.name")
    public String getAppName() {
        return getClass().getSimpleName();
    }

    @Bean
    @Named("application.version")
    public String getAppVersion() {
        return version;
    }

    @Bean
    public ScheduledExecutorService scheduledExecutorService() {
        return Executors.newScheduledThreadPool(1);
    }

    @Bean
    public ConnectionFactory connectionFactory(
            @External @Named("DATABASE_HOST_NAME") String hostName,
            @External @Named("DATABASE_PORT") String port,
            @External @Named("DATABASE_NAME") String databaseName,
            @External @Named("DATABASE_USERNAME") String username,
            @External @Named("DATABASE_PASSWORD") String password
    ) {
        return ConnectionFactories.get(ConnectionFactoryOptions.builder()
                .option(DRIVER, "postgresql")
                .option(HOST, hostName)
                .option(PORT, Integer.parseInt(port))
                .option(DATABASE, databaseName)
                .option(USER, username)
                .option(PASSWORD, password)
                .build());
    }

    @Bean
    public DatabaseClient databaseClient(ConnectionFactory connectionFactory) {
        return DatabaseClient.create(connectionFactory);
    }

    @Bean
    public SharedQueryRepo sharedQueryRepo(DatabaseClient databaseClient) {
        return new SharedQueryRepo(databaseClient);
    }
}
