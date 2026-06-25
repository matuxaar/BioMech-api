package com.biomech.gradle

import org.gradle.api.Plugin
import org.gradle.api.Project
import org.jetbrains.kotlin.gradle.dsl.KotlinMultiplatformExtension

class KmpIosConventionsPlugin : Plugin<Project> {
    override fun apply(target: Project) {
        target.extensions.configure(KotlinMultiplatformExtension::class.java) {
            iosArm64()
            iosSimulatorArm64()
        }
    }
}
