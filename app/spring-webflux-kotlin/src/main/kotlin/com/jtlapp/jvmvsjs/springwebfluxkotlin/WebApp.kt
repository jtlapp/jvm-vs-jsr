package com.jtlapp.jvmvsjs.springwebfluxkotlin

import org.springframework.boot.autoconfigure.SpringBootApplication
import org.springframework.boot.runApplication

@SpringBootApplication
class WebApp

fun main(args: Array<String>) {
	runApplication<WebApp>(*args)
}
