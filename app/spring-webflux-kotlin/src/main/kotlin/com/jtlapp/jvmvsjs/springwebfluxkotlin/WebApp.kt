package com.jtlapp.jvmvsjs.springwebfluxkotlin

import org.springframework.boot.autoconfigure.SpringBootApplication
import org.springframework.boot.runApplication

@SpringBootApplication(scanBasePackages = [
	"com.jtlapp.jvmvsjs.springwebflux", "com.jtlapp.jvmvsjs.webfluxlib"
])
class WebApp

fun main(args: Array<String>) {
	runApplication<WebApp>(*args)
}
