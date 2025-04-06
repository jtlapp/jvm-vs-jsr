package com.jtlapp.jvmvsjs.joobyjdbc.config;

import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.jtlapp.jvmvsjs.javalib.AppProperties;
import com.jtlapp.jvmvsjs.joobylib.ServerConfig;
import com.zaxxer.hikari.HikariConfig;
import com.zaxxer.hikari.HikariDataSource;

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

    public static HikariDataSource createDataSource() {
        var maximumPoolSize = Integer.parseInt(
                AppProperties.get("jooby.hikari.maximumPoolSize"));
        var minimumIdle = Integer.parseInt(
                AppProperties.get("jooby.hikari.minimumIdle"));
        var connectionTimeoutMillis = Integer.parseInt(
                AppProperties.get("jooby.hikari.connectionTimeout"));

        var config = new HikariConfig();
        config.setJdbcUrl("jdbc:" + System.getenv("DATABASE_URL"));
        config.setUsername(System.getenv("DATABASE_USERNAME"));
        config.setPassword(System.getenv("DATABASE_PASSWORD"));
        config.setMaximumPoolSize(maximumPoolSize);
        config.setMinimumIdle(minimumIdle);
        config.setConnectionTimeout(connectionTimeoutMillis);

        return new HikariDataSource(config);
    }

    public JsonNode toJsonNode(ObjectMapper mapper) {
        return mapper.valueToTree(this);
    }
}
