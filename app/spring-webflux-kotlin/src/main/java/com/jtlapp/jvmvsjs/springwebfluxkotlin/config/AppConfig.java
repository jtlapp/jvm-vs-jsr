package com.jtlapp.jvmvsjs.springwebfluxkotlin.config;

import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.jtlapp.jvmvsjs.javalib.CommonServerConfig;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Component;

/**
 * Configurations that may affect load performance.
 */
@Component
public class AppConfig {
    @Autowired
    public ServerConfig server;

    @Autowired
    public R2dbcConfig dbclient;

    public JsonNode toJsonNode(ObjectMapper mapper) {
        return mapper.valueToTree(this);
    }
}
