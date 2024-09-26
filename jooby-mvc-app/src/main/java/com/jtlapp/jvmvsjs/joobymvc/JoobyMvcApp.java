package com.jtlapp.jvmvsjs.joobymvc;

import com.jtlapp.jvmvsjs.joobymvc.controllers.ApiController;
import com.jtlapp.jvmvsjs.joobymvc.controllers.HomeController;
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
public class JoobyMvcApp extends Jooby {

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
        runApp(args, ExecutionMode.EVENT_LOOP, JoobyMvcApp::new);
    }
}
