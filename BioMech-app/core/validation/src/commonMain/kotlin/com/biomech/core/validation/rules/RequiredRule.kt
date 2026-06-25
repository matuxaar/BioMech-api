package com.biomech.core.validation.rules

import com.biomech.core.validation.ValidationResult
import com.biomech.core.validation.ValidationRule

class RequiredRule(private val message: String = "This field is required") : ValidationRule<String> {
    override fun validate(value: String): ValidationResult {
        return if (value.isNotBlank()) ValidationResult.Valid
        else ValidationResult.Invalid(listOf(message))
    }
}
