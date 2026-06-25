package com.biomech.domain.model

data class Device(
    val id: String,
    val type: DeviceType,
    val name: String,
    val hwVersion: String,
)

enum class DeviceType { PROSTHETIC, SENSOR }
