package com.jtlapp.jvmvsjs.springwebflux.config;

import com.jtlapp.jvmvsjs.r2dbcquery.SharedQueryRepo;
import io.r2dbc.spi.ConnectionFactory;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.r2dbc.core.DatabaseClient;

import java.util.concurrent.Executors;
import java.util.concurrent.ScheduledExecutorService;

@Configuration
public class AppConfig {

    @Value("${spring.r2dbc.url}")
    private String databaseURL;

    @Value("${spring.r2dbc.username}")
    private String username;

    @Value("${spring.r2dbc.password}")
    private String password;

    @Bean
    public ScheduledExecutorService scheduledExecutorService() {
        return Executors.newScheduledThreadPool(1);
    }

    // Spring Boot automatically provides the ConnectionFactory based on the
    // sprint.r2dbc.{url, username, password} application properties.
    @Bean
    public DatabaseClient databaseClient(ConnectionFactory connectionFactory) {
        return DatabaseClient.create(connectionFactory);
    }

    @Bean
    public SharedQueryRepo sharedQueryRepo(DatabaseClient databaseClient) {
        return new SharedQueryRepo(databaseClient);
    }
}