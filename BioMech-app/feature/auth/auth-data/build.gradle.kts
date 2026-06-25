plugins {
    id("com.android.kotlin.multiplatform.library")
    kotlin("multiplatform")
    id("com.biomech.gradle.kmp-ios-conventions")
}

kotlin {
    android {
        namespace = "com.biomech.feature.auth.data"
        compileSdk = 36
        minSdk = 26
    }


    sourceSets {
        commonMain.dependencies {
            api(project(":feature:auth:auth-api"))
            api(project(":feature:auth:auth-domain"))
            implementation(project(":core:network"))
            implementation(project(":core:common"))
            implementation(project(":core:storage"))
        }
        commonTest.dependencies {
            implementation(libs.kotlin.test.unit)
            implementation(libs.kotlinx.coroutines.test)
        }
    }
}
