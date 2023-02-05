package runtime

import "errors"

var (
	ErrInvalidActionDataType = errors.New("invalid action data type")
	ErrNoTransition          = errors.New("no transition")
	ErrUnknownAction         = errors.New("unknown acton")
	ErrUnknownState          = errors.New("unknown state")
)
