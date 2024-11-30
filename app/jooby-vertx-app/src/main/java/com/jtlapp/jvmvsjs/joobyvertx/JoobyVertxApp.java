package com.jtlapp.jvmvsjs.joobyvertx;

import com.jtlapp.jvmvsjs.joobyvertx.config.AppConfig;
import com.jtlapp.jvmvsjs.joobyvertx.controllers.ApiController;
import com.jtlapp.jvmvsjs.joobyvertx.controllers.HomeController;
import io.avaje.inject.Factory;
import io.avaje.inject.PreDestroy;
import io.jooby.ExecutionMode;
import io.jooby.Jooby;
import io.jooby.ReactiveSupport;
import io.jooby.ServerOptions;
import io.jooby.avaje.inject.AvajeInjectModule;
import io.jooby.netty.NettyServer;
import jakarta.inject.Inject;
import jakarta.inject.Singleton;

import java.util.concurrent.ScheduledExecutorService;

@Singleton
@Factory
public class JoobyVertxApp extends Jooby {

    @Inject
    ScheduledExecutorService scheduler;

    @Inject
    AppConfig appConfig;

    {
        install(AvajeInjectModule.of());
        var server = new NettyServer();
        appConfig.server.setOptions(server);
        install(server);

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
