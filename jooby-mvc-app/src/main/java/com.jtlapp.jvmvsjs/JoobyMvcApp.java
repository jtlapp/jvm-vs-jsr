package com.jtlapp.jvmvsjs;

import io.jooby.ExecutionMode;
import io.jooby.Jooby;
import io.jooby.ReactiveSupport;
import io.jooby.ServerOptions;
import io.jooby.netty.NettyServer;

import java.util.concurrent.Executors;

public class JoobyMvcApp extends Jooby {

  {
    var scheduler = Executors.newScheduledThreadPool(1);

    install(new NettyServer().setOptions(
            new ServerOptions()
                    .setIoThreads(Runtime.getRuntime().availableProcessors() + 1)
                    .setWorkerThreads(Runtime.getRuntime().availableProcessors() + 1)
    ));
    use(ReactiveSupport.concurrent());

    //mvc(new HomeController());

    onStop(scheduler::shutdown);
  }

  public static void main(final String[] args) {
    runApp(args, ExecutionMode.EVENT_LOOP, JoobyMvcApp::new);
  }
}
