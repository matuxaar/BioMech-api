package com.biomech.core.navigation

sealed class Screen(val route: String) {
    data object Login : Screen("login")
    data object Main : Screen("main")
    data object Training : Screen("training")
}
