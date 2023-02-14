package app_test

import (
	"errors"
	"testing"

	"github.com/brunoluiz/go-lab/core/app"
	"github.com/stretchr/testify/require"
)

func TestStructAssert(t *testing.T) {
	err := &app.ErrValidation{
		Errors: []string{"foo"},
	}

	terr := &app.ErrValidation{}
	require.Equal(t, true, errors.As(err, &terr))
	require.Equal(t, terr, err)

	// errors.Is does not work with structs, as its values might change
	require.Equal(t, false, errors.Is(err, &app.ErrValidation{}))
}

func TestInterfaceAssert(t *testing.T) {
	err := &app.ErrValidation{
		Errors: []string{"foo"},
	}

	// test against provided structure
	var terr app.Err
	require.Equal(t, true, errors.As(err, &terr))
	require.Equal(t, terr, err)
	require.Equal(t, app.ErrCodeValidation, terr.Code())

	// test against possible different structure, should fail
	var xerr interface {
		Code() int
	}
	require.Equal(t, false, errors.As(err, &xerr))
}
