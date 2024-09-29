package com.jtlapp.jvmvsjs.springasync;

import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.scheduling.annotation.EnableAsync;

@SpringBootApplication
@EnableAsync
public class SpringAsyncApp {

	public static void main(String[] args) {
		SpringApplication.run(SpringAsyncApp.class, args);
	}

}
