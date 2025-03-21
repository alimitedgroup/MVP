package persistence

import "sync"

type IdempotentRepositoryImpl struct {
	mutex sync.Mutex
	m     map[string]map[string]struct{}
}

func NewIdempotentRepositoryImpl() *IdempotentRepositoryImpl {
	return &IdempotentRepositoryImpl{
		m: make(map[string]map[string]struct{}),
	}
}

func (i *IdempotentRepositoryImpl) SaveEventID(event string, id string) {
	i.mutex.Lock()
	defer i.mutex.Unlock()

	if _, ok := i.m[event]; !ok {
		i.m[event] = make(map[string]struct{})
	}

	i.m[event][id] = struct{}{}
}

func (i *IdempotentRepositoryImpl) IsAlreadyProcessed(event string, id string) bool {
	i.mutex.Lock()
	defer i.mutex.Unlock()

	if m, ok := i.m[event]; ok {
		if _, ok := m[id]; ok {
			return true
		}
	}

	return false
}
