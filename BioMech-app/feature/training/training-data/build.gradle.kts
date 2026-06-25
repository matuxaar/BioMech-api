plugins {
    id("com.android.kotlin.multiplatform.library")
    kotlin("multiplatform")
    id("com.biomech.gradle.kmp-ios-conventions")
}

kotlin {
    android {
        namespace = "com.biomech.feature.training.data"
        compileSdk = 36
        minSdk = 26
    }


    sourceSets {
        commonMain.dependencies {
            api(project(":feature:training:training-api"))
            api(project(":feature:training:training-domain"))
            implementation(project(":core:network"))
            implementation(project(":core:common"))
        }
        androidMain.dependencies {
            implementation(project(":core:database"))
        }
    }
}
