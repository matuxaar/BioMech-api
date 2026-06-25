plugins {
    `kotlin-dsl`
    `java-gradle-plugin`
}

repositories {
    google()
    mavenCentral()
    gradlePluginPortal()
}

dependencies {
    implementation("com.android.tools.build:gradle:9.2.1")
    implementation("org.jetbrains.kotlin:kotlin-gradle-plugin:2.4.0")
    implementation("org.jetbrains.compose:compose-gradle-plugin:1.11.1")
}

gradlePlugin {
    plugins {
        register("kmpIosConventions") {
            id = "com.biomech.gradle.kmp-ios-conventions"
            implementationClass = "com.biomech.gradle.KmpIosConventionsPlugin"
        }
    }
}
