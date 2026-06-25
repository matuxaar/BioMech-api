package com.biomech.core.notification

import android.app.NotificationChannel
import android.app.NotificationManager as AndroidNotificationManager
import android.app.PendingIntent
import android.app.Notification
import android.content.Context
import android.os.Build
import com.biomech.core.common.PlatformContext

actual class NotificationManager(private val context: Context) {
    private val manager = context.getSystemService(Context.NOTIFICATION_SERVICE) as AndroidNotificationManager

    init {
        if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.O) {
            val channel = NotificationChannel(
                "default", "Default", AndroidNotificationManager.IMPORTANCE_DEFAULT
            )
            manager.createNotificationChannel(channel)
        }
    }

    actual fun schedule(notification: NotificationData) {
        val intent = context.packageManager.getLaunchIntentForPackage(context.packageName) ?: return
        val pendingIntent = PendingIntent.getActivity(
            context, 0, intent, PendingIntent.FLAG_UPDATE_CURRENT or PendingIntent.FLAG_IMMUTABLE
        )
        val built = Notification.Builder(context, notification.channelId)
            .setContentTitle(notification.title)
            .setContentText(notification.body)
            .setSmallIcon(android.R.drawable.ic_dialog_info)
            .setContentIntent(pendingIntent)
            .setAutoCancel(true)
            .build()
        manager.notify(System.currentTimeMillis().toInt(), built)
    }

    actual fun cancelAll() {
        manager.cancelAll()
    }
}

actual fun createNotificationManager(context: PlatformContext): NotificationManager =
    NotificationManager(context.androidContext)
