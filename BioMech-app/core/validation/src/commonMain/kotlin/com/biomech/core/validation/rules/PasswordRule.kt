package com.biomech.core.validation.rules

import com.biomech.core.validation.ValidationResult
import com.biomech.core.validation.ValidationRule

class PasswordRule(private val message: String = "Password must contain at least 8 characters, one uppercase letter, one lowercase letter, and one digit") : ValidationRule<String> {
    override fun validate(value: String): ValidationResult {
        val errors = mutableListOf<String>()
        if (value.length < 8) errors.add("At least 8 characters required")
        if (!value.any { it.isUpperCase() }) errors.add("Must contain an uppercase letter")
        if (!value.any { it.isLowerCase() }) errors.add("Must contain a lowercase letter")
        if (!value.any { it.isDigit() }) errors.add("Must contain a digit")
        return if (errors.isEmpty()) ValidationResult.Valid
        else ValidationResult.Invalid(listOf(message))
    }
}
