package com.biomech.core.di

import com.biomech.core.ble.BleManager
import com.biomech.core.ble.IosBleManager
import com.biomech.core.common.PlatformContext
import com.biomech.core.network.api.EMGApi
import com.biomech.core.network.repository.AuthRepositoryImpl
import com.biomech.core.network.repository.DeviceRepositoryImpl
import com.biomech.core.network.repository.EMGRepositoryImpl
import com.biomech.core.network.repository.TrainingRepositoryImpl
import com.biomech.domain.repository.AuthRepository
import com.biomech.domain.repository.DeviceRepository
import com.biomech.domain.repository.EMGRepository
import com.biomech.domain.repository.TrainingRepository
import org.koin.core.module.Module
import org.koin.dsl.module

actual val databaseModule: Module = module {
}

actual val repositoriesModule: Module = module {
    single<AuthRepository> { AuthRepositoryImpl(get(), get()) }
    single<DeviceRepository> { DeviceRepositoryImpl(get()) }
    single<TrainingRepository> { TrainingRepositoryImpl(get()) }
    single<EMGRepository> { EMGRepositoryImpl(get()) }
}

actual val platformModule: Module = module {
    single<BleManager> { IosBleManager() }
}

actual fun initKoin(context: PlatformContext) {
    com.biomech.core.network.ApiConfig.baseUrl = "http://localhost:8080/api/v1"
    org.koin.core.context.startKoin {
        modules(
            module { single { context } },
            platformModule,
            storageModule,
            databaseModule,
            apisModule,
            repositoriesModule,
            useCasesModule,
            viewModelsModule,
        )
    }
}
