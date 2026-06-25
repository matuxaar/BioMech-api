package com.biomech.core.share

import com.biomech.core.common.PlatformContext

data class ShareData(
    val text: String? = null,
    val title: String? = null,
)

expect class ShareManager {
    fun share(data: ShareData)
}

expect fun createShareManager(context: PlatformContext): ShareManager
