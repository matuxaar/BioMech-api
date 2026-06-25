package com.biomech.core.connectivity

import kotlinx.coroutines.channels.awaitClose
import kotlinx.coroutines.flow.Flow
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.flow.callbackFlow
import platform.Network.*
import kotlin.coroutines.resume
import kotlin.coroutines.suspendCoroutine

class IosConnectivityObserver : ConnectivityObserver {

    private val monitor: nw_path_monitor_t = nw_path_monitor_create()
    private val _status = MutableStateFlow(ConnectionStatus.UNAVAILABLE)
    override val status: StateFlow<ConnectionStatus> = _status.asStateFlow()

    override fun start() {
        nw_path_monitor_set_update_handler(monitor) { path ->
            val status = when (nw_path_get_status(path)) {
                nw_path_status_satisfied -> ConnectionStatus.AVAILABLE
                nw_path_status_satisfiable -> ConnectionStatus.LOSING
                else -> ConnectionStatus.UNAVAILABLE
            }
            _status.value = status
        }
        nw_path_monitor_start(monitor)
    }

    override fun stop() {
        nw_path_monitor_cancel(monitor)
    }
}
