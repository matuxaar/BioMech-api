plugins {
    id("com.android.kotlin.multiplatform.library")
    kotlin("multiplatform")
    id("com.biomech.gradle.kmp-ios-conventions")
}

kotlin {
    android {
        namespace = "com.biomech.core.permission"
        compileSdk = 36
        minSdk = 26
    }


    sourceSets {
        commonMain.dependencies {
            implementation(project(":core:mvi"))
        }
    }
}
