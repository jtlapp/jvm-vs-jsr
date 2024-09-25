package com.jtlapp.jvmvsjs.config;

import com.jtlapp.jvmvsjs.jdbcsharedquery.SharedQueryDB;
import com.jtlapp.jvmvsjs.jdbcsharedquery.SharedQueryRepo;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

@Configuration
public class AppConfig {

    @Value("${spring.datasource.url}")
    private String databaseURL;

    @Value("${spring.datasource.username}")
    private String username;

    @Value("${spring.datasource.password}")
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