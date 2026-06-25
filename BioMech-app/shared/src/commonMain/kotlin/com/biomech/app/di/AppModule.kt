package com.biomech.app.di

import com.biomech.core.ble.BleManager
import com.biomech.core.network.api.AuthApi
import com.biomech.core.network.api.DeviceApi
import com.biomech.core.network.api.EMGApi
import com.biomech.core.network.api.TrainingApi
import com.biomech.core.network.repository.AuthRepositoryImpl
import com.biomech.core.network.repository.DeviceRepositoryImpl
import com.biomech.core.network.repository.EMGRepositoryImpl
import com.biomech.core.network.repository.TrainingRepositoryImpl
import com.biomech.domain.repository.AuthRepository
import com.biomech.domain.repository.DeviceRepository
import com.biomech.domain.repository.EMGRepository
import com.biomech.domain.repository.TrainingRepository
import com.biomech.domain.usecase.LoginUseCase
import com.biomech.domain.usecase.StartEMGSessionUseCase
import com.biomech.feature.auth.LoginViewModel
import com.biomech.feature.dashboard.DashboardViewModel
import com.biomech.feature.devices.DevicesViewModel
import com.biomech.feature.home.HomeViewModel
import com.biomech.feature.profile.ProfileViewModel
import com.biomech.feature.settings.SettingsViewModel
import com.biomech.feature.training.TrainingViewModel
import org.koin.core.module.Module
import org.koin.dsl.module

expect val platformModule: Module

val sharedModule = module {
    single { AuthApi() }
    single { DeviceApi() }
    single { EMGApi() }
    single { TrainingApi() }

    single<AuthRepository> { AuthRepositoryImpl(get(), get()) }
    single<DeviceRepository> { DeviceRepositoryImpl(get()) }
    single<EMGRepository> { EMGRepositoryImpl(get()) }
    single<TrainingRepository> { TrainingRepositoryImpl(get()) }

    single<LoginUseCase> { LoginUseCase(get()) }
    single<StartEMGSessionUseCase> { StartEMGSessionUseCase(get()) }

    factory { LoginViewModel(get()) }
    factory { DashboardViewModel(get(), get()) }
    factory { DevicesViewModel(get()) }
    factory { HomeViewModel(get()) }
    factory { ProfileViewModel(get()) }
    factory { TrainingViewModel(get(), get()) }
    factory { SettingsViewModel(get()) }
}

val appModules = listOf(sharedModule, platformModule)
