rootProject.name = "BioMechApp"
includeBuild("build-logic")

pluginManagement {
    plugins {
        kotlin("plugin.compose") version "2.4.0"
        kotlin("plugin.serialization") version "2.4.0"
        id("com.google.devtools.ksp") version "2.3.9"
    }
    repositories {
        google()
        mavenCentral()
        gradlePluginPortal()
        maven("https://maven.pkg.jetbrains.space/public/p/kotlin/bootstrap")
    }
}

include(":shared")
include(":androidApp")

include(":core:common")
include(":core:base")
include(":core:component")
include(":core:resource")
include(":core:mvi")
include(":core:storage")
include(":core:navigation")
include(":core:analytics")
include(":core:permission")
include(":core:datetime")
include(":core:logger")
include(":core:connectivity")
include(":core:validation")
include(":core:buildconfig")
include(":core:file")
include(":core:deeplink")
include(":core:ui-system")
include(":core:di")
include(":core:image")
include(":core:share")
include(":core:update")
include(":core:notifications")
include(":core:crash")
include(":core:network")
include(":core:firebase")
include(":core:ble")
include(":core:database")

include(":domain")

include(":feature:auth:auth-domain")
include(":feature:auth:auth-data")
include(":feature:auth:auth-api")
include(":feature:auth:auth-ui")
include(":feature:dashboard:dashboard-api")
include(":feature:dashboard:dashboard-ui")
include(":feature:home:home-ui")
include(":feature:profile:profile-ui")
include(":feature:devices:devices-domain")
include(":feature:devices:devices-data")
include(":feature:devices:devices-api")
include(":feature:devices:devices-ui")
include(":feature:training:training-domain")
include(":feature:training:training-data")
include(":feature:training:training-api")
include(":feature:training:training-ui")
include(":feature:settings:settings-api")
include(":feature:settings:settings-ui")
