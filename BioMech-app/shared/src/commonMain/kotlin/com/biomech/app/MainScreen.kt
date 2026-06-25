package com.biomech.app

import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.padding
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Modifier
import androidx.compose.ui.unit.dp
import com.biomech.core.navigation.LocalNavigator
import com.biomech.core.navigation.Screen
import com.biomech.domain.model.Device
import com.biomech.domain.model.DeviceType
import com.biomech.feature.devices.AddDeviceBottomSheet
import com.biomech.feature.devices.DevicesAction
import com.biomech.feature.devices.DevicesViewModel
import com.biomech.feature.home.HomeAction
import com.biomech.feature.home.HomeViewModel
import com.biomech.feature.home.HomeScreen
import com.biomech.feature.profile.ProfileAction
import com.biomech.feature.profile.ProfileEvent
import com.biomech.feature.profile.ProfileScreen
import com.biomech.feature.profile.ProfileViewModel
import com.biomech.feature.settings.SettingsAction
import com.biomech.feature.settings.SettingsEvent
import com.biomech.feature.settings.SettingsScreen
import com.biomech.feature.settings.SettingsViewModel
import kotlinx.coroutines.launch
import org.koin.compose.koinInject

enum class BottomTab { HOME, PROFILE, SETTINGS }

private val offlineDevices = listOf(
    Device(id = "dev-1", name = "MyoBand Pro", type = DeviceType.SENSOR, hwVersion = "2.1.0"),
    Device(id = "dev-2", name = "NeuroFlex", type = DeviceType.SENSOR, hwVersion = "1.4.2"),
    Device(id = "dev-3", name = "BioPulse", type = DeviceType.SENSOR, hwVersion = "3.0.1"),
)

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun MainScreen(isOffline: Boolean = false) {
    val navigator = LocalNavigator.current
    var selectedTab by remember { mutableStateOf(BottomTab.HOME) }
    var showAddDeviceSheet by remember { mutableStateOf(false) }
    var email by remember { mutableStateOf("user@example.com") }

    val homeViewModel: HomeViewModel = koinInject()
    val homeState by homeViewModel.state.collectAsState()

    val profileViewModel: ProfileViewModel = koinInject()
    val profileState by profileViewModel.state.collectAsState()
    val profileScope = rememberCoroutineScope()

    val settingsViewModel: SettingsViewModel = koinInject()
    val settingsState by settingsViewModel.state.collectAsState()
    val settingsScope = rememberCoroutineScope()

    val devicesViewModel: DevicesViewModel = koinInject()
    val devicesState by devicesViewModel.state.collectAsState()

    LaunchedEffect(Unit) {
        if (!isOffline) {
            homeViewModel.dispatch(HomeAction.LoadDevices)
            profileViewModel.dispatch(ProfileAction.LoadProfile)
        }
        profileViewModel.dispatch(ProfileAction.SetEmail(email))
    }

    LaunchedEffect(Unit) {
        launch {
            profileViewModel.event.collect { event ->
                when (event) {
                    ProfileEvent.NavigateToLogin -> navigator.navigateAndClear(Screen.Login)
                }
            }
        }
        launch {
            settingsViewModel.event.collect { event ->
                when (event) {
                    SettingsEvent.NavigateToLogin -> navigator.navigateAndClear(Screen.Login)
                }
            }
        }
    }

    if (showAddDeviceSheet) {
        ModalBottomSheet(
            onDismissRequest = { showAddDeviceSheet = false }
        ) {
            AddDeviceBottomSheet(
                scannedDevices = devicesState.scannedDevices,
                isScanning = devicesState.isScanning,
                onStartScan = { devicesViewModel.dispatch(DevicesAction.StartScan) },
                onStopScan = { devicesViewModel.dispatch(DevicesAction.StopScan) },
                onConnect = { deviceId ->
                    devicesViewModel.dispatch(DevicesAction.Connect(deviceId))
                    showAddDeviceSheet = false
                },
                onDismiss = { showAddDeviceSheet = false },
            )
        }
    }

    Scaffold(
        topBar = {
            if (isOffline) {
                Surface(
                    color = MaterialTheme.colorScheme.errorContainer,
                    modifier = Modifier.fillMaxWidth(),
                ) {
                    Text(
                        text = "Offline mode — server unavailable",
                        modifier = Modifier.padding(8.dp),
                        style = MaterialTheme.typography.bodySmall,
                        color = MaterialTheme.colorScheme.onErrorContainer,
                    )
                }
            }
        },
        bottomBar = {
            NavigationBar {
                NavigationBarItem(
                    selected = selectedTab == BottomTab.HOME,
                    onClick = { selectedTab = BottomTab.HOME },
                    icon = { Text("🏠") },
                    label = { Text("Home") },
                )
                NavigationBarItem(
                    selected = selectedTab == BottomTab.PROFILE,
                    onClick = { selectedTab = BottomTab.PROFILE },
                    icon = { Text("👤") },
                    label = { Text("Profile") },
                )
                NavigationBarItem(
                    selected = selectedTab == BottomTab.SETTINGS,
                    onClick = { selectedTab = BottomTab.SETTINGS },
                    icon = { Text("⚙️") },
                    label = { Text("Settings") },
                )
            }
        }
    ) { padding ->
        when (selectedTab) {
            BottomTab.HOME -> {
                HomeScreen(
                    devices = if (isOffline) offlineDevices else homeState.devices,
                    onAddDevice = { showAddDeviceSheet = true },
                    onDeviceClick = {
                        navigator.navigateTo(Screen.Dashboard)
                    },
                )
            }
            BottomTab.PROFILE -> {
                ProfileScreen(
                    email = if (isOffline) "demo@biomech.app" else profileState.email,
                    onLogout = {
                        profileViewModel.dispatch(ProfileAction.Logout)
                    },
                )
            }
            BottomTab.SETTINGS -> {
                SettingsScreen(
                    serverUrl = settingsState.serverUrl,
                    onServerUrlChange = { settingsViewModel.dispatch(SettingsAction.UpdateServerUrl(it)) },
                    onLogout = {
                        settingsViewModel.dispatch(SettingsAction.Logout)
                    },
                )
            }
        }
    }
}
