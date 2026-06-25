package com.biomech.core.analytics

class AnalyticsLogger : AnalyticsTracker {
    override fun track(event: AnalyticsEvent) {
        println("[Analytics] ${event.name}: ${event.params}")
    }

    override fun trackScreenView(screenName: String) {
        println("[Analytics] Screen: $screenName")
    }
}
