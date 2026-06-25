package com.biomech.core.component

import androidx.compose.foundation.Canvas
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.height
import androidx.compose.runtime.Composable
import androidx.compose.ui.Modifier
import androidx.compose.ui.geometry.Offset
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.graphics.Path
import androidx.compose.ui.graphics.drawscope.Stroke
import androidx.compose.ui.unit.dp

@Composable
fun EMGChart(
    data: List<Float>,
    modifier: Modifier = Modifier,
    lineColor: Color = Color(0xFF006B5E),
    gridColor: Color = Color(0xFFE0E0E0),
) {
    Canvas(
        modifier = modifier
            .fillMaxWidth()
            .height(200.dp)
    ) {
        val spacing = size.width / (data.size - 1).coerceAtLeast(1)
        val midY = size.height / 2
        val amplitude = size.height / 2 * 0.9f

        for (i in 0 until 4) {
            val y = size.height / 4 * (i + 1)
            drawLine(gridColor, Offset(0f, y), Offset(size.width, y), strokeWidth = 0.5f)
        }

        if (data.size < 2) return@Canvas

        val path = Path().apply {
            moveTo(0f, midY - data[0] * amplitude)
            for (i in 1 until data.size) {
                lineTo(i * spacing, midY - data[i] * amplitude)
            }
        }

        drawPath(path, lineColor, style = Stroke(width = 2f))
    }
}
