package service

import "errors"

var (
	ErrInvalidURL     = errors.New("invalid URL")
	ErrUnreachableURL = errors.New("URL is not reachable")
	ErrNotFound       = errors.New("short URL not found")
	ErrCollision      = errors.New("short code collision")
)
