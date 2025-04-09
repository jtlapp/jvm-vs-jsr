package com.jtlapp.jvmvsjs.webfluxlib;

import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Component;

@Component
public class R2dbcConfig {

    @Value("${spring.r2dbc.pool.max-size}")
    public long maxConnections;

    @Value("${spring.r2dbc.pool.initial-size}")
    public int initialConnections;
}
