package com.biomech.core.file

import com.biomech.core.common.PlatformContext

actual class FileManager(context: PlatformContext) {
    private val androidContext = context.androidContext

    actual fun getCacheDir(): String = androidContext.cacheDir.absolutePath
    actual fun getDocumentsDir(): String = androidContext.filesDir.absolutePath
    actual fun getTempDir(): String = androidContext.cacheDir.resolve("temp").also {
        if (!it.exists()) it.mkdirs()
    }.absolutePath
}

actual fun createFileManager(context: PlatformContext): FileManager = FileManager(context)
