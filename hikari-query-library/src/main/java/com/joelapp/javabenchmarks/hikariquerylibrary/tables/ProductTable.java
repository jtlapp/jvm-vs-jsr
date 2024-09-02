package com.joelapp.javabenchmarks.hikariquerylibrary.tables;

import java.sql.Connection;
import java.sql.PreparedStatement;
import java.sql.SQLException;
import java.time.LocalDateTime;

public class ProductTable extends Table {
    public static void createTable(Connection connection) throws SQLException {
        String sql = """
            CREATE TABLE products (
                id VARCHAR PRIMARY KEY,
                name VARCHAR,
                description TEXT,
                price NUMERIC,
                stock_quantity INTEGER,
                created_at TIMESTAMP DEFAULT NOW())
            """;
        try (PreparedStatement statement = connection.prepareStatement(sql)) {
            statement.executeUpdate();
        }
    }

    public static String createID(int productNumber) {
        return createID("PRODUCT_", productNumber);
    }

    public static void insertProduct(
            Connection connection,
            String productID,
            String name,
            String description,
            java.math.BigDecimal price,
            int stockQuantity
    ) throws SQLException {
        String sql = """
            INSERT INTO products (id, name, description, price, stock_quantity, created_at)
                VALUES (?, ?, ?, ?, ?, ?)
            """;
        try (PreparedStatement statement = connection.prepareStatement(sql)) {
            statement.setObject(1, productID);
            statement.setString(2, name);
            statement.setString(3, description);
            statement.setBigDecimal(4, price);
            statement.setInt(5, stockQuantity);
            statement.setObject(6, LocalDateTime.now());
            statement.executeUpdate();
        }
    }
}