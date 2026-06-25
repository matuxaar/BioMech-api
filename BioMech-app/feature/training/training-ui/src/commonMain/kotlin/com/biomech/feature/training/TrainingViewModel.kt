package com.biomech.feature.training

import com.biomech.core.common.AppResult
import com.biomech.core.mvi.BaseAction
import com.biomech.core.mvi.BaseEvent
import com.biomech.core.mvi.BaseState
import com.biomech.core.mvi.BaseViewModel
import com.biomech.domain.model.TrainingJob
import com.biomech.domain.repository.EMGRepository
import com.biomech.domain.repository.TrainingRepository
import kotlinx.coroutines.channels.Channel
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.launch

data class TrainingState(
    val jobs: List<TrainingJob> = emptyList(),
    val sessionLabels: List<String> = emptyList(),
    val selectedSessions: Set<String> = emptySet(),
) : BaseState

sealed class TrainingAction : BaseAction {
    data class ToggleSession(val label: String) : TrainingAction()
    data object StartTraining : TrainingAction()
}

sealed class TrainingEvent : BaseEvent

class TrainingViewModel(
    private val trainingRepository: TrainingRepository,
    private val emgRepository: EMGRepository,
) : BaseViewModel<TrainingState, TrainingAction, TrainingEvent>() {

    override val _state = MutableStateFlow(TrainingState())
    override val _event = Channel<TrainingEvent>(Channel.BUFFERED)

    init {
        loadData()
    }

    private fun loadData() {
        scope.launch {
            val sessions = emgRepository.getSessions()
            val jobs = trainingRepository.getJobs()

            _state.value = _state.value.copy(
                sessionLabels = sessions.getOrNull()?.map { it.label } ?: emptyList(),
                jobs = jobs.getOrNull() ?: emptyList(),
            )
        }
    }

    override suspend fun handleAction(action: TrainingAction) {
        when (action) {
            is TrainingAction.ToggleSession -> {
                val current = _state.value.selectedSessions.toMutableSet()
                if (action.label in current) current.remove(action.label) else current.add(action.label)
                _state.value = _state.value.copy(selectedSessions = current)
            }
            TrainingAction.StartTraining -> {
                // TODO: Create training job with selected session IDs
            }
        }
    }
}

private fun <T> AppResult<T>.getOrNull(): T? =
    when (this) {
        is AppResult.Success -> data
        is AppResult.Error -> null
    }
