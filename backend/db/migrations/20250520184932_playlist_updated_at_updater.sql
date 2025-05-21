-- +goose Up
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION update_playlist_updated_at() 
RETURNS TRIGGER AS $$
BEGIN
    UPDATE playlists
    SET last_updated = CURRENT_TIMESTAMP
    WHERE id = NEW.playlist_id;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql; 

CREATE TRIGGER trigger_update_playlist_updated_at
AFTER INSERT OR UPDATE or Delete ON playlist_tracks
FOR EACH ROW
EXECUTE FUNCTION update_playlist_updated_at();

drop trigger trigger_update_playlist_updated_at on playlist_tracks;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS playlist_updated_at_updater ON playlist_tracks;
DROP FUNCTION IF EXISTS update_playlist_timestamp();
-- +goose StatementEnd
