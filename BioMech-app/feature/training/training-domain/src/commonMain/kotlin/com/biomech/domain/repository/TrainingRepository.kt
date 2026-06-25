package com.biomech.domain.repository

import com.biomech.core.common.AppResult
import com.biomech.domain.model.TrainingJob

interface TrainingRepository {
    suspend fun createJob(sessionIds: List<String>): AppResult<TrainingJob>
    suspend fun getJobs(): AppResult<List<TrainingJob>>
}
