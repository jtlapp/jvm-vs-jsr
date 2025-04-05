package com.jtlapp.jvmvsjs.javalib;

import java.io.IOException;
import java.io.InputStream;
import java.util.Properties;

public class AppProperties {
    private static final Properties properties = new Properties();

    public static void init(ClassLoader classLoader) {
        try (InputStream inputStream =
                     classLoader.getResourceAsStream("application.properties")) {
            if (inputStream == null) {
                throw new RuntimeException("Unable to find application.properties");
            }
            properties.load(inputStream);
        } catch (IOException e) {
            throw new RuntimeException("Error loading application.properties", e);
        }
    }

    public static String get(String key) {
        return properties.getProperty(key);
    }
}