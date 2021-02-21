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
	listRecoverableErrors []error

	configs = map[Config]int{
		ConfigMaxRetries: defaultMaxRetries,
		ConfigDelaySec:   defaultDelaySec,
	}

	ErrorMaxRetriesReached = errors.New("retry max retries reached")
	ErrorUnrecoverable     = errors.New("retry unrecoverable error")
)

type Config string

type Configuration struct {
	Key   Config
	Value int
}

func SetConfigurations(configurations ...Configuration) {
	for _, c := range configurations {
		configs[c.Key] = c.Value
	}
}

func SetRecoverableErrors(errors ...error) {
	for _, err := range errors {
		listRecoverableErrors = append(listRecoverableErrors, err)
	}
}

func Do(f func() interface{}) interface{} {
	var retry = 0
	var continueRecovery = true
	defer panicRecovery(f, &retry, &continueRecovery)
	return execRetry(f, &retry, &continueRecovery)
}

func execRetry(f func() interface{}, retry *int, continueRecovery *bool) interface{} {
	for {
		if *retry >= configs[ConfigMaxRetries] {
			*continueRecovery = false
			return ErrorMaxRetriesReached
		}

		fReturn := f()
		if err, ok := fReturn.(error); ok {
			if err != nil {
				if isRecoverableErrors(err) {
					<-time.After(time.Second * time.Duration(configs[ConfigDelaySec]))
					(*retry)++
				} else {
					*continueRecovery = false
					return ErrorUnrecoverable
				}
			}
		} else {
			*continueRecovery = false
			return fReturn
		}
	}
}

func isRecoverableErrors(err error) bool {
	var isRecoverable = false
	for _, recErr := range listRecoverableErrors {
		if errors.Is(err, recErr) {
			isRecoverable = true
		}
	}
	return isRecoverable
}

func panicRecovery(f func() interface{}, retry *int, continueRecovery *bool) {
	if *continueRecovery {
		defer panicRecovery(f, retry, continueRecovery)
		r := recover()
		if r != nil {
			*retry++
			if *retry <= configs[ConfigMaxRetries] {
				execRetry(f, retry, continueRecovery)
			} else {
				*continueRecovery = false
			}
		}
	}
}
