package pipe

import (
	"errors"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestErrorSliceNoPipeFail(t *testing.T) {
	testCases := []struct {
		name     string
		errs     []error
		expected error
	}{
		{
			name:     "nil",
			errs:     nil,
			expected: nil,
		},
		{
			name:     "last nil",
			errs:     []error{io.EOF, nil},
			expected: nil,
		},
		{
			name:     "last io.EOF",
			errs:     []error{nil, io.EOF},
			expected: NewError(1, io.EOF),
		},
		{
			name:     "io.ErrUnexpectedEOF, nil, io.EOF",
			errs:     []error{io.ErrUnexpectedEOF, nil, io.EOF},
			expected: NewError(1, errors.New("unexpected EOF\nEOF")),
		},
		{
			name:     "nil, Error",
			errs:     []error{nil, NewErrorf(42, "Error")},
			expected: NewError(42, NewErrorf(42, "Error")),
		},
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := errorSlice{errs: tt.errs}.noPipefail(1)
			if tt.expected == nil {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, tt.expected.Error())
			}
		})
	}
}

func TestErrorSlicePipeFail(t *testing.T) {
	testCases := []struct {
		name     string
		errs     []error
		expected error
	}{
		{
			name:     "nil",
			errs:     nil,
			expected: nil,
		},
		{
			name:     "last nil",
			errs:     []error{io.EOF, nil},
			expected: NewError(1, io.EOF),
		},
		{
			name:     "last io.EOF",
			errs:     []error{nil, io.EOF},
			expected: NewError(1, io.EOF),
		},
		{
			name:     "io.ErrUnexpectedEOF, nil, io.EOF",
			errs:     []error{io.ErrUnexpectedEOF, nil, io.EOF},
			expected: NewError(1, errors.New("unexpected EOF\nEOF")),
		},
		{
			name:     "Error, nil",
			errs:     []error{NewErrorf(42, "Error"), nil},
			expected: errors.New("Error{Code: 42, Err: Error{Code: 42, Err: Error}}"),
		},
		{
			name:     "Error, nil, Error2",
			errs:     []error{NewErrorf(42, "Error"), nil, NewErrorf(1, "Error2")},
			expected: errors.New("Error{Code: 42, Err: Error{Code: 42, Err: Error}\nError{Code: 1, Err: Error2}}"),
		},
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := errorSlice{errs: tt.errs}.pipefail(1)
			if tt.expected == nil {
				require.NoError(t, err)
				return
			}

			require.EqualError(t, err, tt.expected.Error())
			for _, e := range Errors(err) {
				if e == nil {
					continue
				}
				require.Contains(t, tt.errs, e)
			}
		})
	}
}
