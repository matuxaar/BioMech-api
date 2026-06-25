package com.biomech.core.deeplink

import com.biomech.core.common.PlatformContext

expect class DeepLinkHandler {
    fun handleDeepLink(uri: String): DeepLink
}

expect fun createDeepLinkHandler(context: PlatformContext): DeepLinkHandler
