package com.biomech.feature.training

package com.biomech.feature.training

import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.items
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.unit.dp
import com.biomech.domain.model.EMGSession
import com.biomech.domain.model.TrainingJob

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun TrainingScreen(
    jobs: List<TrainingJob>,
    sessions: List<EMGSession>,
    selectedSessionIds: Set<String>,
    isCreating: Boolean = false,
    error: String? = null,
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
            if (sessions.isNotEmpty()) {
                Text(
                    text = "Select sessions for training:",
                    style = MaterialTheme.typography.titleMedium
                )
                Spacer(Modifier.height(8.dp))

                sessions.forEach { session ->
                    Row(
                        modifier = Modifier.fillMaxWidth(),
                        verticalAlignment = Alignment.CenterVertically
                    ) {
                        Checkbox(
                            checked = session.id in selectedSessionIds,
                            onCheckedChange = { onToggleSession(session.id) }
                        )
                        Spacer(Modifier.width(8.dp))
                        Text(session.label)
                    }
                }

                Spacer(Modifier.height(16.dp))

                if (error != null) {
                    Text(
                        text = error,
                        color = MaterialTheme.colorScheme.error,
                        style = MaterialTheme.typography.bodySmall,
                    )
                    Spacer(Modifier.height(8.dp))
                }

                Button(
                    onClick = onStartTraining,
                    enabled = selectedSessionIds.isNotEmpty() && !isCreating,
                    modifier = Modifier.fillMaxWidth()
                ) {
                    if (isCreating) {
                        CircularProgressIndicator(modifier = Modifier.size(20.dp))
                    } else {
                        Text("Start Training")
                    }
                }
            }

            if (sessions.isEmpty()) {
                Box(
                    modifier = Modifier.fillMaxWidth().height(120.dp),
                    contentAlignment = Alignment.Center
                ) {
                    Text(
                        "No sessions available",
                        color = MaterialTheme.colorScheme.onSurfaceVariant,
                    )
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
