package com.jtlapp.jvmvsjs.joobyr2dbc.config;

import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.jtlapp.jvmvsjs.joobylib.ServerConfig;

/**
 * Configurations that may affect load performance.
 */
public class AppConfig {

    private static final int threadCount =
            Math.max(Runtime.getRuntime().availableProcessors(), 2) * 8;

    public final ServerConfig server =
            new ServerConfig("Netty", threadCount, threadCount);
    public final R2dbcConfig dbclient = new R2dbcConfig();

    public JsonNode toJsonNode(ObjectMapper mapper) {
        return mapper.valueToTree(this);
    }
}
