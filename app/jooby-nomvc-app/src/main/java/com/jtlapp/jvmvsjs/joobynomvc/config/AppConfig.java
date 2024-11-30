package com.jtlapp.jvmvsjs.joobynomvc.config;

import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.jtlapp.jvmvsjs.joobylib.JoobyServerConfig;

/**
 * Configurations that may affect load performance.
 */
public class AppConfig {
    public JoobyServerConfig server = new JoobyServerConfig();

    public JsonNode toJsonNode(ObjectMapper mapper) {
        return mapper.valueToTree(this);
    }
}
