package com.biomech.domain.model

data class EMGSession(
    val id: String,
    val deviceId: String,
    val label: String,
    val startedAt: Long,
    val endedAt: Long?,
)

data class EMGSample(
    val channel1: Float,
    val channel2: Float,
    val channel3: Float,
    val channel4: Float,
    val channel5: Float,
    val channel6: Float,
    val channel7: Float,
    val channel8: Float,
)
