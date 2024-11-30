package com.jtlapp.jvmvsjs.joobylib;

import com.jtlapp.jvmvsjs.javalib.CommonServerConfig;
import io.jooby.ServerOptions;
import io.jooby.netty.NettyServer;

/**
 * Configuration relevant to all Jooby apps.
 */
public class JoobyServerConfig extends CommonServerConfig {

    public int ioThreadCount = Math.max(Runtime.getRuntime().availableProcessors(), 2) * 8;
    public int workerThreadCount = ioThreadCount;

    public JoobyServerConfig() {
        super("Netty");
    }

    public void setOptions(NettyServer server) {
        server.setOptions(
                new ServerOptions()
                        .setIoThreads(ioThreadCount)
                        .setWorkerThreads(workerThreadCount)
        );
    }
}
