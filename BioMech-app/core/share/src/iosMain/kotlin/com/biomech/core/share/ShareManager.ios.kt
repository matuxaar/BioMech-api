package com.biomech.core.share

import com.biomech.core.common.PlatformContext
import platform.UIKit.UIApplication
import platform.UIKit.UIActivityViewController

actual class ShareManager {
    actual fun share(data: ShareData) {
        val items = mutableListOf<Any>()
        data.text?.let { items.add(it) }
        val controller = UIActivityViewController(
            activityItems = items,
            applicationActivities = null,
        )
        val root = UIApplication.sharedApplication.keyWindow?.rootViewController
        root?.presentViewController(controller, animated = true, completion = null)
    }
}

actual fun createShareManager(context: PlatformContext): ShareManager =
    ShareManager()
