package slogtest

import (
	"context"
	"log/slog"
	"slices"
	"sync"

	"github.com/stretchr/testify/assert"
)

type sharedRecords struct {
	mx      sync.RWMutex
	records []slog.Record
}

func newSharedRecords() *sharedRecords {
	return &sharedRecords{}
}

func (r *sharedRecords) Add(record slog.Record) {
	r.mx.Lock()
	defer r.mx.Unlock()
	r.records = append(r.records, record)
}

func (r *sharedRecords) All() []slog.Record {
	r.mx.RLock()
	defer r.mx.RUnlock()
	return r.records
}

type Handler struct {
	t      assert.TestingT
	groups []string
	attrs  []slog.Attr

	recordCallbackF func(record slog.Record)

	records *sharedRecords
}

func NewHandler(t assert.TestingT, recordCallbackF func(record slog.Record)) *Handler {
	return &Handler{
		t:               t,
		recordCallbackF: recordCallbackF,
		records:         newSharedRecords(),
	}
}

func (h *Handler) Enabled(_ context.Context, _ slog.Level) bool {
	return true
}

func (h *Handler) Handle(_ context.Context, record slog.Record) error {
	h.recordCallbackF(record)
	h.records.Add(record)
	return nil
}

func (h *Handler) Records() []slog.Record {
	return h.records.All()
}

func (h *Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	handler := h.clone()
	handler.attrs = append(handler.attrs, attrs...)
	return handler
}

func (h *Handler) WithGroup(name string) slog.Handler {
	handler := h.clone()
	handler.groups = append(handler.groups, name)
	return handler
}

func (h *Handler) clone() *Handler {
	return &Handler{
		t:               h.t,
		groups:          slices.Clone(h.groups),
		attrs:           slices.Clone(h.attrs),
		recordCallbackF: h.recordCallbackF,
		records:         h.records,
	}
}
