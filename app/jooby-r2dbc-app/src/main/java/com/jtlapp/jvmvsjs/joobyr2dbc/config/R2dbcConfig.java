package com.jtlapp.jvmvsjs.joobyr2dbc.config;

public class R2dbcConfig {
    public final int maximumPoolSize = Integer.parseInt(
            System.getProperty("hikari.maximumPoolSize", "10"));
    public final int minimumIdle = Integer.parseInt(
            System.getProperty("hikari.minimumIdle", String.valueOf(maximumPoolSize)));
    public final int connectionTimeout = Integer.parseInt(
            System.getProperty("jooby.r2dbc.connect-timeout-seconds", "1800"));
}
