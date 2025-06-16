package service

import "errors"

var (
	ErrUserNotFound = errors.New("user not found")
	ErrIncorrectPassword = errors.New("incorrect user password")
	ErrUsernameTaken = errors.New("username is taken")

	ErrTeamNotFound = errors.New("team not found")

	ErrNonexistentCode = errors.New("nonexistent code or key")
	
	ErrTranslationNotFound = errors.New("translation not found")
	ErrTranslationExists = errors.New("translation already exists")
)
