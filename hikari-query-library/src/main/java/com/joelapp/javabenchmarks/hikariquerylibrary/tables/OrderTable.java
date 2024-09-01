package com.joelapp.javabenchmarks.hikariquerylibrary.tables;

import java.sql.Connection;
import java.sql.PreparedStatement;
import java.sql.SQLException;
import java.time.LocalDateTime;

public class OrderTable extends Table {
    public static void createTable(Connection connection) throws SQLException {
        String sql = """
                CREATE TABLE orders (
                    id VARCHAR PRIMARY KEY,
                    user_id VARCHAR REFERENCES users(id),
                    order_date TIMESTAMP,
                    status VARCHAR)
            """;
        try (PreparedStatement statement = connection.prepareStatement(sql)) {
            statement.executeUpdate();
        }
    }

    public static String createID(String userID, int orderNumber) {
        return userID + "_ORDER_" + orderNumber;
    }

    public static void insertOrder(
            Connection connection,
            String orderID,
            String userID,
            LocalDateTime orderDate,
            String status
    ) throws SQLException {
        String sql = "INSERT INTO orders (id, user_id, order_date, status) VALUES (?, ?, ?, ?)";
        try (PreparedStatement statement = connection.prepareStatement(sql)) {
            statement.setObject(1, orderID);
            statement.setObject(2, userID);
            statement.setObject(3, orderDate);
            statement.setString(4, status);
            statement.executeUpdate();
        }
    }
}