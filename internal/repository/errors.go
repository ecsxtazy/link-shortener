package repository

import "errors"

var (
	ErrNotFound          = errors.New("link not found")
	ErrOriginalURLExists = errors.New("original url already exists")
	ErrShortCodeExists   = errors.New("short code already exists")
)
