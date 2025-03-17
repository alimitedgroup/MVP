package port

type IIdempotentPort interface {
	SaveEventID(IdempotentCmd)
	IsAlreadyProcessed(IdempotentCmd) bool
}

type IdempotentCmd struct {
	Event string
	Id    string
}
