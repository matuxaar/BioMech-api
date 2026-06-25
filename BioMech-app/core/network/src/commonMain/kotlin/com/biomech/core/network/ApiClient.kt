package com.biomech.core.network

import io.ktor.client.*
import io.ktor.client.plugins.*
import io.ktor.client.plugins.contentnegotiation.*
import io.ktor.client.plugins.auth.*
import io.ktor.client.plugins.auth.providers.*
import io.ktor.client.request.*
import io.ktor.http.*
import io.ktor.serialization.kotlinx.json.*
import kotlinx.serialization.json.Json

object ApiConfig {
    var baseUrl: String = "http://10.0.2.2:8080/api/v1"
    var token: String? = null

    val healthUrl: String
        get() {
            val clean = baseUrl.trimEnd('/')
            val idx = clean.indexOf("/api")
            return if (idx > 0) clean.substring(0, idx) + "/health" else "$clean/health"
        }
}

fun createHttpClient(): HttpClient {
    return HttpClient {
        install(ContentNegotiation) {
            json(Json {
                ignoreUnknownKeys = true
                isLenient = true
            })
        }

        install(HttpTimeout) {
            requestTimeoutMillis = 30_000
            connectTimeoutMillis = 15_000
        }

        defaultRequest {
            url(ApiConfig.baseUrl)
            contentType(ContentType.Application.Json)
            ApiConfig.token?.let {
                header(HttpHeaders.Authorization, "Bearer $it")
            }
        }
    }
}
