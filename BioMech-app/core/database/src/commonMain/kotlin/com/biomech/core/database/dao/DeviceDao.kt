package com.biomech.core.database.dao

import androidx.room.Dao
import androidx.room.Insert
import androidx.room.OnConflictStrategy
import androidx.room.Query
import com.biomech.core.database.entity.CachedDevice

@Dao
interface DeviceDao {
    @Query("SELECT * FROM cached_devices ORDER BY cachedAt DESC")
    suspend fun getAll(): List<CachedDevice>

    @Query("SELECT * FROM cached_devices WHERE id = :id")
    suspend fun getById(id: String): CachedDevice?

    @Insert(onConflict = OnConflictStrategy.REPLACE)
    suspend fun insertAll(devices: List<CachedDevice>)

    @Insert(onConflict = OnConflictStrategy.REPLACE)
    suspend fun insert(device: CachedDevice)

    @Query("DELETE FROM cached_devices")
    suspend fun clearAll()
}
