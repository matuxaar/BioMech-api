package com.biomech.core.network.api

import com.biomech.core.network.createHttpClient
import com.biomech.core.network.dto.*
import io.ktor.client.call.*
import io.ktor.client.request.*

class TrainingApi {
    private val client = createHttpClient()

    suspend fun createJob(request: CreateTrainingJobRequest): TrainingJobDto {
        return client.post("/training/jobs") {
            setBody(request)
        }.body()
    }

    suspend fun getJobs(): List<TrainingJobDto> {
        return client.get("/training/jobs").body()
    }
}
