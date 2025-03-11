package persistence

import (
	"testing"
)

func TestCatalogRepositoryImpl(t *testing.T) {
	catalog := NewCatalogRepositoryImpl()

	catalog.SetGood("1", "blue_hat", "very beautiful hat")
	good := catalog.GetGood("1")

	if good == nil {
		t.Error("good is nil")
	} else {
		if good.Name != "blue_hat" {
			t.Error("good.Name is wrong")
		}
		if good.Description != "very beautiful hat" {
			t.Error("good.Description is wrong")
		}
	}

	if catalog.GetGood("2") != nil {
		t.Error("good should be nil")
	}

}
