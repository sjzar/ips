package errors

import (
	"errors"
)

var (
	ErrIPVersionNotSupported = errors.New("IP version not supported")
	ErrInvalidFieldsLength   = errors.New("invalid fields length")
	ErrInvalidCIDR           = errors.New("invalid CIDR")
	ErrCIDRConflict          = errors.New("CIDR conflict")
	ErrMetaNotFound          = errors.New("meta not found")
	ErrDBFormatNotSupported  = errors.New("format not supported")
	ErrEmptyFile             = errors.New("empty file")
	ErrDatabaseIsInvalid     = errors.New("database is invalid")
)
