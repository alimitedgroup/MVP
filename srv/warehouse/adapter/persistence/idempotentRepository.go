package persistence

type IIdempotentRepository interface {
	SaveEventID(event string, id string)
	IsAlreadyProcessed(event string, id string) bool
}
