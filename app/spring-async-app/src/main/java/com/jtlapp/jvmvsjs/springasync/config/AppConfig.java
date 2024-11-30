package com.jtlapp.jvmvsjs.springasync.config;

import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.jtlapp.jvmvsjs.tomcatlib.TomcatServerConfig;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Component;

/**
 * Configurations that may affect load performance.
 */
@Component
public class AppConfig {
    @Autowired
    public TomcatServerConfig server;

    public JsonNode toJsonNode(ObjectMapper mapper) {
        return mapper.valueToTree(this);
    }
}
