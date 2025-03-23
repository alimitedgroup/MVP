package serviceobject

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserData(t *testing.T) {
	obj := NewUserData("test")
	assert.Equal(t, "test", obj.GetUsername())
}
