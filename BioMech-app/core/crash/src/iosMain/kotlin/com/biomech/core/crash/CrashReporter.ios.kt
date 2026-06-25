package com.biomech.core.crash

import com.biomech.core.common.PlatformContext
import kotlinx.cinterop.ExperimentalForeignApi

actual class CrashReporter {
    private val tags = mutableMapOf<String, String>()

    actual fun init() {
    }

    actual fun recordException(throwable: Throwable) {
    }

    actual fun setUserId(userId: String) {
    }

    actual fun addTag(key: String, value: String) {
        tags[key] = value
    }
}

actual fun createCrashReporter(context: PlatformContext): CrashReporter =
    CrashReporter()
