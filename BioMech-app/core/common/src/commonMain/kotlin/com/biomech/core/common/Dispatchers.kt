package com.biomech.core.common

import kotlinx.coroutines.CoroutineDispatcher

expect val ioDispatcher: CoroutineDispatcher
expect val mainDispatcher: CoroutineDispatcher
