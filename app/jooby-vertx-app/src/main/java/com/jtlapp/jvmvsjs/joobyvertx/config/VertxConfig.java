package com.jtlapp.jvmvsjs.joobyvertx.config;

import com.jtlapp.jvmvsjs.javalib.AppProperties;

public class VertxConfig {
    public final int maxPoolSize = Integer.parseInt(
            AppProperties.get("vertx.maxPoolSize"));

    public final int maxWaitQueueSize = Integer.parseInt(
            AppProperties.get("vertx.maxWaitQueueSize"));

    public final int connectionTimeout = Integer.parseInt(
            AppProperties.get("vertx.connectionTimeout"));
}
