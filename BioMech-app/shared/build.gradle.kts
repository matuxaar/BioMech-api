plugins {
    id("com.android.kotlin.multiplatform.library")
    kotlin("multiplatform")
    id("org.jetbrains.compose")
    id("org.jetbrains.kotlin.plugin.compose")
}

kotlin {
    android {
        namespace = "com.biomech.shared"
        compileSdk = 36
        minSdk = 26
    }

    listOf(iosArm64(), iosSimulatorArm64()).forEach { target ->
        target.binaries.framework {
            baseName = "ComposeApp"
            isStatic = true
        }
    }

    sourceSets {
        commonMain.dependencies {
            implementation(compose.runtime)
            implementation(compose.foundation)
            implementation(compose.material3)
            implementation(compose.ui)

            implementation(project(":core:base"))
            implementation(project(":core:component"))
            implementation(project(":core:resource"))
            api(project(":core:common"))
            implementation(project(":core:mvi"))
            implementation(project(":core:storage"))
            implementation(project(":core:navigation"))
            implementation(project(":core:connectivity"))
            implementation(project(":core:validation"))
            implementation(project(":core:buildconfig"))
            implementation(project(":core:file"))
            implementation(project(":core:deeplink"))
            implementation(project(":core:ui-system"))
            api(project(":core:di"))
            implementation(project(":core:image"))
            implementation(project(":core:share"))
            implementation(project(":core:update"))
            implementation(project(":core:notifications"))
            implementation(project(":core:crash"))
            implementation(project(":core:network"))
            implementation(project(":core:ble"))
            implementation(project(":domain"))
            implementation(project(":feature:auth:auth-domain"))
            implementation(project(":feature:auth:auth-data"))
            implementation(project(":feature:auth:auth-api"))
            implementation(project(":feature:auth:auth-ui"))
            implementation(project(":feature:home:home-ui"))
            implementation(project(":feature:profile:profile-ui"))
            implementation(project(":feature:dashboard:dashboard-api"))
            implementation(project(":feature:dashboard:dashboard-ui"))
            implementation(project(":feature:devices:devices-domain"))
            implementation(project(":feature:devices:devices-data"))
            implementation(project(":feature:devices:devices-api"))
            implementation(project(":feature:devices:devices-ui"))
            implementation(project(":feature:training:training-domain"))
            implementation(project(":feature:training:training-data"))
            implementation(project(":feature:training:training-api"))
            implementation(project(":feature:training:training-ui"))
            implementation(project(":feature:settings:settings-api"))
            implementation(project(":feature:settings:settings-ui"))

            implementation(libs.koin.core)
            implementation(libs.koin.compose)
            implementation(libs.kotlinx.coroutines.core)
        }
    }
}
