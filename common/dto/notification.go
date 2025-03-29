package dto

import (
	"github.com/alimitedgroup/MVP/common/dto/response"
	"github.com/google/uuid"
)

func InvalidJson() response.ResponseDTO[string] {
	return response.ResponseDTO[string]{
		Message: "The payload you provided is not valid JSON",
		Error:   "bad_json",
	}
}

func RuleNotFound() response.ResponseDTO[string] {
	return response.ResponseDTO[string]{
		Message: "No rule with the given ID found",
		Error:   "not_found",
	}
}

type Rule struct {
	GoodId    string `json:"good_id"`
	Operator  string `json:"operator"`
	Threshold int    `json:"threshold"`
}

type RuleWithId struct {
	RuleId    uuid.UUID `json:"id"`
	GoodId    string    `json:"good_id"`
	Operator  string    `json:"operator"`
	Threshold int       `json:"threshold"`
}

type RuleEdit struct {
	RuleId    uuid.UUID `json:"id"`
	GoodId    *string   `json:"good_id"`
	Operator  *string   `json:"operator"`
	Threshold *int      `json:"threshold"`
}
