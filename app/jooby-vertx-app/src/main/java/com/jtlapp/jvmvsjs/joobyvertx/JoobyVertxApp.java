package com.jtlapp.jvmvsjs.joobyvertx;

import com.jtlapp.jvmvsjs.joobyvertx.controllers.ApiController;
import com.jtlapp.jvmvsjs.joobyvertx.controllers.HomeController;
import io.avaje.inject.Bean;
import io.avaje.inject.Factory;
import io.avaje.inject.PreDestroy;
import io.jooby.ExecutionMode;
import io.jooby.Jooby;
import io.jooby.ReactiveSupport;
import io.jooby.ServerOptions;
import io.jooby.avaje.inject.AvajeInjectModule;
import io.jooby.netty.NettyServer;
import jakarta.inject.Inject;
import jakarta.inject.Named;
import jakarta.inject.Singleton;

import java.util.concurrent.ScheduledExecutorService;

@Singleton
@Factory
public class JoobyVertxApp extends Jooby {
    public static final String version = "0.1.0";

    @Bean
    @Named("application.name")
    public String getAppName() {
        return getClass().getSimpleName();
    }

    @Bean
    @Named("application.version")
    public String getAppVersion() {
        return version;
    }

    @Inject
    ScheduledExecutorService scheduler;

    {
        install(AvajeInjectModule.of());
        install(new NettyServer().setOptions(
                new ServerOptions()
                        .setIoThreads(Runtime.getRuntime().availableProcessors() + 1)
                        .setWorkerThreads(Runtime.getRuntime().availableProcessors() + 1)
        ));

        use(ReactiveSupport.concurrent());

        mvc(HomeController.class);
        mvc(ApiController.class);
    }

    @PreDestroy
    public void shutdown() {
        scheduler.shutdown();
    }

    public static void main(final String[] args) {
        runApp(args, ExecutionMode.EVENT_LOOP, JoobyVertxApp::new);
    }
}
