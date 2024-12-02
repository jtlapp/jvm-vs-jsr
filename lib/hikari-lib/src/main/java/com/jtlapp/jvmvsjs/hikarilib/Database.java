package com.jtlapp.jvmvsjs.hikarilib;

import javax.sql.DataSource;
import java.sql.Connection;
import java.sql.SQLException;

public class Database {
    private final DataSource dataSource;

    public Database(DataSource dataSource) {
        this.dataSource = dataSource;
    }

    public Connection openConnection() throws SQLException {
        return dataSource.getConnection();
    }

    public static void issueSleepQuery(Connection conn, int durationMillis) throws SQLException {
        var statement = conn.prepareStatement("SELECT pg_sleep(?)");
        statement.setDouble(1, durationMillis / 1000.0);
        statement.executeQuery();
    }
}
