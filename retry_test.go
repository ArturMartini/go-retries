package go_retries

import (
	"errors"
	"testing"
	"github.com/stretchr/testify/assert"
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