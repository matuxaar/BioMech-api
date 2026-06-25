package com.biomech.core.update

import com.biomech.core.common.PlatformContext

data class AppVersion(
    val versionName: String,
    val versionCode: Int,
    val minVersionCode: Int,
    val updateUrl: String,
)

sealed class UpdateResult {
    data object UpToDate : UpdateResult()
    data object UpdateRecommended : UpdateResult()
    data class UpdateRequired(val updateUrl: String) : UpdateResult()
    data class Error(val message: String) : UpdateResult()
}

expect class UpdateManager {
    suspend fun checkForUpdate(currentVersionCode: Int): UpdateResult
    fun openStore()
}

expect fun createUpdateManager(context: PlatformContext): UpdateManager
