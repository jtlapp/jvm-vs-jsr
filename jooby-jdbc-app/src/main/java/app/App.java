package app;

import io.jooby.Jooby;
import io.jooby.netty.NettyServer;

public class App extends Jooby {

  {
    install(new NettyServer());
    get("/", ctx -> "Welcome to Jooby!");
  }

  public static void main(final String[] args) {
    runApp(args, App::new);
  }

}
