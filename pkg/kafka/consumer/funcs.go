package consumer

import (
	"context"
	"fmt"
)

func (cg *Consumer) Consume(ctx context.Context, topics []string, handler MessageHandler) error {
	cgh := newConsumerGroupHandler(handler)

	err := cg.conn.Consume(ctx, topics, cgh)
	if err != nil {
		return fmt.Errorf("cg.conn.Consume: %w", err)
	}

	return nil
}
