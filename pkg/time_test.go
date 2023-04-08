package pkg

import (
	"errors"
	"testing"
	"time"
)

func TestRetry(t *testing.T) {
	retries := 3
	sleepDuration := time.Millisecond

	retryFuncCount := 0
	err := Retry(func() error {
		if retryFuncCount < 2 {
			retryFuncCount++
			return errors.New("error")
		}
		return nil
	}, retries, sleepDuration)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}
