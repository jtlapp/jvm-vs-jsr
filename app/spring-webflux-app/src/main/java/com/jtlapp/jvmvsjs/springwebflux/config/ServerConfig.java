package com.jtlapp.jvmvsjs.springwebflux.config;

import com.jtlapp.jvmvsjs.javalib.CommonServerConfig;
import org.springframework.stereotype.Component;

// Netty doesn't appear to be much configurable within Spring Boot.

@Component
public class ServerConfig extends CommonServerConfig {
    public ServerConfig() {
        super("Netty");
    }
}
