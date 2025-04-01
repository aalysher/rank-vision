package crawler

import "errors"

var (
	// ErrInvalidURL возникает, когда URL невалиден
	ErrInvalidURL = errors.New("invalid URL")

	// ErrUnexpectedStatusCode возникает, когда сервер возвращает неожиданный статус код
	ErrUnexpectedStatusCode = errors.New("unexpected status code")
)
