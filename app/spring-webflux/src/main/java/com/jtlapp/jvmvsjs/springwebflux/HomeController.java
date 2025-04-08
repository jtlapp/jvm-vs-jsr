package com.jtlapp.jvmvsjs.springwebflux;

import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

@RestController
public class HomeController {

	@RequestMapping("/")
	String home() {
		return "Running spring-webflux\n";
	}
}
