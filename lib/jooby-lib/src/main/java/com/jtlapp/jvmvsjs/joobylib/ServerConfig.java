package com.jtlapp.jvmvsjs.joobylib;

import com.jtlapp.jvmvsjs.javalib.CommonServerConfig;
import io.jooby.Server;
import io.jooby.ServerOptions;

/**
 * Configuration relevant to all Jooby apps.
 */
public class ServerConfig extends CommonServerConfig {

    public final int ioThreadCount;
    public final int workerThreadCount;

    public ServerConfig(String webServer, int ioThreadCount, int workerThreadCount) {
        super(webServer);
        this.ioThreadCount = ioThreadCount;
        this.workerThreadCount = workerThreadCount;
    }

    public void setOptions(Server server) {
        server.setOptions(
                new ServerOptions()
                        .setIoThreads(ioThreadCount)
                        .setWorkerThreads(workerThreadCount)
        );
    }
}
