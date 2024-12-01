package com.jtlapp.jvmvsjs.joobyjdbc.config;

import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.jtlapp.jvmvsjs.joobylib.ServerConfig;

import javax.sql.DataSource;

/**
 * Configurations that may affect load performance.
 */
public class AppConfig {
    private static final int ioThreadCount = Integer.parseInt(
            System.getProperty("jooby.jetty.ioThreadCount", "200"));
    private static final int workerThreadCount = Integer.parseInt(
            System.getProperty("jooby.jetty.workerThreadCount", "200"));

    public final ServerConfig server =
            new ServerConfig("Jetty", ioThreadCount, workerThreadCount);
    public final JdbcConfig dbclient;

    public AppConfig(DataSource dataSource) {
        dbclient = new JdbcConfig(dataSource);
    }

    public JsonNode toJsonNode(ObjectMapper mapper) {
        return mapper.valueToTree(this);
    }
}
