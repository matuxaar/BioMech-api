package com.biomech.feature.dashboard

import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.unit.dp
import com.biomech.core.component.EMGChart

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun DashboardScreen(
    deviceConnected: Boolean,
    emgData: List<Float>,
    onStartRecording: () -> Unit,
    onStopRecording: () -> Unit,
) {
    Scaffold(
        topBar = {
            TopAppBar(
                title = { Text("Dashboard") },
                actions = {
                    if (deviceConnected) {
                        Text(
                            text = "● Connected",
                            color = MaterialTheme.colorScheme.primary,
                            style = MaterialTheme.typography.labelMedium
                        )
                    }
                }
            )
        }
    ) { padding ->
        LazyColumn(
            modifier = Modifier
                .fillMaxSize()
                .padding(padding)
                .padding(16.dp),
            verticalArrangement = Arrangement.spacedBy(16.dp)
        ) {
            item {
                Card(
                    modifier = Modifier.fillMaxWidth()
                ) {
                    Column(modifier = Modifier.padding(16.dp)) {
                        Text(
                            text = "EMG Signal",
                            style = MaterialTheme.typography.titleMedium
                        )
                        Spacer(Modifier.height(8.dp))
                        EMGChart(
                            data = emgData,
                            modifier = Modifier.height(200.dp)
                        )
                    }
                }
            }

            item {
                Row(
                    modifier = Modifier.fillMaxWidth(),
                    horizontalArrangement = Arrangement.spacedBy(16.dp)
                ) {
                    OutlinedButton(
                        onClick = onStartRecording,
                        modifier = Modifier.weight(1f)
                    ) {
                        Text("Start Recording")
                    }
                    Button(
                        onClick = onStopRecording,
                        modifier = Modifier.weight(1f)
                    ) {
                        Text("Stop")
                    }
                }
            }
        }
    }
}
