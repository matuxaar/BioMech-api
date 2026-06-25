package com.biomech.core.ble

import com.biomech.domain.model.EMGSample
import kotlinx.coroutines.flow.Flow

interface BleManager {
    val scannedDevices: Flow<List<BleDevice>>
    val emgDataStream: Flow<EMGSample>
    val connectionState: Flow<ConnectionState>

    suspend fun startScanning()
    suspend fun stopScanning()
    suspend fun connect(deviceId: String)
    suspend fun disconnect()
}

data class BleDevice(
    val id: String,
    val name: String,
    val rssi: Int,
)

enum class ConnectionState {
    DISCONNECTED, CONNECTING, CONNECTED, DISCONNECTING
}
