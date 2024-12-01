package com.jtlapp.jvmvsjs.joobyjdbc;

import com.fasterxml.jackson.databind.ObjectMapper;
import com.jtlapp.jvmvsjs.joobyjdbc.config.AppConfig;
import com.jtlapp.jvmvsjs.jdbclib.Database;
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
    private final Database db = createDatabase();

    {
        var objectMapper = new ObjectMapper();
        var server = new JettyServer();

        install(new HikariModule());
        var dataSource = require(DataSource.class);
        appConfig = new AppConfig(dataSource);

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

    private Database createDatabase() {
        return new Database(System.getenv("DATABASE_HOST_NAME"),
                System.getenv("DATABASE_USERNAME"),
                System.getenv("DATABASE_PASSWORD"));
    }

    private String toErrorJson(String endpoint, Throwable e) {
        return String.format("{\"endpoint\": \"%s\", \"error\": \"%s: %s\"}",
                endpoint, e.getClass().getSimpleName(), e.getMessage());
    }

    public static void main(final String[] args) {
        runApp(args, ExecutionMode.EVENT_LOOP, JoobyJdbcApp::new);
    }
}
