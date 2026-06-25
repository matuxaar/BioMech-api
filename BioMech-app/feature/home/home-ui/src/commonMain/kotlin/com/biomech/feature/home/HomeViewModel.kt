package com.biomech.feature.home

import com.biomech.core.common.AppResult
import com.biomech.core.mvi.BaseAction
import com.biomech.core.mvi.BaseEvent
import com.biomech.core.mvi.BaseState
import com.biomech.core.mvi.BaseViewModel
import com.biomech.domain.model.Device
import com.biomech.domain.repository.DeviceRepository
import kotlinx.coroutines.channels.Channel
import kotlinx.coroutines.flow.MutableStateFlow

data class HomeState(
    val devices: List<Device> = emptyList(),
    val isLoading: Boolean = false,
    val error: String? = null,
) : BaseState

sealed class HomeAction : BaseAction {
    data object LoadDevices : HomeAction()
}

sealed class HomeEvent : BaseEvent

class HomeViewModel(
    private val deviceRepository: DeviceRepository,
) : BaseViewModel<HomeState, HomeAction, HomeEvent>() {

    override val _state = MutableStateFlow(HomeState())
    override val _event = Channel<HomeEvent>(Channel.BUFFERED)

    override suspend fun handleAction(action: HomeAction) {
        when (action) {
            HomeAction.LoadDevices -> loadDevices()
        }
    }

    private suspend fun loadDevices() {
        _state.value = HomeState(isLoading = true)
        when (val result = deviceRepository.getDevices()) {
            is AppResult.Success -> {
                _state.value = HomeState(devices = result.data)
            }
            is AppResult.Error -> {
                _state.value = HomeState(error = result.message)
            }
        }
    }
}
