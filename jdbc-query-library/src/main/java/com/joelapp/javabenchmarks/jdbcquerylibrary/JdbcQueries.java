package com.joelapp.javabenchmarks.jdbcquerylibrary;

import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.databind.node.ArrayNode;
import com.fasterxml.jackson.databind.node.ObjectNode;
import com.joelapp.javabenchmarks.jdbcquerylibrary.tables.OrderItemTable;
import com.joelapp.javabenchmarks.jdbcquerylibrary.tables.OrderTable;
import com.joelapp.javabenchmarks.jdbcquerylibrary.tables.ProductTable;
import com.joelapp.javabenchmarks.jdbcquerylibrary.tables.UserTable;

import java.sql.*;
import java.time.LocalDateTime;
import java.util.Properties;

public class JdbcQueries {

    private static final int USER_COUNT = 1000;
    private static final int PRODUCT_COUNT = 700;
    private static final int ORDERS_PER_USER = 3;
    private static final int ITEMS_PER_ORDER = 4;

    private final String jdbcURL;
    private final Properties jdbcProperties = new Properties();

    public JdbcQueries(String jdbcURL, String username, String password) {
        this.jdbcURL = jdbcURL;

        jdbcProperties.setProperty("user", username);
        jdbcProperties.setProperty("password", password);

        // TODO: what caching should I do with pgbouncer?
        jdbcProperties.setProperty("prepareThreshold", "1");
        jdbcProperties.setProperty("preparedStatementCacheQueries", "25");
        jdbcProperties.setProperty("preparedStatementCacheSizeMiB", "1");
    }

    public void createTables() throws SQLException {
        try (Connection conn = getConnection()) {
            UserTable.createTable(conn);
            ProductTable.createTable(conn);
            OrderTable.createTable(conn);
            OrderItemTable.createTable(conn);
        }
    }

    public void dropTables() throws SQLException {
        try (Connection conn = getConnection()) {
            String sql = "select tablename from pg_tables where schemaname='public'";
            try (PreparedStatement statement = conn.prepareStatement(sql)) {
                try (var resultSet = statement.executeQuery()) {
                    while (resultSet.next()) {
                        dropTable(conn, resultSet.getString("tablename"));
                    }
                }
            }
        }
    }

    public void populateDatabase() throws SQLException {
        try (Connection conn = getConnection()) {

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

        try (Connection conn = getConnection()) {

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

    public int updateQuery(int userNumber, int orderNumber) throws SQLException {
        String userID = UserTable.createID(userNumber);
        String orderID = OrderTable.createID(userID, orderNumber);

        try (Connection conn = getConnection()) {

            var sql = """
                UPDATE order_items oi
                    SET quantity = quantity + 1
                    FROM orders o
                    WHERE oi.order_id = o.id AND o.id = ?;
                """;

            try (var statement = conn.prepareStatement(sql)) {
                statement.setString(1, orderID);
                return statement.executeUpdate();
            }
        }
    }

    private Connection getConnection() throws SQLException {
        return DriverManager.getConnection(jdbcURL, jdbcProperties);
    }

    private void dropTable(Connection conn, String tableName)
        throws SQLException
    {
        String sql = "drop table if exists \"" + tableName + "\" cascade";
        try (var statement = conn.prepareStatement(sql)) {
            statement.execute();
        }
    }
}
