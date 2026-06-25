package com.biomech.core.buildconfig

expect object BuildConfig {
    val isDebug: Boolean
    val versionName: String
    val versionCode: Int
    val baseUrl: String
}
