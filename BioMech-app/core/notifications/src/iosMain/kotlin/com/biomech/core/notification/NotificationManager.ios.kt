package com.biomech.core.notification

import com.biomech.core.common.PlatformContext

actual class NotificationManager {
    actual fun schedule(notification: NotificationData) {
    }

    actual fun cancelAll() {
    }
}

actual fun createNotificationManager(context: PlatformContext): NotificationManager =
    NotificationManager()
