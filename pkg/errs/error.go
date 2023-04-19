package errs

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
)

func Append(left error, right error) error {
	switch {
	case left == nil:
		return right
	case right == nil:
		return left
	}

	return fmt.Errorf("%w; %w", left, right)
}

func IsCanceledOrDeadline(err error) bool {
	return errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded)
}
