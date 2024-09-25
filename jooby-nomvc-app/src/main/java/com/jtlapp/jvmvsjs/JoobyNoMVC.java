package com.jtlapp.jvmvsjs;

import io.jooby.ExecutionMode;
import io.jooby.Jooby;
import io.jooby.ReactiveSupport;
import io.jooby.ServerOptions;
import io.jooby.netty.NettyServer;

import java.util.concurrent.CompletableFuture;
import java.util.concurrent.Executors;
import java.util.concurrent.TimeUnit;

public class JoobyNoMVC extends Jooby {

  {
    var scheduler = Executors.newScheduledThreadPool(1);

    install(new NettyServer().setOptions(
            new ServerOptions()
                .setIoThreads(Runtime.getRuntime().availableProcessors() + 1)
                .setWorkerThreads(Runtime.getRuntime().availableProcessors() + 1)
    ));

    use(ReactiveSupport.concurrent());

    get("/", ctx -> "Running Jooby without MVC\n");

    get("/api/sleep/{millis}", ctx -> {
      int millis = ctx.path("millis").intValue();
      var future = new CompletableFuture<String>();

      scheduler.schedule(() -> {
        future.complete("");
      }, millis, TimeUnit.MILLISECONDS);

      return future;
    });

    onStop(scheduler::shutdown);
  }

  public static void main(final String[] args) {
    runApp(args, ExecutionMode.EVENT_LOOP, JoobyNoMVC::new);
  }
}
