package com.biomech.core.storage

import platform.Foundation.NSUserDefaults

class IosKeyValueStorage : KeyValueStorage {

    private val defaults = NSUserDefaults.standardUserDefaults

    override fun getString(key: String, default: String): String =
        defaults.stringForKey(key) ?: default

    override fun putString(key: String, value: String) {
        defaults.setObject(value, forKey = key)
    }

    override fun getBoolean(key: String, default: Boolean): Boolean =
        if (defaults.objectForKey(key) != null) defaults.boolForKey(key) else default

    override fun putBoolean(key: String, value: Boolean) {
        defaults.setBool(value, forKey = key)
    }

    override fun getInt(key: String, default: Int): Int =
        if (defaults.objectForKey(key) != null) defaults.integerForKey(key).toInt() else default

    override fun putInt(key: String, value: Int) {
        defaults.setInteger(value.toLong(), forKey = key)
    }

    override fun getLong(key: String, default: Long): Long =
        if (defaults.objectForKey(key) != null) defaults.integerForKey(key) else default

    override fun putLong(key: String, value: Long) {
        defaults.setInteger(value, forKey = key)
    }

    override fun remove(key: String) {
        defaults.removeObjectForKey(key)
    }

    override fun clear() {
        val dict = defaults.dictionaryRepresentation()
        for (key in dict.keys) {
            defaults.removeObjectForKey(key as String)
        }
    }
}
