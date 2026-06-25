package com.biomech.core.navigation

import androidx.compose.runtime.Composable
import androidx.compose.runtime.getValue
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.setValue

class Navigator {
    var currentScreen: Screen by mutableStateOf(Screen.Login)
        private set

    private val backStack = mutableListOf<Screen>()

    fun navigateTo(screen: Screen) {
        backStack.add(currentScreen)
        currentScreen = screen
    }

    fun goBack() {
        if (backStack.isNotEmpty()) {
            currentScreen = backStack.removeLast()
        }
    }

    fun navigateAndClear(screen: Screen) {
        backStack.clear()
        currentScreen = screen
    }
}

val LocalNavigator = androidx.compose.runtime.staticCompositionLocalOf { Navigator() }
