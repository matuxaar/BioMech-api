package com.biomech.feature.dashboard

import com.biomech.core.ble.BleManager
import com.biomech.core.ble.ConnectionState
import com.biomech.core.mvi.BaseAction
import com.biomech.core.mvi.BaseEvent
import com.biomech.core.mvi.BaseState
import com.biomech.core.mvi.BaseViewModel
import com.biomech.domain.usecase.StartEMGSessionUseCase
import kotlinx.coroutines.channels.Channel
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.launch

data class DashboardState(
    val deviceConnected: Boolean = false,
    val emgData: List<Float> = emptyList(),
    val isRecording: Boolean = false,
) : BaseState

sealed class DashboardAction : BaseAction {
    data object StartRecording : DashboardAction()
    data object StopRecording : DashboardAction()
}

sealed class DashboardEvent : BaseEvent

class DashboardViewModel(
    private val bleManager: BleManager,
    private val startSessionUseCase: StartEMGSessionUseCase,
) : BaseViewModel<DashboardState, DashboardAction, DashboardEvent>() {

    override val _state = MutableStateFlow(DashboardState())
    override val _event = Channel<DashboardEvent>(Channel.BUFFERED)

    init {
        scope.launch {
            bleManager.connectionState.collect { state ->
                _state.value = _state.value.copy(
                    deviceConnected = state == ConnectionState.CONNECTED
                )
            }
        }
        scope.launch {
            bleManager.emgDataStream.collect { sample ->
                val data = _state.value.emgData.toMutableList()
                data.add(sample.channel1)
                if (data.size > 200) data.removeAt(0)
                _state.value = _state.value.copy(emgData = data)
            }
        }
    }

    override suspend fun handleAction(action: DashboardAction) {
        when (action) {
            DashboardAction.StartRecording -> {
                _state.value = _state.value.copy(isRecording = true)
            }
            DashboardAction.StopRecording -> {
                _state.value = _state.value.copy(isRecording = false)
            }
        }
    }
}
