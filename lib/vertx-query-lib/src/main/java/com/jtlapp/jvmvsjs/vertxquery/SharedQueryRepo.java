package com.jtlapp.jvmvsjs.vertxquery;

import io.vertx.core.Future;
import io.vertx.pgclient.PgPool;
import io.vertx.sqlclient.Tuple;
import io.vertx.sqlclient.Row;

import java.util.HashMap;

public class SharedQueryRepo {

    private final PgPool db;
    private final HashMap<String, SharedQuery> cache = new HashMap<>();

    public SharedQueryRepo(PgPool db) {
        this.db = db;
    }

    public Future<SharedQuery> get(String queryName) {
        var sharedQuery = cache.get(queryName);
        if (sharedQuery == null) {
            return load(queryName).onSuccess(sq -> cache.put(queryName, sq));
        }
        return Future.succeededFuture(sharedQuery);
    }

    private Future<SharedQuery> load(String queryName) {
        String sql = "SELECT * FROM shared_queries WHERE name = $1";
        Tuple params = Tuple.of(queryName);

        return db.preparedQuery(sql).execute(params).flatMap(rows -> {
            Row row = rows.iterator().next();

            String returns = row.getString("returns");
            SharedQuery sharedQuery = new SharedQuery(
                    queryName,
                    row.getString("query"),
                    SharedQuery.ReturnType.fromLabel(returns)
            );
            return Future.succeededFuture(sharedQuery);
        }).recover(throwable -> {
            var exception = new SharedQueryException(throwable.getMessage());
            return Future.failedFuture(exception);
        });
    }
}
