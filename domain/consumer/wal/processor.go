package wal

import "context"

type Processor interface {
	Processor(ctx context.Context, topic string, message []byte)
}

const WalChangeTopic = "go_wal_postgres.wal_change"

type processor struct {
	h *Handler
}

func NewProcessor(h *Handler) Processor {
	return &processor{
		h: h,
	}
}

func (p *processor) Processor(ctx context.Context, topic string, message []byte) {
	// Do something with the message
	if topic != WalChangeTopic {
		return
	}
	if err := p.h.Handle(ctx, message); err != nil {
		// Handle error
	}
}
