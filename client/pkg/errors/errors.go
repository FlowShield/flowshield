package errors

import (
	"github.com/pkg/errors"
)

// alias
var (
	New          = errors.New
	Wrap         = errors.Wrap
	Wrapf        = errors.Wrapf
	WithStack    = errors.WithStack
	WithMessage  = errors.WithMessage
	WithMessagef = errors.WithMessagef
)

var NewWithStack = func(msg string) error {
	return errors.WithStack(errors.New(msg))
}
