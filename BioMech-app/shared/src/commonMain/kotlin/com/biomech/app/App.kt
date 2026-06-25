package com.biomech.app

import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.Surface
import androidx.compose.runtime.*
import androidx.compose.ui.Modifier
import com.biomech.core.base.theme.BioMechTheme
import com.biomech.core.navigation.LocalNavigator
import com.biomech.core.navigation.Navigator
import com.biomech.core.navigation.Screen
import com.biomech.core.network.ApiConfig
import com.biomech.domain.repository.AuthRepository
import com.biomech.feature.auth.LoginAction
import com.biomech.feature.auth.LoginEvent
import com.biomech.feature.auth.LoginScreen
import com.biomech.feature.auth.LoginViewModel
import com.biomech.feature.dashboard.DashboardAction
import com.biomech.feature.dashboard.DashboardScreen
import com.biomech.feature.dashboard.DashboardViewModel
import com.biomech.feature.training.TrainingAction
import com.biomech.feature.training.TrainingScreen
import com.biomech.feature.training.TrainingViewModel
import io.ktor.client.*
import io.ktor.client.plugins.*
import io.ktor.client.request.*
import kotlinx.coroutines.withTimeout
import org.koin.compose.koinInject

private val healthClient by lazy {
    HttpClient {
        install(HttpTimeout) {
            requestTimeoutMillis = 3_000
            connectTimeoutMillis = 3_000
        }
    }
}

private suspend fun isApiAvailable(): Boolean {
    return try {
        withTimeout(4_000) {
            healthClient.get(ApiConfig.baseUrl)
            true
        }
    } catch (_: Exception) {
        false
    }
}

@Composable
fun App() {
    val navigator = remember { Navigator() }
    val authRepo: AuthRepository = koinInject()
    var startupState by remember { mutableStateOf<AppStartupState>(AppStartupState.Loading) }

    LaunchedEffect(Unit) {
        val apiUp = isApiAvailable()
        if (!apiUp) {
            startupState = AppStartupState.Offline
            navigator.navigateAndClear(Screen.Main)
            return@LaunchedEffect
        }
        val token = authRepo.getToken()
        if (token != null) {
            navigator.navigateAndClear(Screen.Main)
        }
        startupState = AppStartupState.Ready
    }

    if (startupState == AppStartupState.Loading) return

    BioMechTheme {
        CompositionLocalProvider(LocalNavigator provides navigator) {
            Surface(
                modifier = Modifier.fillMaxSize(),
                color = MaterialTheme.colorScheme.background
            ) {
                when (navigator.currentScreen) {
                    Screen.Login -> {
                        val viewModel: LoginViewModel = koinInject()
                        val state by viewModel.state.collectAsState()

                        LaunchedEffect(Unit) {
                            viewModel.event.collect { event ->
                                when (event) {
                                    LoginEvent.NavigateToMain -> navigator.navigateAndClear(Screen.Main)
                                }
                            }
                        }

                        LoginScreen(
                            isLoading = state.isLoading,
                            error = state.error,
                            onLogin = { email, password ->
                                viewModel.dispatch(LoginAction.Login(email, password))
                            },
                            onRegister = { email, password ->
                                viewModel.dispatch(LoginAction.Register(email, password))
                            },
                            onSkip = {
                                navigator.navigateAndClear(Screen.Main)
                            },
                        )
                    }

                    Screen.Main -> {
                        MainScreen(isOffline = startupState == AppStartupState.Offline)
                    }

                    Screen.Dashboard -> {
                        val viewModel: DashboardViewModel = koinInject()
                        val state by viewModel.state.collectAsState()

                        DashboardScreen(
                            deviceConnected = state.deviceConnected,
                            emgData = state.emgData,
                            onStartRecording = { viewModel.dispatch(DashboardAction.StartRecording) },
                            onStopRecording = { viewModel.dispatch(DashboardAction.StopRecording) },
                            onNavigateToTraining = { navigator.navigateTo(Screen.Training) },
                        )
                    }

                    Screen.Training -> {
                        val viewModel: TrainingViewModel = koinInject()
                        val state by viewModel.state.collectAsState()

                        TrainingScreen(
                            jobs = state.jobs,
                            sessions = state.sessions,
                            selectedSessionIds = state.selectedSessionIds,
                            isCreating = state.isCreating,
                            error = state.error,
                            onToggleSession = { viewModel.dispatch(TrainingAction.ToggleSession(it)) },
                            onStartTraining = { viewModel.dispatch(TrainingAction.StartTraining) },
                        )
                    }
                }
            }
        }
    }
}

private sealed interface AppStartupState {
    data object Loading : AppStartupState
    data object Ready : AppStartupState
    data object Offline : AppStartupState
}
