package com.jtlapp.jvmvsjs.springjdbc.config;

import org.springframework.beans.factory.annotation.Value;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import com.jtlapp.jvmvsjs.jdbcquery.Database;

@Configuration
public class AppConfig {

    @Value("${spring.datasource.url}")
    private String databaseURL;

    @Value("${spring.datasource.username}")
    private String username;

    @Value("${spring.datasource.password}")
    private String password;

    @Bean
    public Database database() {
        return new Database(databaseURL, username, password);
    }
}