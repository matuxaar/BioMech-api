package com.biomech.feature.profile

import com.biomech.core.common.AppResult
import com.biomech.core.mvi.BaseAction
import com.biomech.core.mvi.BaseEvent
import com.biomech.core.mvi.BaseState
import com.biomech.core.mvi.BaseViewModel
import com.biomech.domain.repository.AuthRepository
import kotlinx.coroutines.channels.Channel
import kotlinx.coroutines.flow.MutableStateFlow

data class ProfileState(
    val email: String = "",
    val isLoading: Boolean = false,
) : BaseState

sealed class ProfileAction : BaseAction {
    data object LoadProfile : ProfileAction()
    data object Logout : ProfileAction()
    data class SetEmail(val email: String) : ProfileAction()
}

sealed class ProfileEvent : BaseEvent {
    data object NavigateToLogin : ProfileEvent()
}

class ProfileViewModel(
    private val authRepository: AuthRepository,
) : BaseViewModel<ProfileState, ProfileAction, ProfileEvent>() {

    override val _state = MutableStateFlow(ProfileState())
    override val _event = Channel<ProfileEvent>(Channel.BUFFERED)

    override suspend fun handleAction(action: ProfileAction) {
        when (action) {
            ProfileAction.LoadProfile -> loadProfile()
            ProfileAction.Logout -> {
                authRepository.logout()
                _event.send(ProfileEvent.NavigateToLogin)
            }
            is ProfileAction.SetEmail -> {
                _state.value = _state.value.copy(email = action.email)
            }
        }
    }

    private suspend fun loadProfile() {
        _state.value = _state.value.copy(isLoading = true)
        authRepository.refreshToken()
        _state.value = _state.value.copy(isLoading = false)
    }
}
