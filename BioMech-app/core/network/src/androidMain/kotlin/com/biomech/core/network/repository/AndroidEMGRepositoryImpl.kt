package com.biomech.core.network.repository

import com.biomech.core.common.AppResult
import com.biomech.core.common.currentTimeMillis
import com.biomech.core.database.dao.EMGSessionDao
import com.biomech.core.database.entity.CachedEMGSession
import com.biomech.core.network.api.EMGApi
import com.biomech.domain.model.EMGSession

class AndroidEMGRepositoryImpl(
    emgApi: EMGApi,
    private val emgSessionDao: EMGSessionDao,
) : EMGRepositoryImpl(emgApi) {

    override suspend fun afterSessionCreated(session: EMGSession) {
        emgSessionDao.insert(session.toCached())
    }

    override suspend fun onSessionsFetched(sessions: List<EMGSession>) {
        emgSessionDao.clearAll()
        emgSessionDao.insertAll(sessions.map { it.toCached() })
    }

    override suspend fun getCachedSessions(): AppResult<List<EMGSession>>? {
        val cached = emgSessionDao.getAll()
        if (cached.isNotEmpty()) {
            return AppResult.Success(cached.map { it.toDomain() })
        }
        return null
    }
}

private fun EMGSession.toCached() = CachedEMGSession(
    id = id,
    deviceId = deviceId,
    label = label,
    startedAt = startedAt,
    endedAt = endedAt,
    cachedAt = currentTimeMillis(),
)

private fun CachedEMGSession.toDomain() = EMGSession(
    id = id,
    deviceId = deviceId,
    label = label,
    startedAt = startedAt,
    endedAt = endedAt,
)
