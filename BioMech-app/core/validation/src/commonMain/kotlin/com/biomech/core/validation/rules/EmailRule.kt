package com.biomech.core.validation.rules

import com.biomech.core.validation.ValidationResult
import com.biomech.core.validation.ValidationRule

class EmailRule(private val message: String = "Invalid email address") : ValidationRule<String> {
    override fun validate(value: String): ValidationResult {
        val emailRegex = Regex("^[A-Za-z0-9+_.-]+@[A-Za-z0-9.-]+\\.[A-Za-z]{2,}$")
        return if (emailRegex.matches(value.trim())) ValidationResult.Valid
        else ValidationResult.Invalid(listOf(message))
    }
}
