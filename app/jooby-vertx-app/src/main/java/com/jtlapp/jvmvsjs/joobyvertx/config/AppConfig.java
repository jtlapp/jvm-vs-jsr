package com.jtlapp.jvmvsjs.joobyvertx.config;

import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.jtlapp.jvmvsjs.joobylib.JoobyServerConfig;
import io.avaje.inject.Component;

/**
 * Configurations that may affect load performance.
 */
@Component
public class AppConfig {
    public JoobyServerConfig server = new JoobyServerConfig();

    public JsonNode toJsonNode(ObjectMapper mapper) {
        return mapper.valueToTree(this);
    }
}
