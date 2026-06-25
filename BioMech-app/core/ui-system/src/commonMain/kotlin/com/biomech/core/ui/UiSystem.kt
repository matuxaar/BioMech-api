package com.biomech.core.ui

import androidx.compose.runtime.compositionLocalOf
import com.biomech.core.common.PlatformContext

data class WindowInsets(
    val top: Float,
    val bottom: Float,
    val left: Float,
    val right: Float
)

expect class UiSystem {
    val safeAreaInsets: WindowInsets
    fun hideKeyboard()
}

expect fun createUiSystem(context: PlatformContext): UiSystem

val LocalUiSystem = compositionLocalOf<UiSystem> {
    error("UiSystem not provided")
}
