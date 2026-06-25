package com.biomech.core.notification

import com.biomech.core.common.PlatformContext

data class NotificationData(
    val title: String,
    val body: String,
    val channelId: String = "default",
    val delayMillis: Long = 0,
)

expect class NotificationManager {
    fun schedule(notification: NotificationData)
    fun cancelAll()
}

expect fun createNotificationManager(context: PlatformContext): NotificationManager
