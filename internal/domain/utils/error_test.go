package utils

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestHandleErrWithCh(t *testing.T) {
	cases := []struct {
		name          string
		ctx           context.Context
		err           error
		expectSendErr bool
		expectResult  bool
	}{
		{
			name:          "Nil error",
			ctx:           context.Background(),
			err:           nil,
			expectSendErr: false,
			expectResult:  false,
		},
		{
			name:          "Cancelled context",
			ctx:           cancelledContext(),
			err:           fmt.Errorf("test error"),
			expectSendErr: false,
			expectResult:  true,
		},
		{
			name:          "Non-cancelled context with error",
			ctx:           context.Background(),
			err:           fmt.Errorf("test error"),
			expectSendErr: true,
			expectResult:  true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			chErr := make(chan error, 1)
			result := HandleErrWithCh(tc.ctx, chErr, tc.err)

			if result != tc.expectResult {
				t.Errorf("Expected result %v, got %v", tc.expectResult, result)
			}

			if tc.expectSendErr {
				select {
				case err := <-chErr:
					if err == nil || err.Error() != tc.err.Error() {
						t.Errorf("Expected error %v, got %v", tc.err, err)
					}
				case <-time.After(100 * time.Millisecond):
					t.Errorf("Expected error was not received")
				}
			} else {
				select {
				case <-chErr:
					t.Errorf("No error should be sent")
				default:
					// No error expected, so this is the correct path
				}
			}
		})
	}
}

func cancelledContext() context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	return ctx
}
