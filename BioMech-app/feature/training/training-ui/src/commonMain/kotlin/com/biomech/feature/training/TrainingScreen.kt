package com.biomech.feature.training

import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.items
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.unit.dp
import com.biomech.domain.model.TrainingJob

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun TrainingScreen(
    jobs: List<TrainingJob>,
    sessionLabels: List<String>,
    selectedSessions: Set<String>,
    onToggleSession: (String) -> Unit,
    onStartTraining: () -> Unit,
) {
    Scaffold(
        topBar = {
            TopAppBar(title = { Text("Training") })
        }
    ) { padding ->
        Column(
            modifier = Modifier
                .fillMaxSize()
                .padding(padding)
                .padding(16.dp)
        ) {
            if (sessionLabels.isNotEmpty()) {
                Text(
                    text = "Select sessions for training:",
                    style = MaterialTheme.typography.titleMedium
                )
                Spacer(Modifier.height(8.dp))

                sessionLabels.forEach { label ->
                    Row(
                        modifier = Modifier.fillMaxWidth(),
                        verticalAlignment = Alignment.CenterVertically
                    ) {
                        Checkbox(
                            checked = label in selectedSessions,
                            onCheckedChange = { onToggleSession(label) }
                        )
                        Text(label)
                    }
                }

                Spacer(Modifier.height(16.dp))

                Button(
                    onClick = onStartTraining,
                    enabled = selectedSessions.isNotEmpty(),
                    modifier = Modifier.fillMaxWidth()
                ) {
                    Text("Start Training")
                }
            }

            Spacer(Modifier.height(24.dp))

            Text(
                text = "Training History",
                style = MaterialTheme.typography.titleMedium
            )
            Spacer(Modifier.height(8.dp))

            if (jobs.isEmpty()) {
                Box(
                    modifier = Modifier.fillMaxSize(),
                    contentAlignment = Alignment.Center
                ) {
                    Text("No training jobs yet")
                }
            } else {
                LazyColumn(
                    verticalArrangement = Arrangement.spacedBy(8.dp)
                ) {
                    items(jobs) { job ->
                        Card(modifier = Modifier.fillMaxWidth()) {
                            Column(modifier = Modifier.padding(16.dp)) {
                                Text(
                                    text = "Job ${job.id.take(8)}",
                                    style = MaterialTheme.typography.titleSmall
                                )
                                Text("Status: ${job.status.name}")
                                if (job.accuracy > 0) {
                                    Text("Accuracy: ${(job.accuracy * 100).toInt()}%")
                                }
                            }
                        }
                    }
                }
            }
        }
    }
}
