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
	ErrorMaxRetriesReached = errors.New("retry max retries reached")
	ErrorUnrecoverable     = errors.New("retry unrecoverable error")
)

type Config string

type Configuration struct {
	Key   Config
	Value int
}

type retry struct {
	retry                 int
	continueRecovery      bool
	configs               map[Config]int
	listRecoverableErrors []error
}

func New() *retry {
	return &retry{
		retry: 0,
		continueRecovery: true,
		configs: map[Config]int{
			ConfigDelaySec:   defaultDelaySec,
			ConfigMaxRetries: defaultMaxRetries,
		},
		listRecoverableErrors: []error{},
	}
}

func (r *retry) SetConfigurations(configurations ...Configuration) *retry {
	for _, c := range configurations {
		r.configs[c.Key] = c.Value
	}
	return r
}

func (r *retry) SetRecoverableErrors(errors ...error) *retry {
	for _, err := range errors {
		r.listRecoverableErrors = append(r.listRecoverableErrors, err)
	}
	return r
}

func (r *retry) Do(f func() interface{}) interface{} {
	r.retry = 0
	r.continueRecovery = true
	defer r.panicRecovery(f)
	return r.execRetry(f, false)
}

func (r *retry) execRetry(f func() interface{}, originPanic bool) interface{} {
	for {
		if r.retry >= r.configs[ConfigMaxRetries] {
			r.continueRecovery = false
			if originPanic {
				panic(ErrorMaxRetriesReached)
			}
			return ErrorMaxRetriesReached
		}

		fReturn := f()
		if err, ok := fReturn.(error); ok {
			if err != nil {
				if r.isRecoverableErrors(err) {
					<-time.After(time.Second * time.Duration(r.configs[ConfigDelaySec]))
					r.retry++
				} else {
					r.continueRecovery = false
					return ErrorUnrecoverable
				}
			}
		} else {
			r.continueRecovery = false
			return fReturn
		}
	}
}

func (r *retry) isRecoverableErrors(err error) bool {
	var isRecoverable = false
	for _, recErr := range r.listRecoverableErrors {
		if errors.Is(err, recErr) {
			isRecoverable = true
		}
	}
	return isRecoverable
}

func (r *retry) panicRecovery(f func() interface{}) {
	if r.continueRecovery {
		<-time.After(time.Duration(r.configs[ConfigDelaySec]) * time.Second)

		if recover() != nil {
			defer r.panicRecovery(f)
			r.retry++
			r.execRetry(f, true)
		}
	} else {
		if recover() != nil {
			panic(recover())
		}
	}
}
