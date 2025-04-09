package com.jtlapp.jvmvsjs.springwebflux;

import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;

@SpringBootApplication(scanBasePackages = {
		"com.jtlapp.jvmvsjs.springwebflux", "com.jtlapp.jvmvsjs.webfluxlib"
})
public class WebApp {

	public static void main(String[] args) {
		SpringApplication.run(WebApp.class, args);
	}

}
