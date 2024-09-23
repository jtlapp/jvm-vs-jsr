package com.jtlapp.jvmvsjs.config;

import com.jtlapp.jvmvsjs.SharedQueryDB;
import com.jtlapp.jvmvsjs.SharedQueryRepo;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

@Configuration
public class AppConfig {

    @Value("${database.url}")
    private String databaseURL;

    @Value("${database.username}")
    private String username;

    @Value("${database.password}")
    private String password;

    @Bean
    public SharedQueryDB getSharedQueryDB() {
        return new SharedQueryDB(databaseURL, username, password);
    }

    @Bean
    public SharedQueryRepo getSharedQueryRepo() {
        return new SharedQueryRepo();
    }
}