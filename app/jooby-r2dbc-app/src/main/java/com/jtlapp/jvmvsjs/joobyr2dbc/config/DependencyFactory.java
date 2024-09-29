package com.jtlapp.jvmvsjs.joobyr2dbc.config;

import com.jtlapp.jvmvsjs.joobyr2dbc.JoobyR2dbcApp;
import com.jtlapp.jvmvsjs.r2dbcquery.SharedQueryRepo;
import com.typesafe.config.Config;
import io.avaje.inject.Bean;
import io.avaje.inject.External;
import io.avaje.inject.Factory;
import io.jooby.Environment;
import io.r2dbc.spi.ConnectionFactories;
import io.r2dbc.spi.ConnectionFactory;
import io.r2dbc.spi.ConnectionFactoryOptions;
import jakarta.inject.Named;
import org.springframework.r2dbc.core.DatabaseClient;

import java.util.concurrent.Executors;
import java.util.concurrent.ScheduledExecutorService;

import static io.r2dbc.spi.ConnectionFactoryOptions.*;

@Factory
public class DependencyFactory {

    @Bean
    public ScheduledExecutorService scheduledExecutorService() {
        return Executors.newScheduledThreadPool(1);
    }

    @Bean
    public ConnectionFactory connectionFactory() {
        var config = JoobyR2dbcApp.config;
        System.out.println("DATABASE_HOST_NAME: "+ config.getString("DATABASE_HOST_NAME"));
        System.out.println("DATABASE_PORT: "+ config.getString("DATABASE_PORT"));
        System.out.println("DATABASE_NAME: "+ config.getString("DATABASE_NAME"));
        System.out.println("DATABASE_USERNAME: "+ config.getString("DATABASE_USERNAME"));
        System.out.println("DATABASE_PASSWORD: "+ config.getString("DATABASE_PASSWORD"));
        return ConnectionFactories.get(ConnectionFactoryOptions.builder()
                .option(DRIVER, "postgresql")
                .option(HOST, config.getString("DATABASE_HOST_NAME"))
                .option(PORT, config.getInt("DATABASE_PORT"))
                .option(DATABASE, "DATABASE_NAME")
                .option(USER, config.getString("DATABASE_USERNAME"))
                .option(PASSWORD, "DATABASE_PASSWORD")
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
