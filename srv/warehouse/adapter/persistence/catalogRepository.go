package persistence

type CatalogRepository interface {
	GetGood(goodId string) *Good
	SetGood(goodId string, name string, description string) bool
}

type Good struct {
	Id          string
	Name        string
	Description string
}
