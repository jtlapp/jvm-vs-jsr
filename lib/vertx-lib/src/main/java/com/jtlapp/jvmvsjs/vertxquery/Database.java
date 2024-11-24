package com.jtlapp.jvmvsjs.vertxquery;

import io.vertx.core.Future;
import io.vertx.sqlclient.Pool;
import io.vertx.sqlclient.Tuple;

public class Database {

    private final Pool db;

    public Database(Pool db) {
        this.db = db;
    }

    public Future<Void> issueSleepQuery(int durationMillis) {
        String sql = "SELECT pg_sleep($1)";
        Tuple params = Tuple.of(durationMillis / 1000.0);

        return db.preparedQuery(sql).execute(params).mapEmpty();
    }
}
