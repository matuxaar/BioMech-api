plugins {
    id("com.android.kotlin.multiplatform.library")
    kotlin("multiplatform")
    id("com.biomech.gradle.kmp-ios-conventions")
}

kotlin {
    android {
        namespace = "com.biomech.core.di"
        compileSdk = 36
        minSdk = 26
    }


    compilerOptions {
        freeCompilerArgs.add("-Xexpect-actual-classes")
    }

    sourceSets {
        commonMain.dependencies {
            implementation(project(":core:common"))
            implementation(project(":core:mvi"))
            implementation(project(":core:network"))
            implementation(project(":core:ble"))
            implementation(project(":core:storage"))
            implementation(libs.room.runtime)
            implementation(project(":domain"))
            implementation(project(":feature:auth:auth-domain"))
            implementation(project(":feature:auth:auth-api"))
            implementation(project(":feature:auth:auth-data"))
            implementation(project(":feature:auth:auth-ui"))
            implementation(project(":feature:home:home-ui"))
            implementation(project(":feature:profile:profile-ui"))
            implementation(project(":feature:devices:devices-domain"))
            implementation(project(":feature:devices:devices-api"))
            implementation(project(":feature:devices:devices-data"))
            implementation(project(":feature:devices:devices-ui"))
            implementation(project(":feature:training:training-domain"))
            implementation(project(":feature:training:training-api"))
            implementation(project(":feature:training:training-data"))
            implementation(project(":feature:training:training-ui"))
            implementation(project(":feature:settings:settings-ui"))
            implementation(libs.koin.core)
        }
        androidMain.dependencies {
            implementation("io.insert-koin:koin-android:4.2.1")
            implementation(project(":core:database"))
        }
    }
}
