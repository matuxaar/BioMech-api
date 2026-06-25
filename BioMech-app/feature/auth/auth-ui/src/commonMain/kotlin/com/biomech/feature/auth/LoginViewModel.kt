package com.biomech.feature.auth

import com.biomech.core.common.AppResult
import com.biomech.core.mvi.BaseAction
import com.biomech.core.mvi.BaseEvent
import com.biomech.core.mvi.BaseState
import com.biomech.core.mvi.BaseViewModel
import com.biomech.domain.repository.AuthRepository
import kotlinx.coroutines.channels.Channel
import kotlinx.coroutines.flow.MutableStateFlow

data class LoginState(
    val isLoading: Boolean = false,
    val error: String? = null,
) : BaseState

sealed class LoginAction : BaseAction {
    data class Login(val email: String, val password: String) : LoginAction()
    data class Register(val email: String, val password: String) : LoginAction()
}

sealed class LoginEvent : BaseEvent {
    data object NavigateToMain : LoginEvent()
}

class LoginViewModel(
    private val authRepository: AuthRepository,
) : BaseViewModel<LoginState, LoginAction, LoginEvent>() {

    override val _state = MutableStateFlow(LoginState())
    override val _event = Channel<LoginEvent>(Channel.BUFFERED)

    override suspend fun handleAction(action: LoginAction) {
        when (action) {
            is LoginAction.Login -> login(action.email, action.password)
            is LoginAction.Register -> register(action.email, action.password)
        }
    }

    private suspend fun login(email: String, password: String) {
        _state.value = LoginState(isLoading = true)
        when (val result = authRepository.login(email, password)) {
            is AppResult.Success -> {
                _state.value = LoginState()
                _event.send(LoginEvent.NavigateToMain)
            }
            is AppResult.Error -> {
                _state.value = LoginState(error = result.message)
            }
        }
    }

    private suspend fun register(email: String, password: String) {
        _state.value = LoginState(isLoading = true)
        when (val result = authRepository.register(email, password)) {
            is AppResult.Success -> {
                _state.value = LoginState()
                _event.send(LoginEvent.NavigateToMain)
            }
            is AppResult.Error -> {
                _state.value = LoginState(error = result.message)
            }
        }
    }
}
