package com.biomech.core.permission

enum class Permission {
    BLUETOOTH,
    LOCATION,
    NOTIFICATIONS,
}

enum class PermissionStatus {
    GRANTED,
    DENIED,
    PERMANENTLY_DENIED,
}

interface PermissionManager {
    suspend fun checkPermission(permission: Permission): PermissionStatus
    suspend fun requestPermission(permission: Permission): PermissionStatus
}
