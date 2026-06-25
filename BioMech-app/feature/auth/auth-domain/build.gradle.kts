plugins {
    id("com.android.kotlin.multiplatform.library")
    kotlin("multiplatform")
    id("com.biomech.gradle.kmp-ios-conventions")
}

kotlin {
    android {
        namespace = "com.biomech.feature.auth.domain"
        compileSdk = 36
        minSdk = 26
    }


    sourceSets {
        commonMain.dependencies {
            implementation(project(":core:common"))
        }
        commonTest.dependencies {
            implementation(libs.kotlin.test.unit)
            implementation(libs.kotlinx.coroutines.test)
        }
    }
}
