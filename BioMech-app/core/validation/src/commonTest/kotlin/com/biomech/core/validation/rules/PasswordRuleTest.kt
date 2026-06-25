package com.biomech.core.validation.rules

import com.biomech.core.validation.ValidationResult
import kotlin.test.Test
import kotlin.test.assertTrue

class PasswordRuleTest {

    private val rule = PasswordRule()

    @Test
    fun `valid password returns valid`() {
        assertTrue(rule.validate("Password1") is ValidationResult.Valid)
    }

    @Test
    fun `password with special chars returns valid`() {
        assertTrue(rule.validate("Str0ng!Pass") is ValidationResult.Valid)
    }

    @Test
    fun `short password returns invalid`() {
        assertTrue(rule.validate("Ab1") is ValidationResult.Invalid)
    }

    @Test
    fun `password without uppercase returns invalid`() {
        assertTrue(rule.validate("password1") is ValidationResult.Invalid)
    }

    @Test
    fun `password without lowercase returns invalid`() {
        assertTrue(rule.validate("PASSWORD1") is ValidationResult.Invalid)
    }

    @Test
    fun `password without digit returns invalid`() {
        assertTrue(rule.validate("Password") is ValidationResult.Invalid)
    }

    @Test
    fun `empty password returns invalid`() {
        assertTrue(rule.validate("") is ValidationResult.Invalid)
    }
}
