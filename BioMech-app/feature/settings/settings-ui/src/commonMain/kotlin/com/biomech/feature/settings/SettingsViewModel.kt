package com.biomech.feature.settings

import com.biomech.core.mvi.BaseAction
import com.biomech.core.mvi.BaseEvent
import com.biomech.core.mvi.BaseState
import com.biomech.core.mvi.BaseViewModel
import com.biomech.domain.repository.AuthRepository
import kotlinx.coroutines.channels.Channel
import kotlinx.coroutines.flow.MutableStateFlow

data class SettingsState(
    val serverUrl: String = "http://10.0.2.2:8080/api/v1",
) : BaseState

sealed class SettingsAction : BaseAction {
    data class UpdateServerUrl(val url: String) : SettingsAction()
    data object Logout : SettingsAction()
}

sealed class SettingsEvent : BaseEvent {
    data object NavigateToLogin : SettingsEvent()
}

class SettingsViewModel(
    private val authRepository: AuthRepository,
) : BaseViewModel<SettingsState, SettingsAction, SettingsEvent>() {

    override val _state = MutableStateFlow(SettingsState())
    override val _event = Channel<SettingsEvent>(Channel.BUFFERED)

    override suspend fun handleAction(action: SettingsAction) {
        when (action) {
            is SettingsAction.UpdateServerUrl -> {
                _state.value = _state.value.copy(serverUrl = action.url)
            }
            SettingsAction.Logout -> {
                authRepository.logout()
                _event.send(SettingsEvent.NavigateToLogin)
            }
        }
    }
}
