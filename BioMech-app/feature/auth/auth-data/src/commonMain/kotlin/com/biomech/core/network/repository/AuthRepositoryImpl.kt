package com.biomech.core.network.repository

import com.biomech.core.common.AppResult
import com.biomech.core.network.ApiConfig
import com.biomech.core.network.api.AuthApi
import com.biomech.core.network.createHttpClient
import com.biomech.core.storage.KeyValueStorage
import com.biomech.domain.model.User
import com.biomech.domain.repository.AuthRepository
import io.ktor.client.call.*
import io.ktor.client.request.*
import io.ktor.http.*

private const val KEY_ACCESS_TOKEN = "auth_access_token"
private const val KEY_REFRESH_TOKEN = "auth_refresh_token"
private const val KEY_USER_ID = "auth_user_id"
private const val KEY_USER_EMAIL = "auth_user_email"

class AuthRepositoryImpl(
    private val authApi: AuthApi,
    private val storage: KeyValueStorage,
) : AuthRepository {

    override suspend fun login(email: String, password: String): AppResult<User> {
        return try {
            val response = authApi.signInWithPassword(email, password)
            saveTokens(response.idToken, response.refreshToken, response.localId, response.email)
            syncUser()
            AppResult.Success(User(id = response.localId, email = response.email))
        } catch (e: Exception) {
            AppResult.Error(e.message ?: "Login failed")
        }
    }

    override suspend fun register(email: String, password: String): AppResult<User> {
        return try {
            val response = authApi.signUp(email, password)
            saveTokens(response.idToken, response.refreshToken, response.localId, response.email)
            syncUser()
            AppResult.Success(User(id = response.localId, email = response.email))
        } catch (e: Exception) {
            AppResult.Error(e.message ?: "Registration failed")
        }
    }

    override suspend fun refreshToken(): AppResult<User> {
        val savedRefresh = storage.getString(KEY_REFRESH_TOKEN)
        if (savedRefresh.isEmpty()) return AppResult.Error("No refresh token")
        return try {
            val response = authApi.refreshToken(savedRefresh)
            saveTokens(response.idToken, response.refreshToken, response.userId, response.email)
            AppResult.Success(User(id = response.userId, email = response.email))
        } catch (e: Exception) {
            AppResult.Error(e.message ?: "Token refresh failed")
        }
    }

    override suspend fun logout() {
        storage.remove(KEY_ACCESS_TOKEN)
        storage.remove(KEY_REFRESH_TOKEN)
        storage.remove(KEY_USER_ID)
        storage.remove(KEY_USER_EMAIL)
        ApiConfig.token = null
    }

    override suspend fun getToken(): String? {
        val token = storage.getString(KEY_ACCESS_TOKEN)
        return token.ifEmpty { null }
    }

    private suspend fun syncUser() {
        try {
            createHttpClient().post("/auth/firebase")
        } catch (_: Exception) { }
    }

    private fun saveTokens(access: String, refresh: String, userId: String, email: String) {
        storage.putString(KEY_ACCESS_TOKEN, access)
        storage.putString(KEY_REFRESH_TOKEN, refresh)
        storage.putString(KEY_USER_ID, userId)
        storage.putString(KEY_USER_EMAIL, email)
        ApiConfig.token = access
    }
}
