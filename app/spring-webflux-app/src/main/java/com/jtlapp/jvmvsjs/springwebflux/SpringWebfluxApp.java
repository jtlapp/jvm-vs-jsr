package com.jtlapp.jvmvsjs.springwebflux;

import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;

@SpringBootApplication
public class SpringWebfluxApp {
	public final static String version = "0.1.0";

	public static void main(String[] args) {
		System.setProperty("application.name", SpringWebfluxApp.class.getSimpleName());
		System.setProperty("application.version", version);
		SpringApplication.run(SpringWebfluxApp.class, args);
	}

}
