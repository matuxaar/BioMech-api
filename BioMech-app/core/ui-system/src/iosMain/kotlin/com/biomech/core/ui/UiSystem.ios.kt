package com.biomech.core.ui

import com.biomech.core.common.PlatformContext

actual class UiSystem {
    actual val safeAreaInsets: WindowInsets
        get() = WindowInsets(0f, 0f, 0f, 0f)

    actual fun hideKeyboard() {
    }
}

actual fun createUiSystem(context: PlatformContext): UiSystem = UiSystem()
