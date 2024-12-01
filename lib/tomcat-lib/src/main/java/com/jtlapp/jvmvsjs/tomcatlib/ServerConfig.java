package com.jtlapp.jvmvsjs.tomcatlib;

import com.jtlapp.jvmvsjs.javalib.CommonServerConfig;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Component;

/**
 * Configuration relevant to all Spring Tomcat apps.
 */
@Component
public class ServerConfig extends CommonServerConfig {
    @Value("${server.tomcat.threads.min-spare}")
    public int minWebServerThreads; // Tomcat defaults to 10
    @Value("${server.tomcat.threads.max}")
    public int maxWebServerThreads; // Tomcat defaults to 200
    @Value("${server.tomcat.max-connections}")
    public int maxWebServerConns; // Tomcat defaults to 10000
    @Value("${server.compression.enabled}")
    public boolean responseCompression; // defaults to false

    public ServerConfig() {
        super("Tomcat");
    }
}
