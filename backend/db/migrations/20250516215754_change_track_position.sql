-- +goose Up
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION change_track_position(
    p_playlist_id UUID,
    p_track_id UUID,
    p_new_position INT
) RETURNS VOID AS $$
DECLARE
    current_position INT;
BEGIN
    SELECT position INTO current_position
    FROM playlist_tracks
    WHERE playlist_id = p_playlist_id AND track_id = p_track_id;

    IF current_position = p_new_position THEN
        RETURN;
    END IF;

    UPDATE playlist_tracks
    SET position = 1000000
    WHERE playlist_id = p_playlist_id AND track_id = p_track_id;

    IF p_new_position < current_position THEN
        UPDATE playlist_tracks
        SET position = position + 1
        WHERE playlist_id = p_playlist_id
          AND position >= p_new_position
          AND position < current_position;
    ELSE
        UPDATE playlist_tracks
        SET position = position - 1
        WHERE playlist_id = p_playlist_id
          AND position <= p_new_position
          AND position > current_position;
    END IF;

    UPDATE playlist_tracks
    SET position = p_new_position
    WHERE playlist_id = p_playlist_id AND track_id = p_track_id;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP FUNCTION change_track_position(UUID, UUID, INT);
-- +goose StatementEnd
