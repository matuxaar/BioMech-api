package com.biomech.core.database.dao

import androidx.room.Dao
import androidx.room.Insert
import androidx.room.OnConflictStrategy
import androidx.room.Query
import com.biomech.core.database.entity.CachedTrainingJob

@Dao
interface TrainingJobDao {
    @Query("SELECT * FROM cached_training_jobs ORDER BY cachedAt DESC")
    suspend fun getAll(): List<CachedTrainingJob>

    @Query("SELECT * FROM cached_training_jobs WHERE id = :id")
    suspend fun getById(id: String): CachedTrainingJob?

    @Insert(onConflict = OnConflictStrategy.REPLACE)
    suspend fun insertAll(jobs: List<CachedTrainingJob>)

    @Insert(onConflict = OnConflictStrategy.REPLACE)
    suspend fun insert(job: CachedTrainingJob)

    @Query("DELETE FROM cached_training_jobs")
    suspend fun clearAll()
}
