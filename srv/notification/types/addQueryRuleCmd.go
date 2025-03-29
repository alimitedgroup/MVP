package types

import "github.com/google/uuid"

type QueryRule struct {
	GoodId    string
	Operator  string
	Threshold int
}

type QueryRuleWithId struct {
	QueryRule
	RuleId uuid.UUID
}

type EditRule struct {
	GoodId    *string
	Operator  *string
	Threshold *int
}
