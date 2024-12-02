package com.jtlapp.jvmvsjs.joobyjdbc.config;

import com.zaxxer.hikari.HikariDataSource;

import javax.sql.DataSource;

class HikariConfig {

    public final int maxPoolSize;
    public final int minIdleConnections;
    public final long connectionTimeout;

    HikariConfig(DataSource dataSource) {
        var hikariDS = (HikariDataSource) dataSource;

        maxPoolSize = hikariDS.getMaximumPoolSize();
        minIdleConnections = hikariDS.getMinimumIdle();
        connectionTimeout = hikariDS.getConnectionTimeout();
    }
}
