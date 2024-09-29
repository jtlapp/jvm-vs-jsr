package com.jtlapp.jvmvsjs.springjdbc.config;

import com.jtlapp.jvmvsjs.jdbcquery.SharedQueryDB;
import com.jtlapp.jvmvsjs.jdbcquery.SharedQueryRepo;
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