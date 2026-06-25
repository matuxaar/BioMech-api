package com.biomech.core.network.api

import io.ktor.client.*
import io.ktor.client.call.*
import io.ktor.client.plugins.*
import io.ktor.client.plugins.contentnegotiation.*
import io.ktor.client.request.*
import io.ktor.http.*
import io.ktor.serialization.kotlinx.json.*
import kotlinx.serialization.json.Json

import com.biomech.core.firebase.FirebaseConfig
import com.biomech.core.network.dto.*

private val FIREBASE_API_KEY get() = FirebaseConfig.API_KEY

private fun firebaseClient() = HttpClient {
    install(ContentNegotiation) {
        json(Json { ignoreUnknownKeys = true; isLenient = true })
    }
    install(HttpTimeout) {
        requestTimeoutMillis = 30_000
        connectTimeoutMillis = 15_000
    }
    defaultRequest {
        contentType(ContentType.Application.Json)
    }
}

class AuthApi {
    private val client = firebaseClient()

    suspend fun signInWithPassword(email: String, password: String): FirebaseAuthResponse {
        return client.post("https://identitytoolkit.googleapis.com/v1/accounts:signInWithPassword") {
            parameter("key", FIREBASE_API_KEY)
            setBody(FirebaseSignInRequest(email, password, returnSecureToken = true))
        }.body()
    }

    suspend fun signUp(email: String, password: String): FirebaseAuthResponse {
        return client.post("https://identitytoolkit.googleapis.com/v1/accounts:signUp") {
            parameter("key", FIREBASE_API_KEY)
            setBody(FirebaseSignInRequest(email, password, returnSecureToken = true))
        }.body()
    }

    suspend fun refreshToken(refreshToken: String): FirebaseTokenRefreshResponse {
        return client.post("https://securetoken.googleapis.com/v1/token") {
            parameter("key", FIREBASE_API_KEY)
            setBody(FirebaseTokenRefreshRequest(refreshToken, "refresh_token"))
        }.body()
    }
}
