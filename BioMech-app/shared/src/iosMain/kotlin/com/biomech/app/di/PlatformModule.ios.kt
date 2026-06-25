package com.biomech.app.di

import com.biomech.core.ble.BleManager
import com.biomech.core.ble.IosBleManager
import org.koin.dsl.module

actual val platformModule = module {
    single<BleManager> { IosBleManager() }
}
