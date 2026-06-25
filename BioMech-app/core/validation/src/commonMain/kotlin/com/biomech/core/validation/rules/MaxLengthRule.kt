package com.biomech.core.validation.rules

import com.biomech.core.validation.ValidationResult
import com.biomech.core.validation.ValidationRule

class MaxLengthRule(
    private val max: Int,
    private val message: String = "Maximum $max characters allowed"
) : ValidationRule<String> {
    override fun validate(value: String): ValidationResult {
        return if (value.length <= max) ValidationResult.Valid
        else ValidationResult.Invalid(listOf(message))
    }
}
