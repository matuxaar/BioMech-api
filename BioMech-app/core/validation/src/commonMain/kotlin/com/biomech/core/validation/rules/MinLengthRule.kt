package com.biomech.core.validation.rules

import com.biomech.core.validation.ValidationResult
import com.biomech.core.validation.ValidationRule

class MinLengthRule(
    private val min: Int,
    private val message: String = "Minimum $min characters required"
) : ValidationRule<String> {
    override fun validate(value: String): ValidationResult {
        return if (value.length >= min) ValidationResult.Valid
        else ValidationResult.Invalid(listOf(message))
    }
}
