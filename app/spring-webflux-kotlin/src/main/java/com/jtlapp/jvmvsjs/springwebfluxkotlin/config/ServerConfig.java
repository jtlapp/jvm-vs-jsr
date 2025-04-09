package com.jtlapp.jvmvsjs.springwebfluxkotlin.config;

import com.jtlapp.jvmvsjs.javalib.CommonServerConfig;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Component;

@Component
public class ServerConfig extends CommonServerConfig {
    @Value("${reactor.netty.ioWorkerCount}")
    public int ioWorkerCount;
    @Value("${reactor.netty.pool.maxConnections}")
    public int maxWebServerConns;

    public ServerConfig() {
        super("Netty");
    }
}
