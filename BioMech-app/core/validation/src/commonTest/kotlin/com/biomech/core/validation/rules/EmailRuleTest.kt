package com.biomech.core.validation.rules

import com.biomech.core.validation.ValidationResult
import kotlin.test.Test
import kotlin.test.assertTrue

class EmailRuleTest {

    private val rule = EmailRule()

    @Test
    fun `valid email returns valid`() {
        assertTrue(rule.validate("user@example.com") is ValidationResult.Valid)
    }

    @Test
    fun `email with subdomain returns valid`() {
        assertTrue(rule.validate("user@sub.example.com") is ValidationResult.Valid)
    }

    @Test
    fun `email with plus tag returns valid`() {
        assertTrue(rule.validate("user+tag@example.com") is ValidationResult.Valid)
    }

    @Test
    fun `empty email returns invalid`() {
        assertTrue(rule.validate("") is ValidationResult.Invalid)
    }

    @Test
    fun `email without at symbol returns invalid`() {
        assertTrue(rule.validate("userexample.com") is ValidationResult.Invalid)
    }

    @Test
    fun `email without domain returns invalid`() {
        assertTrue(rule.validate("user@") is ValidationResult.Invalid)
    }

    @Test
    fun `email without tld returns invalid`() {
        assertTrue(rule.validate("user@example") is ValidationResult.Invalid)
    }

    @Test
    fun `whitespace trimmed email returns valid`() {
        assertTrue(rule.validate("  user@example.com  ") is ValidationResult.Valid)
    }
}
