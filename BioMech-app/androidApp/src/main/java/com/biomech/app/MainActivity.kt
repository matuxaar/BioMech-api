package com.biomech.app

import android.os.Bundle
import androidx.activity.ComponentActivity
import androidx.activity.compose.setContent
import com.biomech.core.common.PlatformContext
import com.biomech.core.di.initKoin

class MainActivity : ComponentActivity() {
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        initKoin(PlatformContext(applicationContext))
        setContent {
            App()
        }
    }
}
