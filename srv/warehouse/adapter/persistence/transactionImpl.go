package persistence

import "sync"

type TransactionImpl struct {
	m sync.Mutex
}

func NewTransactionImpl() *TransactionImpl {
	return &TransactionImpl{m: sync.Mutex{}}
}

func (t *TransactionImpl) Lock() {
	t.m.Lock()
}

func (t *TransactionImpl) Unlock() {
	t.m.Unlock()
}
