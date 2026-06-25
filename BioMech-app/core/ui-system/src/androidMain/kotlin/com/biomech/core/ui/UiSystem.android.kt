package com.biomech.core.ui

import android.app.Activity
import android.content.Context
import android.view.inputmethod.InputMethodManager
import com.biomech.core.common.PlatformContext

actual class UiSystem(context: PlatformContext) {
    private val androidContext = context.androidContext

    actual val safeAreaInsets: WindowInsets
        get() {
            val res = androidContext.resources
            val statusBarId = res.getIdentifier("status_bar_height", "dimen", "android")
            val navBarId = res.getIdentifier("navigation_bar_height", "dimen", "android")
            val top = if (statusBarId > 0) res.getDimension(statusBarId) / res.displayMetrics.density else 0f
            val bottom = if (navBarId > 0) res.getDimension(navBarId) / res.displayMetrics.density else 0f
            return WindowInsets(top = top, bottom = bottom, left = 0f, right = 0f)
        }

    actual fun hideKeyboard() {
        val imm = androidContext.getSystemService(Context.INPUT_METHOD_SERVICE) as? InputMethodManager
        imm?.let {
            val view = (androidContext as? Activity)?.currentFocus
            view?.let { v -> it.hideSoftInputFromWindow(v.windowToken, 0) }
        }
    }
}

actual fun createUiSystem(context: PlatformContext): UiSystem = UiSystem(context)
