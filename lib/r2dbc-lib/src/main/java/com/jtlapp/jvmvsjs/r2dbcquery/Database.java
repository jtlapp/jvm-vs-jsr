package com.jtlapp.jvmvsjs.r2dbcquery;

import org.springframework.r2dbc.core.DatabaseClient;
import reactor.core.publisher.Mono;

public class Database {

    private final DatabaseClient db;

    public Database(DatabaseClient db) {
        this.db = db;
    }

    public Mono<Void> issueSleepQuery(int durationMillis) {
        return db.sql("SELECT pg_sleep(:duration)")
                .bind("duration", durationMillis / 1000.0)
                .fetch()
                .rowsUpdated()
                .then();
    }
}
