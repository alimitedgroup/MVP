package catalogCommon

import "errors"

var (
	ErrEmptyDescription = errors.New("description is empty")
	ErrEmptyName        = errors.New("name is empty")
	ErrGoodIdNotValid   = errors.New("not a valid goodID")
	ErrRequestNotValid  = errors.New("not a valid request")
	ErrGenericFailure   = errors.New("an error occured")
)
