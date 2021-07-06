package commands_test

import (
	"errors"
	"testing"

	"code.vegaprotocol.io/go-wallet/commands"
	"github.com/stretchr/testify/assert"
)

func TestErrors(t *testing.T) {
	t.Run("Adding errors succeeds", testAddingErrorsSucceeds)
}

func testAddingErrorsSucceeds(t *testing.T) {
	errs := commands.NewErrors()
	prop := "user"
	err1 := errors.New("this is a first error")
	err2 := errors.New("this is a second error")

	errs.AddForProperty(prop, err1)
	errs.AddForProperty(prop, err2)

	assert.Equal(t, []error{err1, err2}, errs.Get(prop))
}
