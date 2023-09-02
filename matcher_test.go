package slogtest

import (
	"fmt"
	"log/slog"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMatcher(t *testing.T) {
	t.Run("find message in logs", func(t *testing.T) {
		matcher := NewMatcher(t).WithMsg("test")
		defer matcher.Finish()

		logger := slog.New(matcher.Handler())
		logger.Info("test")
	})
	t.Run("find message reg exp in logs", func(t *testing.T) {
		matcher := NewMatcher(t).WithMsgRegExp(regexp.MustCompile("test."))
		defer matcher.Finish()

		logger := slog.New(matcher.Handler())
		logger.Info("test5")
	})
	t.Run("find no message reg exp in logs", func(t *testing.T) {
		matcher := NewMatcher(t).WithNoMsgRegExp(regexp.MustCompile("^test.$"))
		defer matcher.Finish()

		logger := slog.New(matcher.Handler())
		logger.Info("test52")
	})
	t.Run("find no message reg exp in logs fail", func(t *testing.T) {
		mockT := &CollectT{}
		matcher := NewMatcher(mockT).WithNoMsgRegExp(regexp.MustCompile("^test.$"))
		defer matcher.Finish()

		logger := slog.New(matcher.Handler())
		logger.Info("test5")
		require.Len(t, mockT.errors, 1)
		assert.Equal(t, mockT.errors[0].Error(), "Expect \"test5\" to NOT match \"^test.$\"")
	})
	t.Run("find no message in logs", func(t *testing.T) {
		matcher := NewMatcher(t).WithNoMsg("test")
		defer matcher.Finish()

		logger := slog.New(matcher.Handler())
		logger.Info("test2")
	})
	t.Run("test immutability", func(t *testing.T) {
		m := NewMatcher(t)
		m.WithMsg("test")
		assert.Nil(t, m.afterAssertF)
	})
	t.Run("test groups find message", func(t *testing.T) {
		matcher := NewMatcher(t).WithMsg("test")
		defer matcher.Finish()

		logger := slog.New(matcher.Handler())
		logger.WithGroup("group").Info("test")
	})
}

type CollectT struct {
	errors []error
}

// Errorf collects the error.
func (c *CollectT) Errorf(format string, args ...interface{}) {
	c.errors = append(c.errors, fmt.Errorf(format, args...))
}
