package go_retries

import (
	"errors"
	"testing"
	"github.com/stretchr/testify/assert"
	"time"
)

func TestRetry(t *testing.T) {
	err := Do(func() interface{} {
		return nil
	})

	assert.Nil(t, err)
}

func TestRetryRecover(t *testing.T) {
	var retry = 0

	recError := errors.New("recoverable")
	SetRecoverableErrors(recError)
	err := Do(func() interface{} {
		if retry < 2 {
			retry++
			return recError
		} else {
			return nil
		}
	})

	assert.Equal(t, 2, retry)
	assert.Nil(t, err)
}

func TestUnrecover(t *testing.T) {
	unrecoverableError := errors.New("unrecoverable")
	err := Do(func() interface{} {
		return unrecoverableError
	})

	assert.Error(t, ErrorUnrecoverable, err)
}

func TestRetryPanicRecovery(t *testing.T) {
	var retry = 0

	recoverableErr := errors.New("recoverable")
	SetRecoverableErrors(recoverableErr)
	err := Do(func() interface{} {
		if retry < 2 {
			retry++
			panic(recoverableErr)
		} else {
			return nil
		}
	})

	assert.Nil(t, err)
	assert.Equal(t, 2, retry)
}

func TestRetryTime(t *testing.T) {
	//The default delay is 3 seconds
	//The default max retries is 3

	start := time.Now()
	recoverableErr := errors.New("recoverable")
	SetRecoverableErrors(recoverableErr)

	err := Do(func() interface{} {
			return recoverableErr
	})

	assert.Equal(t, ErrorMaxRetriesReached, err)
	assert.Equal(t, true, time.Since(start) > time.Second * 9)
	assert.Equal(t, true, time.Since(start) < time.Second * 10)
}

func TestRetryCustomTime(t *testing.T) {
	start := time.Now()
	SetConfigurations(
		Configuration{Key: ConfigMaxRetries, Value: 1},
		Configuration{Key: ConfigDelaySec, Value: 2})

	recoverableErr := errors.New("recoverable")
	SetRecoverableErrors(recoverableErr)
	err := Do(func() interface{} {
		return recoverableErr
	})

	assert.Equal(t, ErrorMaxRetriesReached, err)
	assert.Equal(t, true, time.Since(start) > time.Second * 2)
	assert.Equal(t, true, time.Since(start) < time.Second * 3)
}