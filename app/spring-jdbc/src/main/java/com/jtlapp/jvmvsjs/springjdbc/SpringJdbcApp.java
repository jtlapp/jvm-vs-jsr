package com.jtlapp.jvmvsjs.springjdbc;

import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.context.annotation.ComponentScan;

@SpringBootApplication
@ComponentScan(basePackages = {
	"com.jtlapp.jvmvsjs.springjdbc",
	"com.jtlapp.jvmvsjs.tomcatlib"
})
public class SpringJdbcApp {

	public static void main(String[] args) {
		SpringApplication.run(SpringJdbcApp.class, args);
	}

}
