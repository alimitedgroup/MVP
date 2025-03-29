package types

import "errors"

var (
	ErrRuleExists    = errors.New("rule already exists")
	ErrRuleNotExists = errors.New("no rule found with the given id")
)
