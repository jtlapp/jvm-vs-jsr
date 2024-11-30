package com.jtlapp.jvmvsjs.javalib;

/**
 * Configuration relevant to all requests made to the app
 * and common across all Java apps.
 */
public class CommonServerConfig {
    public final String jvmVendor = System.getProperty("java.vm.vendor");
    public final String jvmName = System.getProperty("java.vm.name");
    public final String jvmVersion = System.getProperty("java.vm.version");
    public final int initialMemoryMB;
    public final int maxMemoryMB;
    private final String webServer;

    private static final long initialMemoryBytes = Runtime.getRuntime().totalMemory();
    private static final long maxMemoryBytes = Runtime.getRuntime().maxMemory();

    public CommonServerConfig(String webServer) {
        initialMemoryMB = (int) (initialMemoryBytes / 1024 / 1024);
        maxMemoryMB = (int) (maxMemoryBytes / 1024 / 1024);
        this.webServer = webServer;
    }
}
