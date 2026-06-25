package com.biomech.core.datetime

import kotlinx.datetime.Clock
import kotlinx.datetime.Instant
import kotlinx.datetime.LocalDate
import kotlinx.datetime.LocalDateTime
import kotlinx.datetime.TimeZone
import kotlinx.datetime.toLocalDateTime

object DateTimeUtil {

    fun now(): Instant = Clock.System.now()

    fun nowLocal(): LocalDateTime = now().toLocalDateTime(TimeZone.currentSystemDefault())

    fun today(): LocalDate = nowLocal().date

    fun formatDate(date: LocalDate, pattern: String = "dd.MM.yyyy"): String {
        val parts = date.toString().split("-")
        return when (pattern) {
            "dd.MM.yyyy" -> "${parts[2]}.${parts[1]}.${parts[0]}"
            "yyyy-MM-dd" -> date.toString()
            else -> date.toString()
        }
    }

    fun formatTimestamp(instant: Instant): String {
        val local = instant.toLocalDateTime(TimeZone.currentSystemDefault())
        return "${local.hour.toString().padStart(2, '0')}:" +
               "${local.minute.toString().padStart(2, '0')} " +
               "${local.dayOfMonth}.${local.monthNumber}.${local.year}"
    }

    fun timeAgo(instant: Instant): String {
        val now = now()
        val diff = now - instant
        val seconds = diff.inWholeSeconds
        return when {
            seconds < 60 -> "только что"
            seconds < 3600 -> "${seconds / 60} мин. назад"
            seconds < 86400 -> "${seconds / 3600} ч. назад"
            seconds < 2592000 -> "${seconds / 86400} дн. назад"
            else -> formatTimestamp(instant)
        }
    }
}
