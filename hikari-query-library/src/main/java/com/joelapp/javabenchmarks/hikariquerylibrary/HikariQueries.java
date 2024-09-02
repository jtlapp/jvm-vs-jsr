package com.joelapp.javabenchmarks.hikariquerylibrary;

import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.databind.node.ArrayNode;
import com.fasterxml.jackson.databind.node.ObjectNode;
import com.joelapp.javabenchmarks.hikariquerylibrary.tables.OrderItemTable;
import com.joelapp.javabenchmarks.hikariquerylibrary.tables.OrderTable;
import com.joelapp.javabenchmarks.hikariquerylibrary.tables.ProductTable;
import com.joelapp.javabenchmarks.hikariquerylibrary.tables.UserTable;
import com.zaxxer.hikari.HikariConfig;
import com.zaxxer.hikari.HikariDataSource;

import java.sql.Connection;
import java.sql.SQLException;
import java.sql.Timestamp;
import java.time.LocalDateTime;

public class HikariQueries {

    private static final int USER_COUNT = 1000;
    private static final int PRODUCT_COUNT = 700;
    private static final int ORDERS_PER_USER = 3;
    private static final int ITEMS_PER_ORDER = 4;

    private HikariDataSource dataSource;

    public HikariQueries(String jdbcUrl, String username, String password) {
        HikariConfig config = new HikariConfig();

        config.setJdbcUrl(jdbcUrl);
        config.setUsername(username);
        config.setPassword(password);
        config.setMaximumPoolSize(10);

        config.addDataSourceProperty("prepareThreshold", "1");
        config.addDataSourceProperty("preparedStatementCacheQueries", "25");
        config.addDataSourceProperty("preparedStatementCacheSizeMiB", "1");

        dataSource = new HikariDataSource(config);
    }

    public void createTables() throws SQLException {
        try (Connection conn = dataSource.getConnection()) {
            UserTable.createTable(conn);
            ProductTable.createTable(conn);
            OrderTable.createTable(conn);
            OrderItemTable.createTable(conn);
        }
    }

    public void populateDatabase() throws SQLException {
        try (Connection conn = dataSource.getConnection()) {

            for (int i = 1; i <= USER_COUNT; i++) {
                UserTable.insertUser(conn, UserTable.createID(i),
                        "user" + i, "user" + i + "@example.com");
            }

            for (int i = 1; i <= PRODUCT_COUNT; i++) {
                ProductTable.insertProduct(conn, ProductTable.createID(i),
                        "Product " + i, "Description of product " + i,
                        new java.math.BigDecimal((i % 50) + ".99"), 100);
            }

            int orderedItemCount = 0;
            for (int i = 1; i <= USER_COUNT; i++) {
                for (int j = 1; j <= ORDERS_PER_USER; ++j) {
                    var userID = UserTable.createID(i);
                    var orderID = OrderTable.createID(userID, j);
                    OrderTable.insertOrder(conn, orderID, userID, LocalDateTime.now(), "Shipped");

                    for (int k = 1; k <= ITEMS_PER_ORDER; ++k) {
                        var orderItemID = orderID + "_ITEM_" + k;
                        var productNumber = (orderedItemCount % PRODUCT_COUNT) + 1;
                        var productID = ProductTable.createID(productNumber);
                        OrderItemTable.insertOrderItem(conn, orderItemID, orderID, productID, 1);
                        ++orderedItemCount;
                    }
                }
            }
        }
    }

    public String selectQuery(int userNumber, int orderNumber) throws SQLException {
        String userID = UserTable.createID(userNumber);
        String orderID = OrderTable.createID(userID, orderNumber);

        ObjectMapper objectMapper = new ObjectMapper();
        ArrayNode arrayNode = objectMapper.createArrayNode();

        try (Connection conn = dataSource.getConnection()) {

            var sql = """
                SELECT o.id AS order_id, o.order_date, o.status, u.username, u.email,
                        p.name, p.description, p.price, oi.quantity
                    FROM orders o
                    JOIN users u ON o.user_id = u.id
                    JOIN order_items oi ON oi.order_id = o.id
                    JOIN products p ON oi.product_id = p.id
                    WHERE o.id = ?;
                """;

            try (var statement = conn.prepareStatement(sql)) {
                statement.setString(1, orderID);
                try (var resultSet = statement.executeQuery()) {
                    while (resultSet.next()) {
                        ObjectNode jsonObject = objectMapper.createObjectNode();

                        jsonObject.put("orderID", resultSet.getString("order_id"));
                        Timestamp orderTimestamp = resultSet.getTimestamp("order_date");
                        jsonObject.put("orderDate", orderTimestamp.toLocalDateTime().toString());
                        jsonObject.put("status", resultSet.getString("status"));
                        jsonObject.put("username", resultSet.getString("username"));
                        jsonObject.put("email", resultSet.getString("email"));
                        jsonObject.put("productName", resultSet.getString("name"));
                        jsonObject.put("productDescription", resultSet.getString("description"));
                        jsonObject.put("price", resultSet.getBigDecimal("price"));
                        jsonObject.put("quantity", resultSet.getInt("quantity"));

                        arrayNode.add(jsonObject);
                    }
                }
            }
        }

        try {
            return objectMapper.writeValueAsString(arrayNode);
        } catch (JsonProcessingException e) {
            throw new RuntimeException(e);
        }
    }

    public void updateQuery(int userNumber, int orderNumber) throws SQLException {
        String userID = UserTable.createID(userNumber);
        String orderID = OrderTable.createID(userID, orderNumber);

        try (Connection conn = dataSource.getConnection()) {

            var sql = """
                UPDATE order_items oi
                    SET quantity = quantity + 1
                    FROM orders o
                    WHERE oi.order_id = o.id AND o.id = ?;
                """;

            try (var statement = conn.prepareStatement(sql)) {
                statement.setString(1, orderID);
                statement.executeQuery();
            }
        }
    }

    public void close() {
        if (dataSource != null) {
            dataSource.close();
        }
    }
}
