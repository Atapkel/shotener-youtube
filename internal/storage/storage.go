package storage

import "errors"

var (
	ErrURLNotfound = errors.New("URL not found")
	ErrURLExists   = errors.New("URL exists")
)
