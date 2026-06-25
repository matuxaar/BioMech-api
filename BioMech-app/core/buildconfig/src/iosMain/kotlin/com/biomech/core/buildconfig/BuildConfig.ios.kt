package com.biomech.core.buildconfig

import platform.Foundation.NSBundle

actual object BuildConfig {
    actual val isDebug: Boolean
        get() {
            val dict = NSBundle.mainBundle.infoDictionary ?: return false
            return dict["Debug"] as? Boolean ?: false
        }

    actual val versionName: String
        get() = NSBundle.mainBundle.infoDictionary
            ?.get("CFBundleShortVersionString") as? String ?: "1.0"

    actual val versionCode: Int
        get() = (NSBundle.mainBundle.infoDictionary
            ?.get("CFBundleVersion") as? String)?.toIntOrNull() ?: 1

    actual val baseUrl: String = "https://api.biomech.app"
}
