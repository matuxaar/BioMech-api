plugins {
    id("com.android.kotlin.multiplatform.library")
    kotlin("multiplatform")
    id("com.biomech.gradle.kmp-ios-conventions")
}

kotlin {
    android {
        namespace = "com.biomech.core.common"
        compileSdk = 36
        minSdk = 26
    }


    compilerOptions {
        freeCompilerArgs.add("-Xexpect-actual-classes")
    }

    sourceSets {
        commonMain.dependencies {
            implementation(libs.kotlinx.coroutines.core)
        }
    }
}
