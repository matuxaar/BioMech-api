plugins {
    id("com.android.kotlin.multiplatform.library")
    kotlin("multiplatform")
    id("com.biomech.gradle.kmp-ios-conventions")
}

kotlin {
    android {
        namespace = "com.biomech.core.network"
        compileSdk = 36
        minSdk = 26
    }


    sourceSets {
        commonMain.dependencies {
            implementation(project(":core:common"))
            implementation(project(":domain"))
            api(libs.ktor.client.core)
            api(libs.ktor.client.content.negotiation)
            api(libs.ktor.serialization.json)
            api(libs.ktor.client.auth)
            implementation(libs.ktor.client.logging)
            implementation(libs.kotlinx.serialization.json)
        }
        androidMain.dependencies {
            implementation(libs.ktor.client.okhttp)
            implementation(project(":core:database"))
        }
        iosMain.dependencies {
            implementation(libs.ktor.client.darwin)
        }
    }
}
