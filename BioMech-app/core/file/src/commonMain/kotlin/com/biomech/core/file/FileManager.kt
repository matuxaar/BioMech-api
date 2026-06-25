package com.biomech.core.file

import com.biomech.core.common.PlatformContext

expect class FileManager {
    fun getCacheDir(): String
    fun getDocumentsDir(): String
    fun getTempDir(): String
}

expect fun createFileManager(context: PlatformContext): FileManager
