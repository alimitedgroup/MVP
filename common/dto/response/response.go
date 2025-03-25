package response

type ResponseDTO[T any] struct {
	Message T      `json:"message"`
	Error   string `json:"error"`
}

type ErrorResponseDTO ResponseDTO[any]
