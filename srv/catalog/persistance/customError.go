package goodRepository

type CustomError struct {
	er string
}

func (c CustomError) Error() string {
	return c.er
}
