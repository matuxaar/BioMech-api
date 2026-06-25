package com.biomech.core.validation

class Validator<T> {
    private val rules = mutableListOf<ValidationRule<T>>()

    fun addRule(rule: ValidationRule<T>) {
        rules.add(rule)
    }

    fun addRules(vararg rules: ValidationRule<T>) {
        this.rules.addAll(rules)
    }

    fun validate(value: T): ValidationResult {
        val errors = rules.mapNotNull { rule ->
            val result = rule.validate(value)
            when (result) {
                is ValidationResult.Invalid -> result.errors
                is ValidationResult.Valid -> null
            }
        }.flatten()

        return if (errors.isEmpty()) ValidationResult.Valid
        else ValidationResult.Invalid(errors)
    }
}
