package deleter

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
	usecase "github.com/hahaclassic/orpheon/backend/internal/domain/usecases/content/track"
	"github.com/hahaclassic/orpheon/backend/pkg/errwrap"
)

var (
	ErrForbidden = errors.New("permission denied error")
)

type PlaylistTrackDeleter interface {
	GetPlaylistIDsWithTrack(ctx context.Context, claims *entity.Claims, trackID uuid.UUID) ([]uuid.UUID, error)
	DeleteTrackFromAllPlaylists(ctx context.Context, claims *entity.Claims, trackID uuid.UUID) error
	RestoreTrackInAllPlaylists(ctx context.Context, claims *entity.Claims, trackID uuid.UUID, playlists []uuid.UUID) error
}

type AudioFileDeleter interface {
	GetAudioFile(ctx context.Context, trackID uuid.UUID) (*entity.AudioChunk, error)
	DeleteAudioFile(ctx context.Context, trackID uuid.UUID) error
	RestoreAudioFile(ctx context.Context, audio *entity.AudioChunk) error
}

type TrackMetaDeleter interface {
	DeleteMeta(ctx context.Context, trackID uuid.UUID) error
}

type TrackDeleter struct {
	audio    AudioFileDeleter
	playlist PlaylistTrackDeleter
	meta     TrackMetaDeleter
}

type OptionFunc func(*TrackDeleter)

func WithMetaDeletion(metaService TrackMetaDeleter) OptionFunc {
	return func(t *TrackDeleter) {
		t.meta = metaService
	}
}

func WithTracksDeletion(playlistTrackService PlaylistTrackDeleter) OptionFunc {
	return func(t *TrackDeleter) {
		t.playlist = playlistTrackService
	}
}

func WithFavoritesDeletion(audioService AudioFileDeleter) OptionFunc {
	return func(t *TrackDeleter) {
		t.audio = audioService
	}
}

func New(options ...OptionFunc) *TrackDeleter {
	deleter := &TrackDeleter{}
	for _, configure := range options {
		configure(deleter)
	}

	return deleter
}

type rollback func() error

func (t *TrackDeleter) DeleteTrack(ctx context.Context, claims *entity.Claims, trackID uuid.UUID) (err error) {
	var rollbacks []rollback

	defer func() {
		if err != nil {
			err = errwrap.Wrap(usecase.ErrDeleteTrack, err)

			for i := len(rollbacks) - 1; i >= 0; i-- {
				_ = rollbacks[i]()
			}
		}
	}()

	if claims.AccessLvl != entity.Admin {
		return ErrForbidden
	}

	if t.playlist != nil {
		rollback, err := t.delete(ctx, claims, playlistID)
		if err != nil {
			return err
		}
		rollbacks = append(rollbacks, rollback)
	}

	if p.cover != nil {
		rollback, err := p.deleteCover(ctx, claims, playlistID)
		if err != nil {
			return err
		}
		rollbacks = append(rollbacks, rollback)
	}

	if p.tracks != nil {
		rollback, err := p.deleteAllTracks(ctx, claims, playlistID)
		if err != nil {
			return err
		}
		rollbacks = append(rollbacks, rollback)
	}

	return nil
}

func (p *TrackDeleter) deleteTrackFromAllPlaylists(ctx context.Context, claims *entity.Claims, trackID uuid.UUID) (rollback, error) {
	trackIDs, err := p.playlist.GetPlaylistIDsWithTrack(ctx, claims, trackID)
	if err != nil {
		return nil, err
	}

	err = p.tracks.DeleteAllTracks(ctx, claims, playlistID)
	if err != nil {
		return nil, err
	}

	return func() error {
		return p.tracks.RestoreAllTracks(ctx, claims, playlistID, trackIDs)
	}, nil
}

func (p *PlaylistDeleter) deleteCover(ctx context.Context, claims *entity.Claims, playlistID uuid.UUID) (rollback, error) {
	cover, err := p.cover.GetCover(ctx, claims, playlistID)
	if err != nil {
		return nil, err
	}

	err = p.cover.DeleteCover(ctx, claims, playlistID)
	if err != nil {
		return nil, err
	}

	return func() error {
		return p.cover.SaveCover(ctx, claims, cover)
	}, nil
}

func (p *PlaylistDeleter) deleteMeta(ctx context.Context, claims *entity.Claims, playlistID uuid.UUID) error {
	return p.meta.DeleteMeta(ctx, claims, playlistID)
}
