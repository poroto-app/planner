package utils

import "context"

func SendOrAbort[T any](ctx context.Context, ch chan<- *T, v *T) bool {
	select {
	case <-ctx.Done():
		return false
	default:
		ch <- v
		return true
	}
}
