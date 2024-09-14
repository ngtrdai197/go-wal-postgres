package wal

import "context"

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) Handle(ctx context.Context, message []byte) error {
	// Do something with the message
	return nil
}
