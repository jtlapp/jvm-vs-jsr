package com.jtlapp.jvmvsjs.springjdbc.controllers;

import org.springframework.beans.factory.annotation.Value;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

@RestController
public class HomeController {

	@Value("${spring.threads.virtual.enabled}")
	private boolean isUsingVirtualThreads;

	@RequestMapping("/")
	String home() {
		var threadKind = isUsingVirtualThreads ? "virtual" : "kernel";
		return "Running spring-jdbc with "+ threadKind +" threads\n";
	}
}
