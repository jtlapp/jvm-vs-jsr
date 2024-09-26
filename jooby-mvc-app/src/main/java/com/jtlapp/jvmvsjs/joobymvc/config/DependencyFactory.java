package com.jtlapp.jvmvsjs.joobymvc.config;

import io.avaje.inject.Bean;
import io.avaje.inject.Factory;

import java.util.concurrent.Executors;
import java.util.concurrent.ScheduledExecutorService;

@Factory
public class DependencyFactory {

    @Bean
    public ScheduledExecutorService scheduledExecutorService() {
        return Executors.newScheduledThreadPool(1);
    }
}
