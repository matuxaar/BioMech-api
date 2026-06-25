package com.biomech.core.network.repository

import com.biomech.core.common.AppResult
import com.biomech.core.network.api.EMGApi
import com.biomech.core.network.dto.CreateSessionRequest
import com.biomech.domain.model.EMGSession
import com.biomech.domain.repository.EMGRepository

open class EMGRepositoryImpl(
    protected val emgApi: EMGApi,
) : EMGRepository {

    override suspend fun startSession(deviceId: String, label: String): AppResult<EMGSession> {
        return try {
            val dto = emgApi.startSession(CreateSessionRequest(deviceId, label))
            val session = dto.toDomain()
            afterSessionCreated(session)
            AppResult.Success(session)
        } catch (e: Exception) {
            AppResult.Error(e.message ?: "Failed to start session")
        }
    }

    override suspend fun endSession(sessionId: String): AppResult<Unit> {
        return try {
            emgApi.endSession(sessionId)
            AppResult.Success(Unit)
        } catch (e: Exception) {
            AppResult.Error(e.message ?: "Failed to end session")
        }
    }

    override suspend fun getSessions(): AppResult<List<EMGSession>> {
        return try {
            val dtos = emgApi.getSessions()
            val sessions = dtos.map { it.toDomain() }
            onSessionsFetched(sessions)
            AppResult.Success(sessions)
        } catch (e: Exception) {
            getCachedSessions() ?: AppResult.Error(e.message ?: "Failed to fetch sessions")
        }
    }

    protected open suspend fun afterSessionCreated(session: EMGSession) {}
    protected open suspend fun onSessionsFetched(sessions: List<EMGSession>) {}
    protected open suspend fun getCachedSessions(): AppResult<List<EMGSession>>? = null
}

internal fun com.biomech.core.network.dto.SessionDto.toDomain() = EMGSession(
    id = id,
    deviceId = device_id,
    label = label,
    startedAt = 0L,
    endedAt = null,
)
