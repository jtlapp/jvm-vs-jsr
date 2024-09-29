package com.jtlapp.jvmvsjs.jdbcquery;

import java.sql.Connection;
import java.sql.SQLException;
import java.util.HashMap;

public class SharedQueryRepo {

    private final HashMap<String, SharedQuery> cache = new HashMap<>();

    public SharedQuery get(SharedQueryDB db, String queryName)
            throws SharedQueryException
    {
        var sharedQuery = cache.get(queryName);
        if (sharedQuery == null) {
            try {
                sharedQuery = load(db, queryName);
                cache.put(queryName, sharedQuery);
            }
            catch (SQLException e) {
                throw new SharedQueryException(e.getMessage());
            }
        }
        return sharedQuery;
    }

    private SharedQuery load(SharedQueryDB db, String queryName)
        throws SharedQueryException, SQLException
    {
        try (Connection conn = db.openConnection()) {

            var sql = "select * from shared_queries where name = ?";

            try (var statement = conn.prepareStatement(sql)) {
                statement.setString(1, queryName);
                var resultSet = statement.executeQuery();
                if (resultSet.next()) {
                    String returns = resultSet.getString("returns");
                    return new SharedQuery(
                            queryName,
                            resultSet.getString("query"),
                            SharedQuery.ReturnType.fromLabel(returns)
                    );
                } else {
                    throw new SharedQueryException(String.format(
                            "Shared query '%s' not found", queryName));
                }
            }
        }
    }

}
