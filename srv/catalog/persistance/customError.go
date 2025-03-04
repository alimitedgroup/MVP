package goodRepository

type CustomError struct {
	err string
}

func (c CustomError) Error() string {
	return c.err
}
