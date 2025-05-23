-- +goose Up
-- +goose StatementBegin
CREATE OR REPLACE PROCEDURE change_track_position(
    p_playlist_id UUID,
    p_track_id UUID,
    p_new_position INT
)
LANGUAGE plpgsql
AS $$
DECLARE
    current_position INT;
BEGIN
    SELECT position INTO current_position
    FROM playlist_tracks
    WHERE playlist_id = p_playlist_id AND track_id = p_track_id;

    IF current_position IS NULL THEN
        RAISE EXCEPTION 'Track not found in playlist';
    END IF;

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
$$;


CREATE OR REPLACE PROCEDURE delete_track_from_playlist(p_playlist_id UUID, p_track_id UUID)
LANGUAGE plpgsql
AS $$
DECLARE
    deleted_position INT;
BEGIN
    SELECT position INTO deleted_position
    FROM playlist_tracks
    WHERE playlist_id = p_playlist_id AND track_id = p_track_id;

    IF NOT FOUND THEN
        RAISE EXCEPTION 'Track with id % not found in playlist %', p_track_id, p_playlist_id;
    END IF;

    DELETE FROM playlist_tracks
    WHERE playlist_id = p_playlist_id AND track_id = p_track_id;

    UPDATE playlist_tracks
    SET position = position + 1000000
    WHERE playlist_id = p_playlist_id AND position > deleted_position;

    UPDATE playlist_tracks
    SET position = position - 1000001
    WHERE playlist_id = p_playlist_id AND position > 1000000;
END;
$$;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP PROCEDURE change_track_position(UUID, UUID, INT);
DROP PROCEDURE delete_track_from_playlist(UUID, UUID);
-- +goose StatementEnd
