package persistence

import "testing"

func TestLockUnLock(t *testing.T) {
	tx := NewTransactionImpl()
	tx.Lock()
	defer tx.Unlock()
}
