package com.biomech.core.update

import android.content.Intent
import android.net.Uri
import com.biomech.core.common.PlatformContext

actual class UpdateManager(private val context: android.content.Context) {
    actual suspend fun checkForUpdate(currentVersionCode: Int): UpdateResult {
        return UpdateResult.UpToDate
    }

    actual fun openStore() {
        val intent = Intent(Intent.ACTION_VIEW).apply {
            data = Uri.parse("market://details?id=${context.packageName}")
            setPackage("com.android.vending")
        }
        context.startActivity(intent)
    }
}

actual fun createUpdateManager(context: PlatformContext): UpdateManager =
    UpdateManager(context.androidContext)
