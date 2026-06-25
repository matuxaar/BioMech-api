package com.biomech.domain.usecase

import com.biomech.core.common.AppResult
import com.biomech.domain.model.User
import com.biomech.domain.repository.AuthRepository
import kotlin.test.Test
import kotlin.test.assertEquals
import kotlin.test.assertTrue
import kotlinx.coroutines.test.runTest

class LoginUseCaseTest {

    private val testUser = User(id = "1", email = "test@example.com")

    @Test
    fun `invoke returns success when auth repository login succeeds`() = runTest {
        val repo = FakeAuthRepository(result = AppResult.Success(testUser))
        val useCase = LoginUseCase(repo)

        val result = useCase("test@example.com", "password123")

        assertTrue(result is AppResult.Success)
        assertEquals(testUser, (result as AppResult.Success).data)
    }

    @Test
    fun `invoke returns error when auth repository login fails`() = runTest {
        val repo = FakeAuthRepository(result = AppResult.Error("Invalid credentials"))
        val useCase = LoginUseCase(repo)

        val result = useCase("test@example.com", "wrong")

        assertTrue(result is AppResult.Error)
        assertEquals("Invalid credentials", (result as AppResult.Error).message)
    }
}

private class FakeAuthRepository(private val result: AppResult<User>) : AuthRepository {
    override suspend fun login(email: String, password: String): AppResult<User> = result
    override suspend fun register(email: String, password: String): AppResult<User> = result
    override suspend fun refreshToken(): AppResult<User> = AppResult.Error("not implemented")
    override suspend fun logout() {}
    override suspend fun getToken(): String? = null
}
