package models

import "fmt"

type GenericOwnershipCheckFields struct {
	UUID    string
	Creator string
}

type Status int

const (
	StatusNone Status = iota
	StatusLoading
	StatusError
)

func (tcgs Status) Name() string {
	switch tcgs {
	case StatusNone:
		return "none"
	case StatusLoading:
		return "loading"
	case StatusError:
		return "error"
	default:
		return fmt.Sprintf("%d", tcgs)
	}
}
