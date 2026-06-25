package com.biomech.core.validation.rules

import com.biomech.core.validation.ValidationResult
import kotlin.test.Test
import kotlin.test.assertTrue

class RequiredRuleTest {

    private val rule = RequiredRule()

    @Test
    fun `non-blank string returns valid`() {
        assertTrue(rule.validate("hello") is ValidationResult.Valid)
    }

    @Test
    fun `empty string returns invalid`() {
        assertTrue(rule.validate("") is ValidationResult.Invalid)
    }

    @Test
    fun `blank string returns invalid`() {
        assertTrue(rule.validate("   ") is ValidationResult.Invalid)
    }
}
