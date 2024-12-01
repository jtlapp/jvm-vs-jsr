package com.jtlapp.jvmvsjs.joobyvertx.config;

public class VertxConfig {
    public final int maxPoolSize = Integer.parseInt(
            System.getProperty("vertx.maxPoolSize", "20"));

    public final int maxWaitQueueSize = Integer.parseInt(
            System.getProperty("vertx.maxWaitQueueSize", "1000"));
    public final int connectionTimeout = Integer.parseInt(
            System.getProperty("vertx.connectionTimeout", "60000"));
}
