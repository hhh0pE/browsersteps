package browsersteps

import (
	"context"
	"github.com/pkg/errors"
	"time"
)

func RunWithTimeout(f func() error, timeout, pingDuration time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	var lastError error
	for {
		select {
		case <-time.After(pingDuration):
			lastError = f()
			if lastError == nil {
				return nil
			}
		case <-ctx.Done():
			return errors.Wrap(lastError, "excedeed timeout")
		}
	}
	return lastError
}
