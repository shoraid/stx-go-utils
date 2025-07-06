package asyncutil

import (
	"errors"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAsyncUtil_SafeGo(t *testing.T) {
	type testCase[T any] struct {
		name     string
		fn       func() (T, error)
		expected Result[T]
		checkErr func(error) bool
	}

	t.Run("int cases", func(t *testing.T) {
		tests := []testCase[int]{
			{
				name: "success",
				fn: func() (int, error) {
					return 42, nil
				},
				expected: Result[int]{Value: 42, Err: nil},
			},
			{
				name: "error",
				fn: func() (int, error) {
					return 0, errors.New("something went wrong")
				},
				expected: Result[int]{Value: 0, Err: errors.New("something went wrong")},
				checkErr: func(err error) bool {
					return assert.ErrorContains(t, err, "something went wrong")
				},
			},
			{
				name: "panic",
				fn: func() (int, error) {
					panic("boom!")
				},
				expected: Result[int]{Value: 0},
				checkErr: func(err error) bool {
					return assert.ErrorContains(t, err, "panic recovered: boom!")
				},
			},
		}

		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				ch := SafeGo(tc.fn)
				result := <-ch

				assert.Equal(t, tc.expected.Value, result.Value)

				if tc.checkErr != nil {
					tc.checkErr(result.Err)
				} else {
					assert.Equal(t, tc.expected.Err, result.Err)
				}
			})
		}
	})

	t.Run("string cases", func(t *testing.T) {
		tests := []testCase[string]{
			{
				name: "success",
				fn: func() (string, error) {
					return "hello", nil
				},
				expected: Result[string]{Value: "hello", Err: nil},
			},
			{
				name: "error",
				fn: func() (string, error) {
					return "", errors.New("string error")
				},
				expected: Result[string]{Value: "", Err: errors.New("string error")},
				checkErr: func(err error) bool {
					return assert.ErrorContains(t, err, "string error")
				},
			},
			{
				name: "panic",
				fn: func() (string, error) {
					panic("string panic")
				},
				expected: Result[string]{Value: ""},
				checkErr: func(err error) bool {
					return assert.ErrorContains(t, err, "panic recovered: string panic")
				},
			},
		}

		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				ch := SafeGo(tc.fn)
				result := <-ch

				assert.Equal(t, tc.expected.Value, result.Value)

				if tc.checkErr != nil {
					tc.checkErr(result.Err)
				} else {
					assert.Equal(t, tc.expected.Err, result.Err)
				}
			})
		}
	})
}

func TestAsyncUtil_SafeGo_WithOnPanic(t *testing.T) {
	type myType struct{ ID int }

	var panicCalled atomic.Bool
	var capturedErr error

	OnPanic = func(err error) {
		panicCalled.Store(true)
		capturedErr = err
	}

	tests := []struct {
		name          string
		fn            func() (myType, error)
		expected      myType
		expectedErr   bool
		expectedPanic bool
		errContains   string
	}{
		{
			name: "success",
			fn: func() (myType, error) {
				return myType{ID: 1}, nil
			},
			expected:      myType{ID: 1},
			expectedErr:   false,
			expectedPanic: false,
		},
		{
			name: "error without panic",
			fn: func() (myType, error) {
				return myType{}, errors.New("some error")
			},
			expected:      myType{},
			expectedErr:   true,
			expectedPanic: false,
			errContains:   "some error",
		},
		{
			name: "panic recovery",
			fn: func() (myType, error) {
				panic("boom!")
			},
			expected:      myType{},
			expectedErr:   true,
			expectedPanic: true,
			errContains:   "panic recovered: boom!",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			panicCalled.Store(false)
			capturedErr = nil

			ch := SafeGo(tt.fn)
			res := <-ch

			assert.Equal(t, tt.expected, res.Value)

			if tt.expectedErr {
				assert.Error(t, res.Err)
				assert.Contains(t, res.Err.Error(), tt.errContains)
			} else {
				assert.NoError(t, res.Err)
			}

			if tt.expectedPanic {
				assert.True(t, panicCalled.Load(), "OnPanic should have been called")
				assert.NotNil(t, capturedErr)
				assert.Contains(t, capturedErr.Error(), tt.errContains)
			} else {
				assert.False(t, panicCalled.Load(), "OnPanic should NOT be called")
				assert.Nil(t, capturedErr)
			}
		})
	}
}

func BenchmarkAsyncUtil_SafeGo(b *testing.B) {
	for b.Loop() {
		ch := SafeGo(func() (int, error) {
			return 42, nil
		})
		<-ch
	}
}
