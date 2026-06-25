package com.biomech.core.storage

import com.biomech.core.common.PlatformContext

actual fun createKeyValueStorage(context: PlatformContext): KeyValueStorage {
    return AndroidKeyValueStorage(context.androidContext)
}
