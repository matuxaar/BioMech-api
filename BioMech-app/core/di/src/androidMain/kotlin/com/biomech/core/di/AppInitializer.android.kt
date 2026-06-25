package com.biomech.core.di

import com.biomech.core.ble.AndroidBleManager
import com.biomech.core.ble.BleManager
import com.biomech.core.common.PlatformContext
import com.biomech.core.database.AppDatabase
import com.biomech.core.database.createRoomDatabase
import com.biomech.core.network.repository.AndroidEMGRepositoryImpl
import com.biomech.core.network.repository.AndroidDeviceRepositoryImpl
import com.biomech.core.network.repository.AndroidTrainingRepositoryImpl
import com.biomech.core.network.repository.AuthRepositoryImpl
import com.biomech.domain.repository.AuthRepository
import com.biomech.domain.repository.DeviceRepository
import com.biomech.domain.repository.EMGRepository
import com.biomech.domain.repository.TrainingRepository
import org.koin.android.ext.koin.androidContext
import org.koin.core.module.Module
import org.koin.dsl.module

actual val databaseModule: Module = module {
    single { get<AppDatabase>().deviceDao() }
    single { get<AppDatabase>().emgSessionDao() }
    single { get<AppDatabase>().trainingJobDao() }
}

actual val repositoriesModule: Module = module {
    single<AuthRepository> { AuthRepositoryImpl(get(), get()) }
    single<DeviceRepository> { AndroidDeviceRepositoryImpl(get(), get()) }
    single<TrainingRepository> { AndroidTrainingRepositoryImpl(get(), get()) }
    single<EMGRepository> { AndroidEMGRepositoryImpl(get(), get()) }
}

actual val platformModule: Module = module {
    single<BleManager> { AndroidBleManager() }
    single<AppDatabase> { createRoomDatabase(get()) }
}

actual fun initKoin(context: PlatformContext) {
    com.biomech.core.network.ApiConfig.baseUrl = "http://10.0.2.2:8080/api/v1"
    org.koin.core.context.startKoin {
        androidContext(context.androidContext)
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
