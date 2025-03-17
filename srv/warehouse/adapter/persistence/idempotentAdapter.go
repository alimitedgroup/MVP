package persistence

import "github.com/alimitedgroup/MVP/srv/warehouse/business/port"

type IDempotentAdapter struct {
	repo IIdempotentRepository
}

func NewIDempotentAdapter(repo IIdempotentRepository) *IDempotentAdapter {
	return &IDempotentAdapter{repo}
}

func (i *IDempotentAdapter) SaveEventID(cmd port.IdempotentCmd) {
	i.repo.SaveEventID(cmd.Event, cmd.Id)
}

func (i *IDempotentAdapter) IsAlreadyProcessed(cmd port.IdempotentCmd) bool {
	return i.repo.IsAlreadyProcessed(cmd.Event, cmd.Id)
}
