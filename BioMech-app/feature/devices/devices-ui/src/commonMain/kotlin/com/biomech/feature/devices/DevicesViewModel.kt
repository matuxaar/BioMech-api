package com.biomech.feature.devices

import com.biomech.core.ble.BleDevice
import com.biomech.core.ble.BleManager
import com.biomech.core.mvi.BaseAction
import com.biomech.core.mvi.BaseEvent
import com.biomech.core.mvi.BaseState
import com.biomech.core.mvi.BaseViewModel
import kotlinx.coroutines.channels.Channel
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.launch

data class DevicesState(
    val scannedDevices: List<BleDevice> = emptyList(),
    val isScanning: Boolean = false,
) : BaseState

sealed class DevicesAction : BaseAction {
    data object StartScan : DevicesAction()
    data object StopScan : DevicesAction()
    data class Connect(val deviceId: String) : DevicesAction()
}

sealed class DevicesEvent : BaseEvent

class DevicesViewModel(
    private val bleManager: BleManager,
) : BaseViewModel<DevicesState, DevicesAction, DevicesEvent>() {

    override val _state = MutableStateFlow(DevicesState())
    override val _event = Channel<DevicesEvent>(Channel.BUFFERED)

    init {
        scope.launch {
            bleManager.scannedDevices.collect { devices ->
                _state.value = _state.value.copy(scannedDevices = devices)
            }
        }
    }

    override suspend fun handleAction(action: DevicesAction) {
        when (action) {
            DevicesAction.StartScan -> {
                _state.value = _state.value.copy(isScanning = true)
                bleManager.startScanning()
            }
            DevicesAction.StopScan -> {
                bleManager.stopScanning()
                _state.value = _state.value.copy(isScanning = false)
            }
            is DevicesAction.Connect -> {
                bleManager.connect(action.deviceId)
            }
        }
    }
}
