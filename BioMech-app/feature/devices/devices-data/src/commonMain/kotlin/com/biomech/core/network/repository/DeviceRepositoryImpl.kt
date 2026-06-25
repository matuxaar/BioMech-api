package com.biomech.core.network.repository

import com.biomech.core.common.AppResult
import com.biomech.core.network.api.DeviceApi
import com.biomech.core.network.dto.CreateDeviceRequest
import com.biomech.domain.model.Device
import com.biomech.domain.model.DeviceType
import com.biomech.domain.repository.DeviceRepository

open class DeviceRepositoryImpl(
    protected val deviceApi: DeviceApi,
) : DeviceRepository {

    override suspend fun getDevices(): AppResult<List<Device>> {
        return try {
            val dtos = deviceApi.getDevices()
            val devices = dtos.map { it.toDomain() }
            afterDevicesFetched(devices)
            AppResult.Success(devices)
        } catch (e: Exception) {
            getCachedDevices() ?: AppResult.Error(e.message ?: "Failed to fetch devices")
        }
    }

    override suspend fun createDevice(type: String, name: String, hwVersion: String): AppResult<Device> {
        return try {
            val dto = deviceApi.createDevice(CreateDeviceRequest(type, name, hwVersion))
            val device = dto.toDomain()
            afterDeviceCreated(device)
            AppResult.Success(device)
        } catch (e: Exception) {
            AppResult.Error(e.message ?: "Failed to create device")
        }
    }

    protected open suspend fun afterDevicesFetched(devices: List<Device>) {}
    protected open suspend fun getCachedDevices(): AppResult<List<Device>>? = null
    protected open suspend fun afterDeviceCreated(device: Device) {}
}

internal fun com.biomech.core.network.dto.DeviceDto.toDomain() = Device(
    id = id,
    type = if (type == "prosthetic") DeviceType.PROSTHETIC else DeviceType.SENSOR,
    name = name,
    hwVersion = hw_version,
)
