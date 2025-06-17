package genericutil_test

import (
	"testing"
	"time"

	"github.com/shoraid/stx-go-utils/genericutil"

	"github.com/stretchr/testify/assert"
)

func TestGenericUtil_Ptr(t *testing.T) {
	t.Run("String", func(t *testing.T) {
		s := "hello"
		ptr := genericutil.Ptr(s)
		assert.NotNil(t, ptr)
		assert.Equal(t, s, *ptr)
	})

	t.Run("Int", func(t *testing.T) {
		i := 42
		ptr := genericutil.Ptr(i)
		assert.NotNil(t, ptr)
		assert.Equal(t, i, *ptr)
	})

	t.Run("Float64", func(t *testing.T) {
		f := 3.14
		ptr := genericutil.Ptr(f)
		assert.NotNil(t, ptr)
		assert.Equal(t, f, *ptr)
	})

	t.Run("Bool", func(t *testing.T) {
		b := true
		ptr := genericutil.Ptr(b)
		assert.NotNil(t, ptr)
		assert.Equal(t, b, *ptr)
	})

	t.Run("Time", func(t *testing.T) {
		now := time.Now()
		ptr := genericutil.Ptr(now)
		assert.NotNil(t, ptr)
		assert.Equal(t, now, *ptr)
	})

	t.Run("Struct", func(t *testing.T) {
		type Example struct {
			Name string
		}
		ex := Example{Name: "Go"}
		ptr := genericutil.Ptr(ex)
		assert.NotNil(t, ptr)
		assert.Equal(t, ex, *ptr)
	})
}

func BenchmarkPtr(b *testing.B) {
	b.Run("String", func(b *testing.B) {
		s := "benchmark"
		for i := 0; i < b.N; i++ {
			_ = genericutil.Ptr(s)
		}
	})

	b.Run("Int", func(b *testing.B) {
		data := 123
		for i := 0; i < b.N; i++ {
			_ = genericutil.Ptr(data)
		}
	})

	b.Run("Struct", func(b *testing.B) {
		type Example struct {
			A int
			B string
		}
		ex := Example{A: 1, B: "data"}
		for i := 0; i < b.N; i++ {
			_ = genericutil.Ptr(ex)
		}
	})
}
