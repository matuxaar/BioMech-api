package com.biomech.core.validation

fun interface ValidationRule<T> {
    fun validate(value: T): ValidationResult
}
