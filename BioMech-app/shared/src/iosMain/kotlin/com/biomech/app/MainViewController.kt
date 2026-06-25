package com.biomech.app

import androidx.compose.ui.window.ComposeUIViewController
import com.biomech.core.common.PlatformContext
import com.biomech.core.di.initKoin
import platform.UIKit.UIViewController

fun MainViewController(): UIViewController {
    initKoin(PlatformContext())
    return ComposeUIViewController(
        configure = {
            enforceStrictPlistSanityCheck = false
        },
    ) { App() }
}
