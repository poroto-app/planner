package utils

import (
	"context"
)

func HandleErrWithCh(ctx context.Context, chErr chan<- error, err error) bool {
	if err == nil {
		return false
	}

	select {
	case <-ctx.Done():
	// context is done, no need to send error to channel
	default:
		chErr <- err
	}

	return true
}

func HandleWrappedErrWithCh(ctx context.Context, chErr chan<- error, err error, wrappedErr error) bool {
	if err == nil {
		return false
	}

	if wrappedErr != nil {
		err = wrappedErr
	}

	return HandleErrWithCh(ctx, chErr, err)
}
