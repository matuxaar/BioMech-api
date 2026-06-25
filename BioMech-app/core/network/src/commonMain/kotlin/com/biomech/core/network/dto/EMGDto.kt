package com.biomech.core.network.dto

import kotlinx.serialization.Serializable

@Serializable
data class CreateSessionRequest(
    val device_id: String,
    val label: String,
)

@Serializable
data class SessionDto(
    val id: String,
    val device_id: String,
    val label: String,
    val started_at: String,
    val ended_at: String? = null,
)

@Serializable
data class SampleRequest(
    val timestamp: String,
    val channel_1: Float,
    val channel_2: Float,
    val channel_3: Float,
    val channel_4: Float,
    val channel_5: Float,
    val channel_6: Float,
    val channel_7: Float,
    val channel_8: Float,
)

@Serializable
data class BatchSamplesRequest(
    val samples: List<SampleRequest>,
)

@Serializable
data class SampleDto(
    val id: String,
    val session_id: String,
    val timestamp: String,
    val channel_1: Float,
    val channel_2: Float,
    val channel_3: Float,
    val channel_4: Float,
    val channel_5: Float,
    val channel_6: Float,
    val channel_7: Float,
    val channel_8: Float,
)
