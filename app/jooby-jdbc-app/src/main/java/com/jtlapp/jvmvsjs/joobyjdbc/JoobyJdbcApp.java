package com.jtlapp.jvmvsjs.joobyjdbc;

import com.fasterxml.jackson.databind.ObjectMapper;
import com.jtlapp.jvmvsjs.joobyjdbc.config.AppConfig;
import com.jtlapp.jvmvsjs.hikarilib.Database;
import com.zaxxer.hikari.HikariConfig;
import io.jooby.ExecutionMode;
import io.jooby.Jooby;
import io.jooby.StatusCode;
import io.jooby.exception.StatusCodeException;
import io.jooby.hikari.HikariModule;
import io.jooby.jetty.JettyServer;

import javax.sql.DataSource;
import java.sql.SQLException;

public class JoobyJdbcApp extends Jooby {
    public final String appName = System.getenv("APP_NAME");
    public final String appVersion = System.getenv("APP_VERSION");

    private final AppConfig appConfig;
    private final Database db;

    {
        var objectMapper = new ObjectMapper();

        // HikariCP uses JDBC under the hood.
        var hikariConfig = createHikariConfig();
        install(new HikariModule(hikariConfig));
        var dataSource = require(DataSource.class);
        db = new Database(dataSource);
        appConfig = new AppConfig(dataSource);

        var server = new JettyServer();
        appConfig.server.setOptions(server);
        install(server);

        get("/", ctx -> "Running Jooby with Jetty and JDBC\n");

        get("/api/info", ctx -> {
            var jsonObj = objectMapper.createObjectNode()
                    .put("appName", appName)
                    .put("appVersion", appVersion)
                    .set("appConfig", appConfig.toJsonNode(objectMapper));
            return jsonObj.toString();
        });

        get("/api/app-sleep", ctx -> {
            int millis = ctx.query("millis").intValue(0);
            Thread.sleep(millis);
            return "{}";
        });

        get("/api/pg-sleep", ctx -> {
            int millis = ctx.query("millis").intValue(0);
            try {
                // TODO: Do I need to close this connection, here and elsewhere?
                var conn = db.openConnection();
                Database.issueSleepQuery(conn, millis);
                return "{}";
            }
            catch (SQLException e) {
                throw new StatusCodeException(StatusCode.BAD_REQUEST,
                        toErrorJson("pg-sleep", e));
            }
        });
    }

    private HikariConfig createHikariConfig() {
        var config = new HikariConfig();
        config.setJdbcUrl("jdbc:" + System.getenv("DATABASE_URL"));
        config.setUsername(System.getenv("DATABASE_USERNAME"));
        config.setPassword(System.getenv("DATABASE_PASSWORD"));
        return config;
    }

    private String toErrorJson(String endpoint, Throwable e) {
        return String.format("{\"endpoint\": \"%s\", \"error\": \"%s: %s\"}",
                endpoint, e.getClass().getSimpleName(), e.getMessage());
    }

    public static void main(final String[] args) {
        runApp(args, ExecutionMode.EVENT_LOOP, JoobyJdbcApp::new);
    }
}
