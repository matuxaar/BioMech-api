package com.biomech.core.deeplink

import com.biomech.core.common.PlatformContext

actual class DeepLinkHandler {
    actual fun handleDeepLink(uri: String): DeepLink {
        val path = uri.trim('/')
        return when {
            path == "dashboard" -> DeepLink.OpenDashboard
            path == "devices" -> DeepLink.OpenDevices
            path == "settings" -> DeepLink.OpenSettings
            path.startsWith("training/") -> {
                val id = path.removePrefix("training/")
                DeepLink.OpenTraining(id)
            }
            path.startsWith("device/") -> {
                val id = path.removePrefix("device/")
                DeepLink.OpenDevice(id)
            }
            else -> DeepLink.Unknown(uri)
        }
    }
}

actual fun createDeepLinkHandler(context: PlatformContext): DeepLinkHandler = DeepLinkHandler()
