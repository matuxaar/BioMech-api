package com.biomech.core.validation.rules

import com.biomech.core.validation.ValidationResult
import com.biomech.core.validation.ValidationRule

class MatchRule(
    private val otherValue: () -> String,
    private val message: String = "Values do not match"
) : ValidationRule<String> {
    override fun validate(value: String): ValidationResult {
        return if (value == otherValue()) ValidationResult.Valid
        else ValidationResult.Invalid(listOf(message))
    }
}
