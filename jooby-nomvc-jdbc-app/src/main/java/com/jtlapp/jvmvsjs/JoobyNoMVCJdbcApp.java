package com.jtlapp.jvmvsjs;

import io.jooby.Jooby;
import io.jooby.exception.BadRequestException;
import io.jooby.netty.NettyServer;

public class JoobyNoMVCJdbcApp extends Jooby {

  {
    install(new NettyServer());

    SharedQueryDB sharedQueryDB = new SharedQueryDB(
            System.getenv("DATABASE_URL"),
            System.getenv("DATABASE_USERNAME"),
            System.getenv("DATABASE_PASSWORD"));
    SharedQueryRepo sharedQueryRepo = new SharedQueryRepo();

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
  }

  public static void main(final String[] args) {
    runApp(args, JoobyNoMVCJdbcApp::new);
  }
}
