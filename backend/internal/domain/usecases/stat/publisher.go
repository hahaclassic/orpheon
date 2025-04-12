package stats

import (
	"context"

	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
)

type ListeningEventPublisher interface {
	SendListeningEvent(ctx context.Context, event *entity.ListeningEvent) error
}
