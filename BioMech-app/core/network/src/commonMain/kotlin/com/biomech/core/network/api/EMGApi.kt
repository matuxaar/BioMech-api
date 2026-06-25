package com.biomech.core.network.api

import com.biomech.core.network.createHttpClient
import com.biomech.core.network.dto.*
import io.ktor.client.call.*
import io.ktor.client.request.*

class EMGApi {
    private val client = createHttpClient()

    suspend fun startSession(request: CreateSessionRequest): SessionDto {
        return client.post("/emg/sessions") {
            setBody(request)
        }.body()
    }

    suspend fun endSession(sessionId: String) {
        client.post("/emg/sessions/$sessionId/end")
    }

    suspend fun getSessions(): List<SessionDto> {
        return client.get("/emg/sessions").body()
    }

    suspend fun addSamplesBatch(sessionId: String, request: BatchSamplesRequest) {
        client.post("/emg/sessions/$sessionId/samples/batch") {
            setBody(request)
        }
    }

    suspend fun getSamples(sessionId: String): List<SampleDto> {
        return client.get("/emg/sessions/$sessionId/samples").body()
    }
}
