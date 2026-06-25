package com.biomech.core.database.entity

import androidx.room.Entity
import androidx.room.PrimaryKey

@Entity(tableName = "cached_emg_sessions")
data class CachedEMGSession(
    @PrimaryKey val id: String,
    val deviceId: String,
    val label: String,
    val startedAt: Long,
    val endedAt: Long?,
    val cachedAt: Long,
)
