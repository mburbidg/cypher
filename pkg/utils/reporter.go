package utils

type Reporter interface {
	Error(line int, msg string)
}
