package servicecmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAllObjects(t *testing.T) {
	obj := NewGetTokenCmd("test")
	obj2 := NewCheckPemKeyPairExistence()
	obj3 := NewGetPemPrivateKeyCmd()
	obj4 := NewGetPemPublicKeyCmd()
	a := []byte{7}
	obj5 := NewPublishPublicKeyCmd(&a, "test")
	obj6 := NewStorePemKeyPairCmd(&a, &a, "issuer")
	assert.Equal(t, obj.GetUsername(), "test")
	require.NotNil(t, obj2)
	require.NotNil(t, obj3)
	require.NotNil(t, obj4)
	assert.Equal(t, *obj5.GetKey(), []byte{7})
	assert.Equal(t, obj5.GetIssuer(), "test")
	assert.Equal(t, *obj6.GetPemPrivateKey(), []byte{7})
	assert.Equal(t, *obj6.GetPemPublicKey(), []byte{7})
	assert.Equal(t, obj6.GetIssuer(), "issuer")
}
