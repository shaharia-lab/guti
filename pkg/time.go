// Package pkg gotil contains packages
package pkg

import "time"

// Retry allows you to retry a given function a certain number of times until it succeeds or the retries are exhausted
func Retry(f func() error, retries int, sleep time.Duration) error {
	var err error
	for i := 0; i < retries; i++ {
		err = f()
		if err == nil {
			return nil
		}
		time.Sleep(sleep)
	}
	return err
}
