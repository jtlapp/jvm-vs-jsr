package com.jtlapp.jvmvsjs;

import io.jooby.ExecutionMode;
import io.jooby.Jooby;
import io.jooby.ReactiveSupport;
import io.jooby.StatusCode;
import io.jooby.exception.BadRequestException;
import io.jooby.netty.NettyServer;

import java.util.concurrent.CompletableFuture;
import java.util.concurrent.Executors;
import java.util.concurrent.ScheduledExecutorService;
import java.util.concurrent.TimeUnit;

public class JoobyNoMVCJdbcApp extends Jooby {

  {
    ScheduledExecutorService scheduler = Executors.newScheduledThreadPool(1);
    SharedQueryDB sharedQueryDB = new SharedQueryDB(
            System.getenv("DATABASE_URL"),
            System.getenv("DATABASE_USERNAME"),
            System.getenv("DATABASE_PASSWORD"));
    SharedQueryRepo sharedQueryRepo = new SharedQueryRepo();

    install(new NettyServer());

    use(ReactiveSupport.concurrent());

    get("/", ctx -> "Running Jooby without MVC, with JDBC\n");

    post("/api/query/{queryName}", ctx -> {
      String queryName = ctx.path("queryName").value();
      String jsonBody = ctx.body().value();
      try {
        SharedQuery query = sharedQueryRepo.get(sharedQueryDB, queryName);
        String jsonResponse = query.executeUsingGson(sharedQueryDB, jsonBody);
        return ctx.send(jsonResponse);
      } catch (SharedQueryException e) {
        throw new BadRequestException(String.format("{\"error\": \"%s\"}", e.getMessage()));
      }
    });

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
    runApp(args, ExecutionMode.EVENT_LOOP, JoobyNoMVCJdbcApp::new);
  }
}
