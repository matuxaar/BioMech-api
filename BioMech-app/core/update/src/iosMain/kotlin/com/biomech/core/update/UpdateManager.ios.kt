package com.biomech.core.update

import com.biomech.core.common.PlatformContext
import platform.Foundation.NSURL
import platform.UIKit.UIApplication

actual class UpdateManager {
    actual suspend fun checkForUpdate(currentVersionCode: Int): UpdateResult {
        return UpdateResult.UpToDate
    }

    actual fun openStore() {
        val url = NSURL.URLWithString("https://apps.apple.com/app/idXXXXXXXX")
        url?.let { UIApplication.sharedApplication.openURL(it) }
    }
}

actual fun createUpdateManager(context: PlatformContext): UpdateManager =
    UpdateManager()
