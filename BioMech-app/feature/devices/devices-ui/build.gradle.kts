plugins {
    id("com.android.kotlin.multiplatform.library")
    kotlin("multiplatform")
    id("com.biomech.gradle.kmp-ios-conventions")
    id("org.jetbrains.compose")
    id("org.jetbrains.kotlin.plugin.compose")
}

kotlin {
    android {
        namespace = "com.biomech.feature.devices.ui"
        compileSdk = 36
        minSdk = 26
    }


    sourceSets {
        commonMain.dependencies {
            implementation(project(":feature:devices:devices-domain"))
            implementation(project(":core:mvi"))
            implementation(project(":core:base"))
            implementation(project(":core:ble"))
            implementation(libs.kotlinx.coroutines.core)
            implementation(compose.runtime)
            implementation(compose.foundation)
            implementation(compose.material3)
            implementation(compose.ui)
        }
    }
}
