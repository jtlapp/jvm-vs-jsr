package com.joelapp.javabenchmarks.hikariquerylibrary.tables;

import java.sql.Connection;
import java.sql.PreparedStatement;
import java.sql.SQLException;

public class OrderItemTable extends Table {
    public static void createTable(Connection connection) throws SQLException {
        String sql = """
                CREATE TABLE order_items (
                    id VARCHAR PRIMARY KEY,
                    order_id VARCHAR REFERENCES orders(id),
                    product_id VARCHAR REFERENCES products(id),
                    quantity INTEGER)
            """;
        try (PreparedStatement statement = connection.prepareStatement(sql)) {
            statement.executeUpdate();
        }
    }

    public static void insertOrderItem(
            Connection connection,
            String orderItemID,
            String orderID,
            String productID,
            int quantity
    ) throws SQLException {
        String sql = "INSERT INTO order_items (id, order_id, product_id, quantity) VALUES (?, ?, ?, ?)";
        try (PreparedStatement statement = connection.prepareStatement(sql)) {
            statement.setObject(1, orderItemID);
            statement.setObject(2, orderID);
            statement.setObject(3, productID);
            statement.setInt(4, quantity);
            statement.executeUpdate();
        }
    }
}