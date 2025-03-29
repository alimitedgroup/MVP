package dto

import "github.com/alimitedgroup/MVP/common/dto/response"

func InvalidJson() response.ResponseDTO[string] {
	return response.ResponseDTO[string]{
		Message: "The payload you provided is not valid JSON",
		Error:   "bad_json",
	}
}
