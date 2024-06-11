package errors

import "errors"

var (
	ErrorStorageGenShort = errors.New("error generate random short uri")
	ErrorStorageGetShort = errors.New("error get url by short from storage")
)
