package serviceresponse

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetPemPrivateKeyResponse(t *testing.T) {
	a := []byte{7}
	obj := NewGetPemPrivateKeyResponse(&a, "test", nil)
	assert.Equal(t, a, *obj.GetPemPrivateKey())
	assert.Equal(t, "test", obj.GetIssuer())
	require.Nil(t, obj.GetError())
}

func TestCheckKeyPairExistenceResponse(t *testing.T) {
	obj := NewCheckKeyPairExistenceResponse(nil)
	require.Nil(t, obj.GetError())
}

func TestGetPemPublicKeyResponse(t *testing.T) {
	a := []byte{7}
	obj := NewGetPemPublicKeyResponse(&a, "test", nil)
	assert.Equal(t, a, *obj.GetPemPublicKey())
	assert.Equal(t, "test", obj.GetIssuer())
	require.Nil(t, obj.GetError())
}

func TestGetTokenResponse(t *testing.T) {
	obj := NewGetTokenResponse("test", nil)
	assert.Equal(t, "test", obj.GetToken())
	require.Nil(t, obj.GetError())
}

func TestPublishPublicKeyResponse(t *testing.T) {
	obj := NewPublishPublicKeyResponse(nil)
	require.Nil(t, obj.GetError())
}

func TestStorePemKeyPairResponse(t *testing.T) {
	obj := NewStorePemKeyPairResponse(nil)
	require.Nil(t, obj.GetError())
}
