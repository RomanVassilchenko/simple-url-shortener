package repository

import "errors"

var (
	ErrObjectNotFound = errors.New("object not found")
	ErrURLNotFound    = errors.New("url not found")
	ErrURLExists      = errors.New("url exists")
)
