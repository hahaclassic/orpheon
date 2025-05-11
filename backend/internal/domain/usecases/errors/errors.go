package commonerr

import "errors"

var (
	ErrForbidden = errors.New("permission denied error")
)
