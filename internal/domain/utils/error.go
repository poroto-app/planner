package utils

import (
	"context"
)

// HandleErrWithCh chErr が close したあとにエラーを送信しないようにするための関数
// 適切に動作させるために、この関数を呼び出す場合は以下のようにcontextを扱う必要がある
// ```
// ctx, cancel := context.WithCancel(ctx)
// defer cancel()
// ```
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
