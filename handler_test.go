package slogtest

import (
	"context"
	"log/slog"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandlerHandle(t *testing.T) {
	h := NewHandler(t, func(record slog.Record) {
		assert.Equal(t, "test", record.Message)
	})
	record := slog.NewRecord(time.Now(), slog.LevelInfo, "test", 0)
	err := h.Handle(context.Background(), record)
	require.NoError(t, err)
	records := h.Records()
	require.Len(t, records, 1)
	assert.Equal(t, record, records[0])
}

func TestHandlerHandleConc(t *testing.T) {
	h := NewHandler(t, func(record slog.Record) {})
	var wg sync.WaitGroup
	const iterations = 100
	wg.Add(iterations)
	for i := 0; i < iterations; i++ {
		i := i
		go func() {
			defer wg.Done()

			record := slog.NewRecord(time.Now(), slog.LevelInfo, strconv.Itoa(i), 0)
			err := h.Handle(context.Background(), record)
			require.NoError(t, err)
		}()
	}
	wg.Wait()
	records := h.Records()
	require.Len(t, records, 100)
}
