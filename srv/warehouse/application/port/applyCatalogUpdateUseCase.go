package port

import "context"

type ApplyCatalogUpdateUseCase interface {
	ApplyCatalogUpdate(ctx context.Context, cmd CatalogUpdateCmd) error
}

type CatalogUpdateCmd struct {
	GoodId      string
	Name        string
	Description string
}
