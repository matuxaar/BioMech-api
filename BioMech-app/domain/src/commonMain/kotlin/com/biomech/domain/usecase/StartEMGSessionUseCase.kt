package com.biomech.domain.usecase

import com.biomech.core.common.AppResult
import com.biomech.domain.model.EMGSession
import com.biomech.domain.repository.EMGRepository

class StartEMGSessionUseCase(private val emgRepository: EMGRepository) {
    suspend operator fun invoke(deviceId: String, label: String): AppResult<EMGSession> {
        return emgRepository.startSession(deviceId, label)
    }
}
