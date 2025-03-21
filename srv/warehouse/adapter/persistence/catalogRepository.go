package persistence

type ICatalogRepository interface {
	GetGood(goodId string) *Good
	SetGood(goodId string, name string, description string) bool
}

type Good struct {
	ID          string
	Name        string
	Description string
}
