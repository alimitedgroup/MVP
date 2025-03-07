package persistence

import (
	"testing"
)

func TestStockRepositoryIml(t *testing.T) {
	stock := NewStockRepositoryIml()
	stock.SetStock("1", 0)
	stock.AddStock("1", 10)
	good := stock.GetStock("1")

	if good != 10 {
		t.Errorf("Expected 10, got %d", good)
	}

	if stock.GetStock("2") != 0 {
		t.Error("good should be 0")
	}
}
