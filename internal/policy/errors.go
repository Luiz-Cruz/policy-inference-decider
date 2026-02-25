package policy

import "errors"

var (
	ErrNoStartNode      = errors.New("graph has no node named start")
	ErrInvalidPolicyDot = errors.New("invalid policy dot")
	ErrInvalidCondition = errors.New("invalid condition")
)
