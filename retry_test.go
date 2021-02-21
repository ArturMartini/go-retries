package go_retries

import (
	"errors"
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestDefaultRetry(t *testing.T) {

	err := Do(func() error {
		return nil
	})

	assert.Nil(t, err)


	var retry = 0

	err = Do(func() error {
		if retry < 2 {
			retry++
			return errors.New("recoverable")
		} else {
			return nil
		}
	})

	assert.Equal(t, 2, retry)
	assert.Nil(t, err)


	unrecoverableError := errors.New("unrecoverable")
	Setting(nil, []error{unrecoverableError})
	err = Do(func() error {
		return unrecoverableError
	})

	assert.Error(t, ErrorUnrecoverable, err)
}
