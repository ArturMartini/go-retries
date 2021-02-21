package go_retries

import (
	"errors"
	"time"
)

const (
	defaultMaxRetries int = 3
	defaultDelaySec   int = 3

	ConfigMaxRetries Config = "max.retries"
	ConfigDelaySec   Config = "delay.sec"
)

var (
	listUnrecoverableErrors []error

	configs map[Config]int = map[Config]int{
		ConfigMaxRetries: defaultMaxRetries,
		ConfigDelaySec: defaultDelaySec,
	}

	ErrorMaxRetriesReached = errors.New("max retries reached")
	ErrorUnrecoverable     = errors.New("unrecoverable error")
)

type Config string

func Setting(configurations map[Config]int, unrecoverableErrors []error) {
	configs = configurations
	listUnrecoverableErrors = unrecoverableErrors
}

func Do(f func() error) error {
	var retry = 0

	for {
		if retry >= configs[ConfigMaxRetries] {
			return ErrorMaxRetriesReached
		}

		err := f()
		if err != nil {
			if isUnrecoverableErrors(err) {
				return ErrorUnrecoverable
			} else {
				<-time.After(time.Second * time.Duration(configs[ConfigDelaySec]))
				retry++
			}
		} else {
			return nil
		}
	}
}

func isUnrecoverableErrors(err error) bool {
	var isUnrecoverable = false
	for _, recErr := range listUnrecoverableErrors {
		if errors.Is(err, recErr) {
			isUnrecoverable = true
		}
	}
	return isUnrecoverable
}
