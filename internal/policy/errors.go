package policy

import "errors"

var ErrNoStartNode = errors.New("graph has no node named start")
var ErrInvalidPolicyDot = errors.New("invalid policy dot")
