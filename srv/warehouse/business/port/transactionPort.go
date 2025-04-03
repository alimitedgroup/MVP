package port

type ITransactionPort interface {
	Lock()
	Unlock()
}
