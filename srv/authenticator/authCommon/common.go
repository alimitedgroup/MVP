package common

import "errors"

var (
	ErrNoPublicKey     = errors.New("no public key is present")
	ErrNoKeyPair       = errors.New("missing key pair")
	ErrKeyPairNotValid = errors.New("not a valid key pair")
)
