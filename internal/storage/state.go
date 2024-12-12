package storage

import "sync"

type UserState struct {
	PendingImage     []byte
	PendingText      string
	WaitingForDetail bool
}

type UserStateManager struct {
	states map[int64]UserState
	mutex  sync.RWMutex
}

func NewUserStateManager() *UserStateManager {
	return &UserStateManager{
		states: make(map[int64]UserState),
	}
}

func (m *UserStateManager) GetOrCreateState(userID int64) UserState {
	m.mutex.RLock()
	state, exists := m.states[userID]
	m.mutex.RUnlock()

	if !exists {
		m.mutex.Lock()
		m.states[userID] = UserState{}
		m.mutex.Unlock()
	}

	return state
}

func (m *UserStateManager) UpdateState(userID int64, state UserState) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.states[userID] = state
}
