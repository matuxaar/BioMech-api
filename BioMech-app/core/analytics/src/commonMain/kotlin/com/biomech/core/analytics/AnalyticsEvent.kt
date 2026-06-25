package com.biomech.core.analytics

data class AnalyticsEvent(
    val name: String,
    val params: Map<String, Any> = emptyMap(),
)
