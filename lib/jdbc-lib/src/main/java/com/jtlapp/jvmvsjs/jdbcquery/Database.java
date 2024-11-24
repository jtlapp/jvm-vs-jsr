package com.jtlapp.jvmvsjs.jdbcquery;

import java.sql.Connection;
import java.sql.DriverManager;
import java.sql.SQLException;
import java.util.Properties;

public class Database {
    private final String databaseURL;
    private final Properties jdbcProperties = new Properties();

    public Database(String databaseURL, String username, String password) {
        this.databaseURL = databaseURL;

        jdbcProperties.setProperty("user", username);
        jdbcProperties.setProperty("password", password);

        // TODO: what caching should I do with pgbouncer?
        jdbcProperties.setProperty("prepareThreshold", "1");
        jdbcProperties.setProperty("preparedStatementCacheQueries", "25");
        jdbcProperties.setProperty("preparedStatementCacheSizeMiB", "1");
    }

    public Connection openConnection() throws SQLException {
        return DriverManager.getConnection(databaseURL, jdbcProperties);
    }

    public static void issueSleepQuery(Connection conn, int durationMillis) throws SQLException {
        var statement = conn.prepareStatement("SELECT pg_sleep(?)");
        statement.setDouble(1, durationMillis / 1000.0);
        statement.executeQuery();
    }
}
