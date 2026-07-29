package inputconv_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/internal/api"
	"hexletbasics/internal/inputconv"
)

func TestPtr(t *testing.T) {
	t.Run("value", func(t *testing.T) {
		got := inputconv.Ptr(api.NewNilString("hello"))

		require.NotNil(t, got)
		assert.Equal(t, "hello", *got)
	})

	t.Run("null", func(t *testing.T) {
		got := inputconv.Ptr(api.NilBool{Null: true})

		assert.Nil(t, got)
	})
}

func TestStringPtr(t *testing.T) {
	t.Run("enum value", func(t *testing.T) {
		got := inputconv.StringPtr(api.NewNilCourseLearnAs(api.CourseLearnAsFirstLanguage))

		require.NotNil(t, got)
		assert.Equal(t, string(api.CourseLearnAsFirstLanguage), *got)
	})

	t.Run("null", func(t *testing.T) {
		got := inputconv.StringPtr(api.NilReviewState{Null: true})

		assert.Nil(t, got)
	})
}

func TestIntPtr(t *testing.T) {
	t.Run("value", func(t *testing.T) {
		got := inputconv.IntPtr(api.NewNilInt32(42))

		require.NotNil(t, got)
		assert.Equal(t, 42, *got)
	})

	t.Run("null", func(t *testing.T) {
		got := inputconv.IntPtr(api.NilInt32{Null: true})

		assert.Nil(t, got)
	})
}
