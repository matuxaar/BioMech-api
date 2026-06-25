plugins {
    id("com.android.kotlin.multiplatform.library")
    kotlin("multiplatform")
    id("com.biomech.gradle.kmp-ios-conventions")
}

kotlin {
    android {
        namespace = "com.biomech.feature.devices.data"
        compileSdk = 36
        minSdk = 26
    }


    sourceSets {
        commonMain.dependencies {
            api(project(":feature:devices:devices-api"))
            api(project(":feature:devices:devices-domain"))
            implementation(project(":core:network"))
            implementation(project(":core:common"))
        }
        androidMain.dependencies {
            implementation(project(":core:database"))
        }
    }
}
