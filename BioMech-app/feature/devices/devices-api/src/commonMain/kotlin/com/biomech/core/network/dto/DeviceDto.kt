package com.biomech.core.network.dto

import kotlinx.serialization.Serializable

@Serializable
data class CreateDeviceRequest(
    val type: String,
    val name: String,
    val hw_version: String,
)

@Serializable
data class DeviceDto(
    val id: String,
    val type: String,
    val name: String,
    val hw_version: String,
)
