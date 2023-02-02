// Copyright 2022 Michal Vyskocil. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package pipe_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	. "github.com/gomoni/gio/pipe"
)

func TestError(t *testing.T) {
	t.Parallel()

	t.Run("error", func(t *testing.T) {
		var err error
		var e Error
		t.Parallel()
		err = NewError(42, errors.New("pipe.Error"))
		e = FromError(err)
		require.EqualValues(t, 42, e.Code)
		require.EqualError(t, e.Err, "pipe.Error")
		require.EqualError(t, err, "Error{Code: 42, Err: pipe.Error}")
	})

	t.Run("errorf", func(t *testing.T) {
		var err error
		var e Error
		t.Parallel()
		err = NewErrorf(142, "pipe: %w", errors.New("Errorf"))
		e = FromError(err)
		require.EqualValues(t, 142, e.Code)
		require.EqualError(t, e.Err, "pipe: Errorf")
		require.EqualError(t, err, "Error{Code: 142, Err: pipe: Errorf}")
	})

	t.Run("as error", func(t *testing.T) {
		var err error
		var e Error
		t.Parallel()
		err = errors.New("random error")
		e = FromError(err)
		require.EqualValues(t, UnknownError, e.Code)
		require.EqualError(t, e.Err, "random error")
		require.EqualError(t, e, "Error{Code: 250, Err: random error}")
	})
}
