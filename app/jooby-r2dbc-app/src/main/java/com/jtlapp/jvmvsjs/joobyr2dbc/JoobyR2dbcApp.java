package com.jtlapp.jvmvsjs.joobyr2dbc;

import com.jtlapp.jvmvsjs.joobyr2dbc.config.AppConfig;
import com.jtlapp.jvmvsjs.joobyr2dbc.controllers.ApiController;
import com.jtlapp.jvmvsjs.joobyr2dbc.controllers.HomeController;
import io.avaje.inject.Factory;
import io.avaje.inject.PreDestroy;
import io.jooby.ExecutionMode;
import io.jooby.Jooby;
import io.jooby.ReactiveSupport;
import io.jooby.ServerOptions;
import io.jooby.avaje.inject.AvajeInjectModule;
import io.jooby.netty.NettyServer;
import io.jooby.reactor.Reactor;
import jakarta.inject.Inject;
import jakarta.inject.Singleton;

import java.util.concurrent.ScheduledExecutorService;

@Singleton
@Factory
public class JoobyR2dbcApp extends Jooby {

    static final String appName = System.getenv("APP_NAME");;
    static final String appVersion = System.getenv("APP_VERSION");;

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
        use(Reactor.reactor());

        mvc(HomeController.class);
        mvc(ApiController.class);
    }

    @PreDestroy
    public void shutdown() {
        scheduler.shutdown();
    }

    public static void main(final String[] args) {
        runApp(args, ExecutionMode.EVENT_LOOP, JoobyR2dbcApp::new);
    }
}
