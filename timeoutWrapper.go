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
	for i := 0; ; i++ {
		lastError = f()
		if lastError == nil {
			return nil
		}

		select {
		case <-time.After(pingDuration):
			continue
		case <-ctx.Done():
			return errors.Wrapf(lastError, "exceeded timeout, tried %d times", i+1)
		}
	}
	return lastError
}
