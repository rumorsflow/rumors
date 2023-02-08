package errs

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"go.uber.org/multierr"
	"runtime"
	"strings"
)

var _ error = (*Error)(nil)

type Op string

type Error struct {
	Err error
	Op  Op
	ID  uuid.UUID
}

func (e *Error) Errors() []error {
	return multierr.Errors(e.Err)
}

func (e *Error) Error() string {
	str := ""

	if e.Op != "" {
		str += string(e.Op) + " -> "
	}

	if e.ID != uuid.Nil {
		str += e.ID.String() + " -> "
	}

	if e.Err != nil {
		str += e.Err.Error()
	} else {
		str += "unknown error"
	}

	return str
}

func (e *Error) Unwrap() error {
	return e.Err
}

func (e *Error) Cause() error {
	return e.Unwrap()
}

func E(args ...any) error {
	e := &Error{}

	if len(args) == 0 {
		msg := "errors.E called with 0 args"
		_, file, line, ok := runtime.Caller(1)
		if ok {
			msg = fmt.Sprintf("%v - %v:%v", msg, file, line)
		}
		e.Err = errors.New(msg)
	}

	for _, arg := range args {
		if arg == nil {
			continue
		}

		switch arg := arg.(type) {
		case Op:
			e.Op = arg
		case string:
			if arg != "" {
				e.Err = multierr.Append(e.Err, errors.New(arg))
			}
		case *Error:
			eCopy := *arg
			e.Err = multierr.Append(e.Err, &eCopy)
		case error:
			e.Err = multierr.Append(e.Err, arg)
		case []error:
			e.Err = multierr.Combine(append(e.Errors(), arg...)...)
		case uuid.UUID:
			e.ID = arg
		default:
			e.Err = multierr.Append(e.Err, errors.Errorf("unknown type %T, value %v in error call", arg, arg))
		}
	}

	return e
}

func Errorf(op Op, format string, a ...any) error {
	return E(op, errors.WithStack(fmt.Errorf(format, a...)))
}

func New(text string) error {
	return errors.New(text)
}

func Is(err, target error) bool {
	return errors.Is(err, target)
}

func As(err error, target any) bool {
	return errors.As(err, target)
}

func Unwrap(err error) error {
	return errors.Unwrap(err)
}

func HasOp(err error, op Op) bool {
	return err != nil && string(op) != "" && strings.Contains(err.Error(), string(op))
}

func IsCanceledOrDeadline(err error) bool {
	if Is(err, context.Canceled) || Is(err, context.DeadlineExceeded) {
		return true
	}
	return false
}
