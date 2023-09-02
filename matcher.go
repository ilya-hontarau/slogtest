package slogtest

import (
	"log/slog"
	"regexp"
	"slices"

	"github.com/stretchr/testify/assert"
)

type Matcher struct {
	inplaceAssertF []func(t assert.TestingT, record slog.Record)
	afterAssertF   []func(t assert.TestingT, records []slog.Record)
	t              assert.TestingT
	handlerAssertF []func()
}

func NewMatcher(t assert.TestingT) *Matcher {
	return &Matcher{t: t}
}

func (m Matcher) WithNoLevel(level slog.Level) *Matcher {
	m.inplaceAssertF = append(m.inplaceAssertF, func(t assert.TestingT, record slog.Record) {
		assert.NotEqual(m.t, level, record.Level)
	})
	return &m
}

func (m Matcher) WithNoMsg(msg string) *Matcher {
	m.inplaceAssertF = append(m.inplaceAssertF, func(t assert.TestingT, record slog.Record) {
		assert.NotEqual(m.t, msg, record.Message)
	})
	return &m
}

func (m Matcher) WithNoMsgRegExp(msg *regexp.Regexp) *Matcher {
	m.inplaceAssertF = append(m.inplaceAssertF, func(t assert.TestingT, record slog.Record) {
		ok := msg.MatchString(record.Message)
		if ok {
			m.t.Errorf("Expect \"%s\" to NOT match \"%s\"", record.Message, msg.String())
		}
	})
	return &m
}

func (m Matcher) WithMsg(msg string) *Matcher {
	m.afterAssertF = append(m.afterAssertF, func(t assert.TestingT, records []slog.Record) {
		ok := slices.ContainsFunc(records, func(r slog.Record) bool {
			return r.Message == msg
		})
		if !ok {
			m.t.Errorf("Can't find msg \"%s\" in logs", msg)
		}
	})
	return &m
}

func (m Matcher) WithMsgRegExp(msg *regexp.Regexp) *Matcher {
	m.afterAssertF = append(m.afterAssertF, func(t assert.TestingT, records []slog.Record) {
		ok := slices.ContainsFunc(records, func(r slog.Record) bool {
			return msg.MatchString(r.Message)
		})
		if !ok {
			m.t.Errorf("Can't find regexp \"%s\" in logs", msg.String())
		}
	})
	return &m
}

func (m *Matcher) Handler() *Handler {
	h := NewHandler(m.t, func(record slog.Record) {
		for _, f := range m.inplaceAssertF {
			f(m.t, record)
		}
	})
	m.handlerAssertF = append(m.handlerAssertF, func() {
		records := h.Records()
		for _, f := range m.afterAssertF {
			f(m.t, records)
		}
	})
	return h
}

func (m *Matcher) Finish() {
	for _, f := range m.handlerAssertF {
		f()
	}
}
