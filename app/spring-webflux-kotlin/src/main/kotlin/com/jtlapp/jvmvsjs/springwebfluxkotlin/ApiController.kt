package com.jtlapp.jvmvsjs.springwebfluxkotlin

import com.fasterxml.jackson.databind.JsonNode
import com.fasterxml.jackson.databind.ObjectMapper
import com.jtlapp.jvmvsjs.r2dbclib.Database
import com.jtlapp.jvmvsjs.springwebfluxkotlin.config.AppConfig
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.delay
import kotlinx.coroutines.withContext
import org.springframework.beans.factory.annotation.Autowired
import org.springframework.http.HttpStatus
import org.springframework.http.ResponseEntity
import org.springframework.web.bind.annotation.*

@RestController
@RequestMapping("/api")
class ApiController @Autowired constructor(
    private val appConfig: AppConfig,
    private val db: Database
) {

    companion object {
        private val APP_NAME = System.getenv("APP_NAME") ?: ""
        private const val APP_VERSION = "0.1.0"
        private val objectMapper = ObjectMapper()
    }

    @GetMapping("/info")
    suspend fun info(): String {
        val jsonObj: JsonNode = objectMapper.createObjectNode()
            .put("appName", APP_NAME)
            .put("appVersion", APP_VERSION)
            .set("appConfig", appConfig.toJsonNode(objectMapper))
        return jsonObj.toString()
    }

    @GetMapping("/app-sleep")
    suspend fun appSleep(@RequestParam("millis") millis: Int): String {
        delay(millis.toLong())
        return ""
    }

    @GetMapping("/pg-sleep")
    suspend fun pgSleep(@RequestParam("millis") millis: Int): ResponseEntity<String> {
        return withContext(Dispatchers.IO) {
            try {
                db.issueSleepQuery(millis).block()
                ResponseEntity.ok().body("{}")
            } catch (e: Exception) {
                ResponseEntity
                    .status(HttpStatus.INTERNAL_SERVER_ERROR)
                    .body(toErrorJson("pg-sleep", e))
            }
        }
    }

    private fun toErrorJson(endpoint: String, e: Throwable): String {
        return """{"endpoint": "$endpoint", "error": "${e.javaClass.simpleName}: ${e.message}"}"""
    }
}
