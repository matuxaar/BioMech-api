package com.biomech.core.storage

import com.biomech.core.common.PlatformContext

expect fun createKeyValueStorage(context: PlatformContext): KeyValueStorage
