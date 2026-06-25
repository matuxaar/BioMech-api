plugins {
    id("com.android.kotlin.multiplatform.library")
    kotlin("multiplatform")
    id("com.biomech.gradle.kmp-ios-conventions")
    kotlin("plugin.serialization")
}

kotlin {
    android {
        namespace = "com.biomech.feature.devices.api"
        compileSdk = 36
        minSdk = 26
    }


    sourceSets {
        commonMain.dependencies {
            implementation(project(":core:network"))
            implementation(libs.ktor.client.core)
            implementation(libs.kotlinx.serialization.json)
        }
    }
}
