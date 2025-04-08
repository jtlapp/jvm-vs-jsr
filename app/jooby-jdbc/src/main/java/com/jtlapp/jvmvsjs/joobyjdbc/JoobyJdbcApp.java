package com.jtlapp.jvmvsjs.joobyjdbc;

import com.fasterxml.jackson.databind.ObjectMapper;
import com.jtlapp.jvmvsjs.javalib.AppProperties;
import com.jtlapp.jvmvsjs.joobyjdbc.config.AppConfig;
import com.jtlapp.jvmvsjs.hikarilib.Database;
import io.jooby.ExecutionMode;
import io.jooby.Jooby;
import io.jooby.StatusCode;
import io.jooby.exception.StatusCodeException;
import io.jooby.hikari.HikariModule;
import io.jooby.jetty.JettyServer;

import java.sql.SQLException;

public class JoobyJdbcApp extends Jooby {

    public final String APP_NAME = System.getenv("APP_NAME");
    public final String APP_VERSION = "0.1.0";

    private final AppConfig appConfig;
    private final Database db;

    {
        AppProperties.init(JoobyJdbcApp.class.getClassLoader());
        var objectMapper = new ObjectMapper();

        // HikariCP uses JDBC under the hood.
        var dataSource = AppConfig.createDataSource();
        install(new HikariModule(dataSource));
        db = new Database(dataSource);
        appConfig = new AppConfig(dataSource);

        var server = new JettyServer();
        appConfig.server.setOptions(server);
        install(server);

        get("/", ctx -> "Running Jooby with Jetty and JDBC\n");

        get("/api/info", ctx -> {
            var jsonObj = objectMapper.createObjectNode()
                    .put("appName", APP_NAME)
                    .put("appVersion", APP_VERSION)
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
            try(var conn = db.openConnection()) {
                Database.issueSleepQuery(conn, millis);
                return "{}";
            }
            catch (SQLException e) {
                throw new StatusCodeException(StatusCode.BAD_REQUEST,
                        toErrorJson("pg-sleep", e));
            }
        });
    }

    private String toErrorJson(String endpoint, Throwable e) {
        return String.format("{\"endpoint\": \"%s\", \"error\": \"%s: %s\"}",
                endpoint, e.getClass().getSimpleName(), e.getMessage());
    }

    public static void main(final String[] args) {
        runApp(args, ExecutionMode.EVENT_LOOP, JoobyJdbcApp::new);
    }
}
