package com.biomech.core.database.entity

import androidx.room.Entity
import androidx.room.PrimaryKey

@Entity(tableName = "cached_training_jobs")
data class CachedTrainingJob(
    @PrimaryKey val id: String,
    val sessionIds: String,
    val status: String,
    val accuracy: Double,
    val cachedAt: Long,
)
