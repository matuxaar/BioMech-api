package com.biomech.core.mvi

import kotlinx.coroutines.CoroutineScope
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.SupervisorJob
import kotlinx.coroutines.cancel
import kotlinx.coroutines.channels.Channel
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.flow.receiveAsFlow
import kotlinx.coroutines.launch

abstract class BaseViewModel<S : BaseState, A : BaseAction, E : BaseEvent> {

    protected val scope = CoroutineScope(SupervisorJob() + Dispatchers.Main)

    protected abstract val _state: MutableStateFlow<S>
    val state: StateFlow<S> by lazy { _state.asStateFlow() }

    protected abstract val _event: Channel<E>
    val event by lazy { _event.receiveAsFlow() }

    protected abstract suspend fun handleAction(action: A)

    fun dispatch(action: A) {
        scope.launch {
            handleAction(action)
        }
    }

    protected fun updateState(transform: S.() -> S) {
        _state.value = _state.value.transform()
    }

    protected fun sendEvent(event: E) {
        scope.launch {
            _event.send(event)
        }
    }

    open fun onCleared() {
        scope.cancel()
    }
}
