package com.jtlapp.jvmvsjs.springjdbc.config;

import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.jtlapp.jvmvsjs.javalib.CommonServerConfig;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Component;

/**
 * Configurations that may affect load performance.
 */
@Component
public class AppConfig {
    @Autowired
    public ServerConfig server;

    public JsonNode toJsonNode(ObjectMapper mapper) {
        return mapper.valueToTree(this);
    }
}

/**
 * Configuration relevant to all requests made to the app.
 */
@Component
class ServerConfig extends CommonServerConfig {
    @Value("${server.tomcat.threads.min-spare}")
    public int minWebServerThreads; // Tomcat defaults to 10
    @Value("${server.tomcat.threads.max}")
    public int maxWebServerThreads; // Tomcat defaults to 200
    @Value("${server.tomcat.max-connections}")
    public int maxWebServerConns; // Tomcat defaults to 10000
}
