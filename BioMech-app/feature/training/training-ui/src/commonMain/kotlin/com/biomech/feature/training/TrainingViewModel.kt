package com.biomech.feature.training

import com.biomech.core.common.AppResult
import com.biomech.core.mvi.BaseAction
import com.biomech.core.mvi.BaseEvent
import com.biomech.core.mvi.BaseState
import com.biomech.core.mvi.BaseViewModel
import com.biomech.domain.model.EMGSession
import com.biomech.domain.model.TrainingJob
import com.biomech.domain.repository.EMGRepository
import com.biomech.domain.repository.TrainingRepository
import kotlinx.coroutines.channels.Channel
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.launch

data class TrainingState(
    val jobs: List<TrainingJob> = emptyList(),
    val sessions: List<EMGSession> = emptyList(),
    val selectedSessionIds: Set<String> = emptySet(),
    val isCreating: Boolean = false,
    val error: String? = null,
) : BaseState

sealed class TrainingAction : BaseAction {
    data class ToggleSession(val sessionId: String) : TrainingAction()
    data object StartTraining : TrainingAction()
}

sealed class TrainingEvent : BaseEvent {
    data object TrainingCreated : TrainingEvent()
}

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
                sessions = sessions.getOrNull() ?: emptyList(),
                jobs = jobs.getOrNull() ?: emptyList(),
            )
        }
    }

    override suspend fun handleAction(action: TrainingAction) {
        when (action) {
            is TrainingAction.ToggleSession -> {
                val current = _state.value.selectedSessionIds.toMutableSet()
                if (action.sessionId in current) {
                    current.remove(action.sessionId)
                } else {
                    current.add(action.sessionId)
                }
                _state.value = _state.value.copy(selectedSessionIds = current, error = null)
            }
            TrainingAction.StartTraining -> {
                val sessionIds = _state.value.selectedSessionIds.toList()
                if (sessionIds.isEmpty()) return
                _state.value = _state.value.copy(isCreating = true, error = null)
                when (val result = trainingRepository.createJob(sessionIds)) {
                    is AppResult.Success -> {
                        _state.value = _state.value.copy(isCreating = false, selectedSessionIds = emptySet())
                        loadData()
                        _event.send(TrainingEvent.TrainingCreated)
                    }
                    is AppResult.Error -> {
                        _state.value = _state.value.copy(isCreating = false, error = result.message)
                    }
                }
            }
        }
    }
}

private fun <T> AppResult<T>.getOrNull(): T? =
    when (this) {
        is AppResult.Success -> data
        is AppResult.Error -> null
    }
