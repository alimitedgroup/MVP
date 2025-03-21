package port

type IApplyCatalogUpdateUseCase interface {
	ApplyCatalogUpdate(cmd CatalogUpdateCmd)
}

type CatalogUpdateCmd struct {
	GoodID      string
	Name        string
	Description string
}
