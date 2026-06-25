package com.biomech.core.image

import coil3.ImageLoader
import coil3.PlatformContext
import com.biomech.core.common.PlatformContext as AppContext

actual fun createImageLoader(context: AppContext): ImageLoader {
    return ImageLoader(PlatformContext.INSTANCE)
}
