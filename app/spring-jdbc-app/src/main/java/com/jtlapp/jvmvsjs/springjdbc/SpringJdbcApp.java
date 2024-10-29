package com.jtlapp.jvmvsjs.springjdbc;

import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;

@SpringBootApplication
public class SpringJdbcApp {
	public final static String version = "0.1.0";

	public static void main(String[] args) {
		System.setProperty("application.name", SpringJdbcApp.class.getSimpleName());
		System.setProperty("application.version", version);
		SpringApplication.run(SpringJdbcApp.class, args);
	}

}
