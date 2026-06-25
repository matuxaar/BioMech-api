package com.biomech.core.network.repository

import com.biomech.core.common.AppResult
import com.biomech.core.common.currentTimeMillis
import com.biomech.core.database.dao.TrainingJobDao
import com.biomech.core.database.entity.CachedTrainingJob
import com.biomech.core.network.api.TrainingApi
import com.biomech.domain.model.TrainingJob
import com.biomech.domain.model.TrainingStatus

class AndroidTrainingRepositoryImpl(
    trainingApi: TrainingApi,
    private val trainingJobDao: TrainingJobDao,
) : TrainingRepositoryImpl(trainingApi) {

    override suspend fun afterJobCreated(job: TrainingJob) {
        trainingJobDao.insert(job.toCached())
    }

    override suspend fun afterJobsFetched(jobs: List<TrainingJob>) {
        trainingJobDao.clearAll()
        trainingJobDao.insertAll(jobs.map { it.toCached() })
    }

    override suspend fun getCachedJobs(): AppResult<List<TrainingJob>>? {
        val cached = trainingJobDao.getAll()
        if (cached.isNotEmpty()) {
            return AppResult.Success(cached.map { it.toDomain() })
        }
        return null
    }
}

private fun TrainingJob.toCached() = CachedTrainingJob(
    id = id,
    sessionIds = sessionIds.joinToString(","),
    status = when (status) {
        TrainingStatus.PENDING -> "pending"
        TrainingStatus.RUNNING -> "running"
        TrainingStatus.COMPLETED -> "completed"
        TrainingStatus.FAILED -> "failed"
    },
    accuracy = accuracy,
    cachedAt = currentTimeMillis(),
)

private fun CachedTrainingJob.toDomain() = TrainingJob(
    id = id,
    sessionIds = sessionIds.split(",").filter { it.isNotBlank() },
    status = when (status) {
        "pending" -> TrainingStatus.PENDING
        "running" -> TrainingStatus.RUNNING
        "completed" -> TrainingStatus.COMPLETED
        "failed" -> TrainingStatus.FAILED
        else -> TrainingStatus.PENDING
    },
    accuracy = accuracy,
)
