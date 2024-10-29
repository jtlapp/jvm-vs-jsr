package com.jtlapp.jvmvsjs.springasync;

import org.springframework.beans.factory.annotation.Value;
import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.scheduling.annotation.EnableAsync;

@SpringBootApplication
@EnableAsync
public class SpringAsyncApp {
	public final static String version = "0.1.0";

	public static void main(String[] args) {
		System.setProperty("application.name", SpringAsyncApp.class.getSimpleName());
		System.setProperty("application.version", version);
		SpringApplication.run(SpringAsyncApp.class, args);
	}

}
