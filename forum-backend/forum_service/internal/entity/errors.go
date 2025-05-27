package entity

import "errors"

var (
	ErrInvalidAuthorID = errors.New("invalid author id")
	ErrInvalidPostID   = errors.New("invalid post id")
	ErrInvalidUserID   = errors.New("invalid user id")
	ErrEmptyTitle      = errors.New("empty title")
	ErrEmptyContent    = errors.New("empty content")
	ErrEmptyUsername   = errors.New("empty username")
)
