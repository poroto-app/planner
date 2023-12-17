package utils

import (
	"context"
	"testing"
	"time"
)

func TestSendOrAbort(t *testing.T) {
	type testEntity struct {
		Value string
	}

	cases := []struct {
		name       string
		cancelCtx  bool
		expectSent bool
	}{
		{
			name:       "Non-cancelled context",
			cancelCtx:  false,
			expectSent: true,
		},
		{
			name:       "Cancelled context",
			cancelCtx:  true,
			expectSent: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			if tc.cancelCtx {
				var cancel context.CancelFunc
				ctx, cancel = context.WithCancel(ctx)
				cancel()
			}

			ch := make(chan *testEntity, 1)
			entity := &testEntity{"test"}

			if sent := SendOrAbort(ctx, ch, entity); sent != tc.expectSent {
				t.Errorf("sendOrAbort returned %v, expected %v", sent, tc.expectSent)
			}

			if tc.expectSent {
				select {
				case e := <-ch:
					if e != entity {
						t.Errorf("Expected entity to be sent to the channel")
					}
				case <-time.After(time.Millisecond * 100):
					t.Errorf("Expected entity was not sent to the channel")
				}
			} else {
				select {
				case <-ch:
					t.Errorf("No entity should be sent to the channel")
				case <-time.After(time.Millisecond * 100):
					// This is the expected path, as the context is cancelled
				}
			}
		})
	}
}
