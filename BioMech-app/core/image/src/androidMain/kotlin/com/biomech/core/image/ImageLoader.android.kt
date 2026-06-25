package com.biomech.core.image

import coil3.ImageLoader
import com.biomech.core.common.PlatformContext

actual fun createImageLoader(context: PlatformContext): ImageLoader {
    return ImageLoader.Builder(context.androidContext).build()
}
