package com.biomech.core.connectivity

import android.content.Context
import android.net.ConnectivityManager
import android.net.Network
import android.net.NetworkCapabilities
import android.net.NetworkRequest
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow

class AndroidConnectivityObserver(context: Context) : ConnectivityObserver {

    private val connectivityManager =
        context.getSystemService(Context.CONNECTIVITY_SERVICE) as ConnectivityManager

    private val _status = MutableStateFlow(ConnectionStatus.UNAVAILABLE)
    override val status: StateFlow<ConnectionStatus> = _status.asStateFlow()

    private val networkCallback = object : ConnectivityManager.NetworkCallback() {
        override fun onAvailable(network: Network) {
            _status.value = ConnectionStatus.AVAILABLE
        }

        override fun onLost(network: Network) {
            _status.value = ConnectionStatus.UNAVAILABLE
        }

        override fun onCapabilitiesChanged(network: Network, capabilities: NetworkCapabilities) {
            val hasInternet = capabilities.hasCapability(
                NetworkCapabilities.NET_CAPABILITY_INTERNET
            )
            _status.value = if (hasInternet) ConnectionStatus.AVAILABLE
            else ConnectionStatus.UNAVAILABLE
        }
    }

    override fun start() {
        val request = NetworkRequest.Builder()
            .addCapability(NetworkCapabilities.NET_CAPABILITY_INTERNET)
            .build()
        connectivityManager.registerNetworkCallback(request, networkCallback)
        checkCurrentStatus()
    }

    override fun stop() {
        connectivityManager.unregisterNetworkCallback(networkCallback)
    }

    private fun checkCurrentStatus() {
        val network = connectivityManager.activeNetwork
        val capabilities = connectivityManager.getNetworkCapabilities(network)
        val hasInternet = capabilities?.hasCapability(
            NetworkCapabilities.NET_CAPABILITY_INTERNET
        ) ?: false
        _status.value = if (hasInternet) ConnectionStatus.AVAILABLE
        else ConnectionStatus.UNAVAILABLE
    }
}
