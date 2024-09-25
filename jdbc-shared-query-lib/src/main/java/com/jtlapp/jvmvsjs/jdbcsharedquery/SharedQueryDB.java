package com.jtlapp.jvmvsjs.jdbcsharedquery;

import java.sql.Connection;
import java.sql.DriverManager;
import java.sql.SQLException;
import java.util.Properties;

public class SharedQueryDB {
    private final String databaseURL;
    private final Properties jdbcProperties = new Properties();


    public SharedQueryDB(String databaseURL, String username, String password) {
        this.databaseURL = databaseURL;

        jdbcProperties.setProperty("user", username);
        jdbcProperties.setProperty("password", password);

        // TODO: what caching should I do with pgbouncer?
        jdbcProperties.setProperty("prepareThreshold", "1");
        jdbcProperties.setProperty("preparedStatementCacheQueries", "25");
        jdbcProperties.setProperty("preparedStatementCacheSizeMiB", "1");
    }

    Connection openConnection() throws SQLException {
        return DriverManager.getConnection(databaseURL, jdbcProperties);
    }
}
