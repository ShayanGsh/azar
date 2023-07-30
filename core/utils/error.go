package utils

import (
	"strings"

)


type Error interface{
	Error() string
	String() string
	Errorf(format string, a ...interface{String() string}) string
}

type CustomError struct {
	msg string
}

func (e CustomError) Error() string {
	var sb strings.Builder
	sb.WriteString("Error: ")
	sb.WriteString(e.msg)
	return sb.String()
}

func (e CustomError) String() string {
	return e.Error()
}

func (e CustomError) Errorf(format string, a ...interface{String() string}) string {
	var sb strings.Builder
	sb.WriteString("Error: ")
	sb.WriteString(e.msg)
	sb.WriteString(" ")
	sb.WriteString(format)
	sb.WriteString(" ")
	for _, arg := range a {
		sb.WriteString(arg.String())
	}
	return sb.String()
}
// TODO: rework error handling entirely