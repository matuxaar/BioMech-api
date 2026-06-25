package com.biomech.core.network.dto

import kotlinx.serialization.Serializable

@Serializable
data class CreateTrainingJobRequest(
    val session_ids: List<String>,
)

@Serializable
data class TrainingJobDto(
    val id: String,
    val user_id: String,
    val session_ids: List<String>,
    val status: String,
    val accuracy: Double,
    val created_at: String,
)
