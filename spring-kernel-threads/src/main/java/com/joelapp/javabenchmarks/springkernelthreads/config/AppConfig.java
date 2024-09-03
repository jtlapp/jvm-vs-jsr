package com.joelapp.javabenchmarks.springkernelthreads.config;

import com.joelapp.javabenchmarks.jdbcquerylibrary.JdbcQueries;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

@Configuration
public class AppConfig {

    @Value("${database.url}")
    private String jdbcURL;

    @Value("${database.username}")
    private String username;

    @Value("${database.password}")
    private String password;

    @Bean
    public JdbcQueries queries() {
        return new JdbcQueries(jdbcURL, username, password);
    }
}