package com.jtlapp.jvmvsjs.joobyr2dbc.config;

import com.jtlapp.jvmvsjs.javalib.AppProperties;

public class R2dbcConfig {
    public final int maximumPoolSize = Integer.parseInt(
            AppProperties.get("jooby.hikari.maximumPoolSize"));
    public final int minimumIdle = Integer.parseInt(
            AppProperties.get("jooby.hikari.minimumIdle"));
    public final int connectionTimeout = Integer.parseInt(
            AppProperties.get("jooby.r2dbc.connect-timeout-seconds"));
}
