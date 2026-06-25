package com.biomech.core.deeplink

sealed class DeepLink {
    data object OpenDashboard : DeepLink()
    data object OpenDevices : DeepLink()
    data class OpenTraining(val sessionId: String) : DeepLink()
    data class OpenDevice(val deviceId: String) : DeepLink()
    data object OpenSettings : DeepLink()
    data class Unknown(val uri: String) : DeepLink()
}
