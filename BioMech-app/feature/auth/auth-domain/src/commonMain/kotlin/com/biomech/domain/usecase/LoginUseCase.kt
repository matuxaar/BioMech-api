package com.biomech.domain.usecase

import com.biomech.core.common.AppResult
import com.biomech.domain.model.User
import com.biomech.domain.repository.AuthRepository

class LoginUseCase(private val authRepository: AuthRepository) {
    suspend operator fun invoke(email: String, password: String): AppResult<User> {
        return authRepository.login(email, password)
    }
}
