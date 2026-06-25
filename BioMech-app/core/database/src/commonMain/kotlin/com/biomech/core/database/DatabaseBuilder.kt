package com.biomech.core.database

import com.biomech.core.common.PlatformContext

expect fun createRoomDatabase(context: PlatformContext): AppDatabase
