package port

type TransactionPort interface {
	Lock()
	Unlock()
}
