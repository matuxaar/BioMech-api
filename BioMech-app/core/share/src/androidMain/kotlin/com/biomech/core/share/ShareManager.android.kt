package com.biomech.core.share

import android.content.Intent
import com.biomech.core.common.PlatformContext

actual class ShareManager(private val context: android.content.Context) {
    actual fun share(data: ShareData) {
        val intent = Intent(Intent.ACTION_SEND).apply {
            type = "text/plain"
            putExtra(Intent.EXTRA_TEXT, data.text)
            data.title?.let { putExtra(Intent.EXTRA_SUBJECT, it) }
        }
        val chooser = Intent.createChooser(intent, data.title ?: "Share")
        context.startActivity(chooser)
    }
}

actual fun createShareManager(context: PlatformContext): ShareManager =
    ShareManager(context.androidContext)
