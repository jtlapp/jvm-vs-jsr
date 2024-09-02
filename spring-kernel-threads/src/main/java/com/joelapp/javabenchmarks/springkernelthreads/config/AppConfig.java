package com.joelapp.javabenchmarks.springkernelthreads.config;

import com.joelapp.javabenchmarks.hikariquerylibrary.HikariQueries;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

@Configuration
public class AppConfig {

    @Value("${spring.datasource.url}")
    private String jdbcUrl;

    @Value("${spring.datasource.username}")
    private String username;

    @Value("${spring.datasource.password}")
    private String password;

    @Bean
    public HikariQueries hikariQueries() {
        return new HikariQueries(jdbcUrl, username, password);
    }
}