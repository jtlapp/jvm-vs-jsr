package com.jtlapp.jvmvsjs.springjdbc.config;

import com.zaxxer.hikari.HikariDataSource;
import org.springframework.stereotype.Component;

import javax.sql.DataSource;

@Component
class JdbcConfig {

    public final int maxPoolSize;
    public final int minIdleConnections;
    public final long connectionTimeout;

    JdbcConfig(DataSource dataSource) {
        var hikariDS = (HikariDataSource) dataSource;

        maxPoolSize = hikariDS.getMaximumPoolSize();
        minIdleConnections = hikariDS.getMinimumIdle();
        connectionTimeout = hikariDS.getConnectionTimeout();
    }
}
