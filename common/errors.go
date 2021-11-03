package common

import (
	"fmt"
	"strings"
)

func CombineErrors(errs ...error) error {
	if len(errs) == 0 {
		return nil
	}
	var result errors
	for _, x := range errs {
		result = append(result, flat(x)...)
	}
	return result
}

func flat(e error) []error {
	if m, ok := e.(interface{ Errors() []error }); ok {
		return m.Errors()
	}
	return []error{e}
}

type errors []error

func (e errors) Error() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("(%d errors): [", len(e)))
	for i, x := range e {
		sb.WriteString(fmt.Sprintf("%s", x))
		if i > 0 {
			sb.WriteString(", ")
		}
	}
	sb.WriteString("]")
	return sb.String()
}

func (e errors) Errors() []error {
	return e
}
