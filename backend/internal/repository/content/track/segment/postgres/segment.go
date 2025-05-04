package segment_postgres

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TrackSegmentRepository struct {
	pool *pgxpool.Pool
}

func NewTrackSegmentRepository(pool *pgxpool.Pool) *TrackSegmentRepository {
	return &TrackSegmentRepository{pool: pool}
}

func (r *TrackSegmentRepository) IncrementSegmentStreamCount(ctx context.Context, trackID uuid.UUID, segmentsIdxs []int) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer func() {
		err := tx.Rollback(ctx)
		if err != nil {
			slog.Error("err", "rollback err", err)
		}
	}()

	for _, idx := range segmentsIdxs {
		_, err := tx.Exec(ctx, `
			UPDATE track_segments
			SET stream_count = stream_count + 1
			WHERE track_id = $1 AND idx = $2
		`, trackID, idx)

		if err != nil {
			return fmt.Errorf("increment stream count: %w", err)
		}
	}

	return tx.Commit(ctx)
}

func (r *TrackSegmentRepository) GetSegments(ctx context.Context, trackID uuid.UUID) ([]*entity.Segment, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT idx, stream_count, range_start, range_end
		FROM track_segments
		WHERE track_id = $1
		ORDER BY idx
	`, trackID)
	if err != nil {
		return nil, fmt.Errorf("get segments: %w", err)
	}
	defer rows.Close()

	var segments []*entity.Segment

	for rows.Next() {
		var seg entity.Segment
		var start, end int
		err := rows.Scan(&seg.Idx, &seg.StreamCount, &start, &end)
		if err != nil {
			return nil, fmt.Errorf("scan segment: %w", err)
		}
		seg.TrackID = trackID
		seg.Range = &entity.Range{Start: start, End: end}
		segments = append(segments, &seg)
	}

	return segments, nil
}

func (r *TrackSegmentRepository) DeleteSegments(ctx context.Context, trackID uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `
		DELETE FROM track_segments WHERE track_id = $1
	`, trackID)
	if err != nil {
		return fmt.Errorf("delete segments: %w", err)
	}
	return nil
}

func (r *TrackSegmentRepository) CreateSegments(ctx context.Context, trackID uuid.UUID, numOfSegments int) error {
	if numOfSegments <= 0 {
		return fmt.Errorf("number of segments must be positive")
	}

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer func() {
		err := tx.Rollback(ctx)
		if err != nil {
			slog.Error("err", "rollback err", err)
		}
	}()

	for i := 0; i < numOfSegments; i++ {
		// здесь диапазон можно задать произвольно — условно по 10 сек или т.п.
		start := i * 10
		end := start + 10

		_, err := tx.Exec(ctx, `
			INSERT INTO track_segments (track_id, idx, stream_count, range_start, range_end)
			VALUES ($1, $2, 0, $3, $4)
		`, trackID, i, start, end)
		if err != nil {
			return fmt.Errorf("create segment %d: %w", i, err)
		}
	}

	return tx.Commit(ctx)
}
