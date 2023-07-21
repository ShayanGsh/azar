package utils

import "strings"

type CustomError string

func (e CustomError) Error() string {
	var sb strings.Builder
	sb.WriteString("Error: ")
	sb.WriteString(string(e))
	return sb.String()
}