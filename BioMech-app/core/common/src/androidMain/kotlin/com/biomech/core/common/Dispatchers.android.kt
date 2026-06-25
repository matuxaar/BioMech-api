package com.biomech.core.common

import kotlinx.coroutines.CoroutineDispatcher
import kotlinx.coroutines.Dispatchers

actual val ioDispatcher: CoroutineDispatcher = Dispatchers.IO
actual val mainDispatcher: CoroutineDispatcher = Dispatchers.Main
