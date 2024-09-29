package com.jtlapp.jvmvsjs.r2dbcquery;

import org.springframework.r2dbc.core.DatabaseClient;
import reactor.core.publisher.Mono;

import java.sql.SQLException;
import java.util.HashMap;

public class SharedQueryRepo {

    private final DatabaseClient db;

    private final HashMap<String, SharedQuery> cache = new HashMap<>();

    public SharedQueryRepo(DatabaseClient db) {
        this.db = db;
    }

    public Mono<SharedQuery> get(String queryName)
            throws SharedQueryException
    {
        var sharedQuery = cache.get(queryName);
        if (sharedQuery == null) {
            try {
                var sharedQueryMono = load(queryName);
                return sharedQueryMono.map(sq -> {
                    cache.put(queryName, sq);
                    return sq;
                });
            }
            catch (SQLException e) {
                throw new SharedQueryException(e.getMessage());
            }
        }
        return Mono.just(sharedQuery);
    }

    private Mono<SharedQuery> load(String queryName)
        throws SharedQueryException, SQLException
    {
        return db.sql("select * from shared_queries where name = ?")
                .bind(1, queryName)
                .map(row -> {
                    String returns = row.get("returns", String.class);
                    return new SharedQuery(
                            queryName,
                            row.get("query", String.class),
                            SharedQuery.ReturnType.fromLabel(returns)
                    );
                })
                .one();
    }
}
