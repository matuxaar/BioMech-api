package com.biomech.core.crash

import com.biomech.core.common.PlatformContext

data class CrashReport(
    val message: String,
    val throwable: Throwable? = null,
    val tags: Map<String, String> = emptyMap(),
)

expect class CrashReporter {
    fun init()
    fun recordException(throwable: Throwable)
    fun setUserId(userId: String)
    fun addTag(key: String, value: String)
}

expect fun createCrashReporter(context: PlatformContext): CrashReporter
