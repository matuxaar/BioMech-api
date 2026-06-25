plugins {
    id("com.android.kotlin.multiplatform.library")
    kotlin("multiplatform")
    id("com.biomech.gradle.kmp-ios-conventions")
}

kotlin {
    android {
        namespace = "com.biomech.core.crash"
        compileSdk = 36
        minSdk = 26
    }


    compilerOptions {
        freeCompilerArgs.add("-Xexpect-actual-classes")
    }

    sourceSets {
        commonMain.dependencies {
            implementation(project(":core:common"))
            implementation(project(":core:logger"))
        }
    }
}
