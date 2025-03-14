package catalogCommon

import "errors"

var (
	ErrGoodIdNotValid  = errors.New("not a valid goodID")
	ErrRequestNotValid = errors.New("not a valid request")
	ErrGenericFailure  = errors.New("an error occured")
)
