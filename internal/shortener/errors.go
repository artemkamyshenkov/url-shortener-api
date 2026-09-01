package shortener

import "errors"

var ErrNotFound = errors.New("short URL not found")
var ErrConflict = errors.New("short code already exists")
var ErrInvalidURL = errors.New("invalid url")
var ErrCreateAttemptsExhausted = errors.New("could not create short URL")
var ErrCacheMiss = errors.New("cache miss")
