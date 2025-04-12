package stats

import (
	"context"

	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
)

type ListeningEventConsumer interface {
	ConsumeListeningEvent(ctx context.Context, event *entity.ListeningEvent) error
}
