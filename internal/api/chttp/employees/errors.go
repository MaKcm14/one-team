package employees

import "errors"

var (
	// Error for the responses.
	ErrUnknownCitizenship = errors.New("unknown citizenship was got")
	ErrUnknownTitle       = errors.New("unknown title was got")
)
