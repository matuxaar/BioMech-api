package com.biomech.core.analytics

interface AnalyticsTracker {
    fun track(event: AnalyticsEvent)
    fun trackScreenView(screenName: String)
}
