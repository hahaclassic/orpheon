package publisher

import (
	"context"
	"errors"

	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
	"github.com/hahaclassic/orpheon/backend/pkg/errwrap"
)

const (
	MinTotalDuration = 10
)

var (
	ErrShortListeningTime    = errors.New("error: the listening time is too short")
	ErrPublishListeningEvent = errors.New("publish listening event error")
)

type EventBus interface {
	Publish(ctx context.Context, event *entity.ListeningEvent) error
}

type ListeningEventPublisher struct {
	bus EventBus
}

func New(bus EventBus) *ListeningEventPublisher {
	return &ListeningEventPublisher{bus: bus}
}

func (p *ListeningEventPublisher) PublishListeningEvent(ctx context.Context, event *entity.ListeningEvent) error {
	total := totalDuration(event)
	if total < MinTotalDuration {
		return errwrap.Wrap(ErrPublishListeningEvent, ErrShortListeningTime)
	}

	if err := p.bus.Publish(ctx, event); err != nil {
		return errwrap.Wrap(ErrPublishListeningEvent, err)
	}

	return nil
}

func totalDuration(event *entity.ListeningEvent) int {
	total := 0
	for _, r := range event.Ranges {
		total += r.End - r.Start
	}

	return total
}
