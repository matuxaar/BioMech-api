package com.biomech.core.file

import com.biomech.core.common.PlatformContext
import platform.Foundation.NSCachesDirectory
import platform.Foundation.NSDocumentDirectory
import platform.Foundation.NSSearchPathForDirectoriesInDomains
import platform.Foundation.NSTemporaryDirectory
import platform.Foundation.NSUserDomainMask

actual class FileManager {
    actual fun getCacheDir(): String = NSSearchPathForDirectoriesInDomains(
        NSCachesDirectory, NSUserDomainMask, true
    ).firstOrNull() as? String ?: ""

    actual fun getDocumentsDir(): String = NSSearchPathForDirectoriesInDomains(
        NSDocumentDirectory, NSUserDomainMask, true
    ).firstOrNull() as? String ?: ""

    actual fun getTempDir(): String = NSTemporaryDirectory() ?: ""
}

actual fun createFileManager(context: PlatformContext): FileManager = FileManager()
