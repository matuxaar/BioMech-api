package com.biomech.app.di

import com.biomech.core.ble.AndroidBleManager
import com.biomech.core.ble.BleManager
import org.koin.dsl.module

actual val platformModule = module {
    single<BleManager> { AndroidBleManager() }
}
