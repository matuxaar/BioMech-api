package com.biomech.core.database.entity

import androidx.room.Entity
import androidx.room.PrimaryKey

@Entity(tableName = "cached_devices")
data class CachedDevice(
    @PrimaryKey val id: String,
    val type: String,
    val name: String,
    val hwVersion: String,
    val cachedAt: Long,
)
