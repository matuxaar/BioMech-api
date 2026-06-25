package com.biomech.core.network.repository

import com.biomech.core.common.AppResult
import com.biomech.core.network.api.TrainingApi
import com.biomech.core.network.dto.CreateTrainingJobRequest
import com.biomech.domain.model.TrainingJob
import com.biomech.domain.model.TrainingStatus
import com.biomech.domain.repository.TrainingRepository

open class TrainingRepositoryImpl(
    protected val trainingApi: TrainingApi,
) : TrainingRepository {

    override suspend fun createJob(sessionIds: List<String>): AppResult<TrainingJob> {
        return try {
            val dto = trainingApi.createJob(CreateTrainingJobRequest(sessionIds))
            val job = dto.toDomain()
            afterJobCreated(job)
            AppResult.Success(job)
        } catch (e: Exception) {
            AppResult.Error(e.message ?: "Failed to create training job")
        }
    }

    override suspend fun getJobs(): AppResult<List<TrainingJob>> {
        return try {
            val dtos = trainingApi.getJobs()
            val jobs = dtos.map { it.toDomain() }
            afterJobsFetched(jobs)
            AppResult.Success(jobs)
        } catch (e: Exception) {
            getCachedJobs() ?: AppResult.Error(e.message ?: "Failed to fetch training jobs")
        }
    }

    protected open suspend fun afterJobCreated(job: TrainingJob) {}
    protected open suspend fun afterJobsFetched(jobs: List<TrainingJob>) {}
    protected open suspend fun getCachedJobs(): AppResult<List<TrainingJob>>? = null
}

internal fun com.biomech.core.network.dto.TrainingJobDto.toDomain() = TrainingJob(
    id = id,
    sessionIds = session_ids,
    status = when (status) {
        "pending" -> TrainingStatus.PENDING
        "running" -> TrainingStatus.RUNNING
        "completed" -> TrainingStatus.COMPLETED
        "failed" -> TrainingStatus.FAILED
        else -> TrainingStatus.PENDING
    },
    accuracy = accuracy,
)
