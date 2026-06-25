package com.biomech.core.validation

sealed class ValidationResult {
    data object Valid : ValidationResult()
    data class Invalid(val errors: List<String>) : ValidationResult() {
        val firstError: String get() = errors.first()
    }
}

fun ValidationResult.isValid(): Boolean = this is ValidationResult.Valid
fun ValidationResult.isInvalid(): Boolean = this is ValidationResult.Invalid
