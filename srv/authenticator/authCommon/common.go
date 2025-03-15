package common

import "errors"

var (
	ErrNoPrivateKey    = errors.New("no private key is present")
	ErrNoPublicKey     = errors.New("no public key is present")
	ErrNoKeyPair       = errors.New("missing key pair")
	ErrKeyPairNotValid = errors.New("not a valid key pair")
	ErrNoIssuer        = errors.New("issuer is missing, generate key first")
	ErrPublish         = errors.New("cannot publish key")
	ErrUserNotLegit    = errors.New("user has wrong login data")
	ErrNoToken         = errors.New("cannot generate new token")
)
