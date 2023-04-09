// Package guti contains packages
package guti

import (
	"fmt"
	"time"
)

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

// RetryWithExponentialBackoff retries a function with exponential backoff in case of errors
func RetryWithExponentialBackoff(fn func() error, maxRetries int, initialBackoffSeconds int) error {
	backoff := time.Duration(initialBackoffSeconds) * time.Second
	subsequentBackoff := 2

	for i := 0; ; i++ {
		fmt.Println("Retrying....")
		err := fn()
		if err == nil {
			return nil
		}

		if i == maxRetries {
			return fmt.Errorf("maximum number of retries exceeded: %w", err)
		}

		//log.Printf("Error occurred: %v. Retrying in %v", err, backoff)
		time.Sleep(backoff)

		backoff *= time.Duration(subsequentBackoff)
	}
}
