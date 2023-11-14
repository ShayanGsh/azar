package utils

import (
	"strings"

)

func Error(msg ErrorMsg, data ...string) error {
	sb := strings.Builder{}
	sb.WriteString("Error: ")
	sb.WriteString(string(msg))
	if len(data) > 0 {
		sb.WriteString(" ")
		sb.WriteString(strings.Join(data, " - "))
	}
	return CustomError(sb.String())
}

type ErrorMsg string

type CustomError string

func (e CustomError) Error() string {
	return string(e)
}