package com.jtlapp.jvmvsjs.springjdbc.config;

import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.ObjectMapper;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Component;

/**
 * Configuration that may affect load performance.
 */
@Component
public class AppConfig {
    public final String jvmVendor = System.getProperty("java.vm.vendor");
    public final String jvmName = System.getProperty("java.vm.name");
    public final String jvmVersion = System.getProperty("java.vm.version");
    public final int initialMemoryMB;
    public final int maxMemoryMB;
    @Value("${server.tomcat.threads.min-spare}")
    public int minWebServerThreads; // Tomcat defaults to 10
    @Value("${server.tomcat.threads.max}")
    public int maxWebServerThreads; // Tomcat defaults to 200
    @Value("${server.tomcat.max-connections}")
    public int maxWebServerConns; // Tomcat defaults to 10000

    private static final long initialMemoryBytes = Runtime.getRuntime().totalMemory();
    private static final long maxMemoryBytes = Runtime.getRuntime().maxMemory();

    public AppConfig() {
        initialMemoryMB = (int) (initialMemoryBytes / 1024 / 1024);
        maxMemoryMB = (int) (maxMemoryBytes / 1024 / 1024);
    }

    public JsonNode toJsonNode(ObjectMapper mapper) {
        return mapper.valueToTree(this);
    }
}
