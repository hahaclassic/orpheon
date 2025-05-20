import { List, Divider } from '@mui/material';
import TrackItem from './TrackItem';

interface Artist {
  id: string;
  name: string;
  description?: string;
  country?: string;
}

interface Track {
  id: string;
  name: string;
  duration: number;
  track_number: number;
  artists: Artist[];
  album_id: string;
  coverUrl?: string;
  license?: {
    id: string;
    title: string;
    description: string;
    url: string;
  };
  total_streams?: number;
}

interface TrackListProps {
  tracks: Track[];
  currentTrackId?: string;
  isPlaying?: boolean;
  onTrackClick: (trackId: string) => void;
  onAddToPlaylist?: (event: React.MouseEvent<HTMLElement>, track: Track) => void;
  showTrackNumber?: boolean;
  showAlbumLink?: boolean;
}

const TrackList = ({
  tracks,
  currentTrackId,
  isPlaying,
  onTrackClick,
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
            index={index}
            currentTrackId={currentTrackId}
            isPlaying={isPlaying}
            onTrackClick={onTrackClick}
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