import { List, Divider } from '@mui/material';
import TrackItem from './TrackItem';
import type { Track } from '../types';

interface TrackListProps {
  tracks: Track[];
  onAddToPlaylist?: (event: React.MouseEvent<HTMLElement>, track: Track) => void;
  showTrackNumber?: boolean;
  showAlbumLink?: boolean;
}

const TrackList = ({
  tracks,
  onAddToPlaylist,
  showTrackNumber = true,
  showAlbumLink = true,
}: TrackListProps) => {
  return (
    <List>
      {tracks.map((track, index) => (
        <div key={track.id}>
          <TrackItem
            track={track}
            tracks={tracks}
            index={index}
            onAddToPlaylist={onAddToPlaylist}
            showTrackNumber={showTrackNumber}
            showAlbumLink={showAlbumLink}
          />
          {index < tracks.length - 1 && <Divider />}
        </div>
      ))}
    </List>
  );
};

export default TrackList; 