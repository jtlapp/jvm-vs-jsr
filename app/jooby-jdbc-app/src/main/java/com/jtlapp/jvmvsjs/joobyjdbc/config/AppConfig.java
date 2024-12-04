package com.jtlapp.jvmvsjs.joobyjdbc.config;

import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.jtlapp.jvmvsjs.javalib.AppProperties;
import com.jtlapp.jvmvsjs.joobylib.ServerConfig;

import javax.sql.DataSource;

/**
 * Configurations that may affect load performance.
 */
public class AppConfig {
    private static final int ioThreadCount = Integer.parseInt(
            AppProperties.get("jooby.jetty.ioThreadCount"));
    private static final int workerThreadCount = Integer.parseInt(
            AppProperties.get("jooby.jetty.workerThreadCount"));

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
