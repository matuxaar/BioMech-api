package com.biomech.core.database.dao

import androidx.room.Dao
import androidx.room.Insert
import androidx.room.OnConflictStrategy
import androidx.room.Query
import com.biomech.core.database.entity.CachedEMGSession

@Dao
interface EMGSessionDao {
    @Query("SELECT * FROM cached_emg_sessions ORDER BY cachedAt DESC")
    suspend fun getAll(): List<CachedEMGSession>

    @Query("SELECT * FROM cached_emg_sessions WHERE id = :id")
    suspend fun getById(id: String): CachedEMGSession?

    @Insert(onConflict = OnConflictStrategy.REPLACE)
    suspend fun insertAll(sessions: List<CachedEMGSession>)

    @Insert(onConflict = OnConflictStrategy.REPLACE)
    suspend fun insert(session: CachedEMGSession)

    @Query("DELETE FROM cached_emg_sessions")
    suspend fun clearAll()
}
