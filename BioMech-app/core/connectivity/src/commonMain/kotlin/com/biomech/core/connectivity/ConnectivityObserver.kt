package com.biomech.core.connectivity

import kotlinx.coroutines.flow.StateFlow

enum class ConnectionStatus {
    AVAILABLE,
    UNAVAILABLE,
    LOSING,
}

interface ConnectivityObserver {
    val status: StateFlow<ConnectionStatus>
    fun start()
    fun stop()
}
