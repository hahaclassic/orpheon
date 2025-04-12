package consumer

import (
	"context"
	"errors"

	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
	"github.com/hahaclassic/orpheon/backend/pkg/errwrap"
)

var (
	ErrSetupConsumer         = errors.New("setup consumer error")
	ErrConsumeListeningEvent = errors.New("consume listening event")
)

type ListeningStatService interface {
	UpdateTrackStat(ctx context.Context, event *entity.ListeningEvent) error
}

type EventBus interface {
	Subscribe(ctx context.Context, handler func(ctx context.Context, event *entity.ListeningEvent) error) error
}

type ListeningEventConsumer struct {
	stat ListeningStatService
}

func Setup(ctx context.Context, bus EventBus, statService ListeningStatService) error {
	consumer := &ListeningEventConsumer{stat: statService}

	var err error

	go func() {
		err = bus.Subscribe(ctx, consumer.ConsumeListeningEvent)
	}()

	return errwrap.WrapIfErr(ErrSetupConsumer, err)
}

func (c *ListeningEventConsumer) ConsumeListeningEvent(ctx context.Context, event *entity.ListeningEvent) error {
	if err := c.stat.UpdateTrackStat(ctx, event); err != nil {
		return errwrap.Wrap(ErrConsumeListeningEvent, err)
	}

	return nil
}
