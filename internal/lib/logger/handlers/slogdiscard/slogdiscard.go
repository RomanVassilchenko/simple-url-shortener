package slogdiscard

import (
	"context"
	"golang.org/x/exp/slog"
)

func NewDiscardLogger() *slog.Logger {
	return slog.New(NewDiscardHandler())
}

type DiscardHandler struct{}

func NewDiscardHandler() *DiscardHandler {
	return &DiscardHandler{}
}

func (h *DiscardHandler) Handle(_ context.Context, _ slog.Record) error {
	// Just ignore writing in journal
	return nil
}

func (h *DiscardHandler) WithAttrs(_ []slog.Attr) slog.Handler {
	// return the same handler, because there is not attrs to save
	return h
}

func (h *DiscardHandler) Enabled(_ context.Context, _ slog.Level) bool {
	// Always return false, because the journal writing is ignoring
	return false
}

func (h *DiscardHandler) WithGroup(_ string) slog.Handler {
	// return the same handler, because there is no group to save
	return h
}
