plugins {
    id("com.android.kotlin.multiplatform.library")
    kotlin("multiplatform")
    id("com.biomech.gradle.kmp-ios-conventions")
}

kotlin {
    android {
        namespace = "com.biomech.core.firebase"
        compileSdk = 36
        minSdk = 26
    }


    sourceSets {
        commonMain.dependencies {
        }
    }
}
