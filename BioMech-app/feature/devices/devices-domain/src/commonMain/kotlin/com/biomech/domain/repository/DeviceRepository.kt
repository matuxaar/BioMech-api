package com.biomech.domain.repository

import com.biomech.core.common.AppResult
import com.biomech.domain.model.Device

interface DeviceRepository {
    suspend fun getDevices(): AppResult<List<Device>>
    suspend fun createDevice(type: String, name: String, hwVersion: String): AppResult<Device>
}
