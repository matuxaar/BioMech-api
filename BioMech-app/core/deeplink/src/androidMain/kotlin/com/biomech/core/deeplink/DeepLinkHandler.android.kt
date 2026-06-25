package com.biomech.core.deeplink

import android.net.Uri
import com.biomech.core.common.PlatformContext

actual class DeepLinkHandler(context: PlatformContext) {
    private val androidContext = context.androidContext

    actual fun handleDeepLink(uri: String): DeepLink {
        val parsed = Uri.parse(uri)
        val path = parsed.path?.trim('/') ?: ""
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

actual fun createDeepLinkHandler(context: PlatformContext): DeepLinkHandler = DeepLinkHandler(context)
