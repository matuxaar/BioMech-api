package com.biomech.core.validation

import com.biomech.core.validation.rules.EmailRule
import com.biomech.core.validation.rules.RequiredRule
import kotlin.test.Test
import kotlin.test.assertFalse
import kotlin.test.assertTrue

class ValidatorTest {

    @Test
    fun `all rules valid returns valid`() {
        val validator = Validator<String>()
        validator.addRules(RequiredRule(), EmailRule())

        val result = validator.validate("user@example.com")

        assertTrue(result is ValidationResult.Valid)
    }

    @Test
    fun `one rule fails returns invalid`() {
        val validator = Validator<String>()
        validator.addRules(RequiredRule(), EmailRule())

        val result = validator.validate("not-an-email")

        assertTrue(result is ValidationResult.Invalid)
        assertFalse((result as ValidationResult.Invalid).errors.isEmpty())
    }

    @Test
    fun `no rules always returns valid`() {
        val validator = Validator<String>()

        val result = validator.validate("anything")

        assertTrue(result is ValidationResult.Valid)
    }
}
