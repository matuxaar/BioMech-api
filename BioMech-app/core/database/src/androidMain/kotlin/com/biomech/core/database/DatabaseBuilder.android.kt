package com.biomech.core.database

import androidx.room.Room
import com.biomech.core.common.PlatformContext

actual fun createRoomDatabase(context: PlatformContext): AppDatabase {
    return Room.databaseBuilder(
        context.androidContext,
        AppDatabase::class.java,
        "biomech_cache.db",
    ).build()
}
