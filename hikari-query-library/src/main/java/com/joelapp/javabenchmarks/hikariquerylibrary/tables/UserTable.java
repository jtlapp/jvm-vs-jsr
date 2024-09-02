package com.joelapp.javabenchmarks.hikariquerylibrary.tables;


import java.sql.Connection;
import java.sql.PreparedStatement;
import java.sql.SQLException;
import java.time.LocalDateTime;

public class UserTable extends Table {
    public static void createTable(Connection connection) throws SQLException {
        String sql = """
            CREATE TABLE users (
                id VARCHAR PRIMARY KEY,
                username VARCHAR UNIQUE NOT NULL,
                email VARCHAR UNIQUE NOT NULL,
                created_at TIMESTAMP DEFAULT NOW())
            """;
        try (PreparedStatement statement = connection.prepareStatement(sql)) {
            statement.executeUpdate();
        }
    }

    public static String createID(int userNumber) {
        return createID("USER_", userNumber);
    }

    public static void insertUser(
            Connection connection,
            String userID,
            String username,
            String email
    ) throws SQLException {
        String sql = """
            INSERT INTO users (id, username, email, created_at) VALUES (?, ?, ?, ?)
            """;
        try (PreparedStatement statement = connection.prepareStatement(sql)) {
            statement.setObject(1, userID);
            statement.setString(2, username);
            statement.setString(3, email);
            statement.setObject(4, LocalDateTime.now());
            statement.executeUpdate();
        }
    }
}