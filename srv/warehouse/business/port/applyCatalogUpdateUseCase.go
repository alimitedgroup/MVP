package port

import "context"

type IApplyCatalogUpdateUseCase interface {
	ApplyCatalogUpdate(ctx context.Context, cmd CatalogUpdateCmd) error
}

type CatalogUpdateCmd struct {
	GoodID      string
	Name        string
	Description string
}
