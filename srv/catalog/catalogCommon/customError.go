package catalogCommon

type CustomError struct {
	err string
}

func NewCustomError(text string) *CustomError {
	return &CustomError{err: text}
}

func (c CustomError) Error() string {
	return c.err
}
