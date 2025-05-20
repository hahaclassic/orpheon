import { useState } from 'react';
import { 
  Box, 
  List, 
  ListItem, 
  ListItemText, 
  Typography, 
  IconButton, 
  Divider,
  Menu,
  MenuItem,
} from '@mui/material';
import { PlayArrow, Pause, Add as AddIcon, MoreVert } from '@mui/icons-material';
import { useNavigate } from 'react-router-dom';

interface Artist {
  id: string;
  name: string;
}

interface Track {
  id: string;
  name: string;
  duration: number;
  track_number: number;
  artists: Artist[];
  album_id: string;
}

interface TrackListProps {
  tracks: Track[];
  currentTrackId?: string;
  isPlaying?: boolean;
  onTrackClick: (trackId: string) => void;
  onAddToPlaylist?: (event: React.MouseEvent<HTMLElement>, track: Track) => void;
  showTrackNumber?: boolean;
  showAlbumLink?: boolean;
  albumId?: string;
}

const formatDuration = (seconds: number) => {
  const minutes = Math.floor(seconds / 60);
  const remainingSeconds = seconds % 60;
  return `${minutes}:${remainingSeconds.toString().padStart(2, '0')}`;
};

const TrackList = ({
  tracks,
  currentTrackId,
  isPlaying,
  onTrackClick,
  onAddToPlaylist,
  showTrackNumber = true,
  showAlbumLink = false,
  albumId,
}: TrackListProps) => {
  const navigate = useNavigate();
  const [trackMenuAnchorEl, setTrackMenuAnchorEl] = useState<null | HTMLElement>(null);
  const [selectedTrackForMenu, setSelectedTrackForMenu] = useState<Track | null>(null);

  const handleTrackMenuOpen = (event: React.MouseEvent<HTMLElement>, track: Track) => {
    event.stopPropagation();
    setSelectedTrackForMenu(track);
    setTrackMenuAnchorEl(event.currentTarget);
  };

  const handleTrackMenuClose = () => {
    setTrackMenuAnchorEl(null);
    setSelectedTrackForMenu(null);
  };

  return (
    <>
      <List>
        {tracks.map((track, index) => (
          <Box key={track.id}>
            <ListItem
              sx={{
                cursor: 'pointer',
                '&:hover': {
                  backgroundColor: 'action.hover',
                },
              }}
              onClick={() => onTrackClick(track.id)}
            >
              <Box
                sx={{
                  width: 40,
                  height: 40,
                  display: 'flex',
                  alignItems: 'center',
                  justifyContent: 'center',
                  mr: 2,
                  position: 'relative',
                  '&:hover .play-icon': {
                    opacity: 1,
                  },
                  '&:hover .track-number': {
                    opacity: 0,
                  },
                }}
              >
                {showTrackNumber && (
                  <Typography
                    className="track-number"
                    sx={{
                      position: 'absolute',
                      transition: 'opacity 0.2s',
                    }}
                  >
                    {track.track_number || index + 1}
                  </Typography>
                )}
                <IconButton
                  size="small"
                  className="play-icon"
                  sx={{
                    position: 'absolute',
                    opacity: currentTrackId === track.id ? 1 : 0,
                    transition: 'opacity 0.2s',
                  }}
                >
                  {currentTrackId === track.id && isPlaying ? <Pause /> : <PlayArrow />}
                </IconButton>
              </Box>
              <ListItemText
                primary={track.name}
                secondary={
                  <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                    <Typography variant="body2" color="text.secondary">
                      {formatDuration(track.duration)}
                    </Typography>
                    <Typography variant="body2" color="text.secondary">
                      •
                    </Typography>
                    <Box sx={{ display: 'flex', gap: 0.5 }}>
                      {track.artists.map((artist) => (
                        <Typography
                          key={artist.id}
                          variant="body2"
                          component="a"
                          href={`/artists/${artist.id}`}
                          onClick={(e) => {
                            e.stopPropagation();
                            navigate(`/artists/${artist.id}`);
                          }}
                          sx={{
                            color: 'text.secondary',
                            textDecoration: 'none',
                            '&:hover': {
                              textDecoration: 'underline',
                            },
                            '&:not(:last-child)::after': {
                              content: '", "',
                              color: 'text.secondary',
                            }
                          }}
                        >
                          {artist.name}
                        </Typography>
                      ))}
                    </Box>
                  </Box>
                }
              />
              {onAddToPlaylist && (
                <IconButton
                  onClick={(e) => onAddToPlaylist(e, track)}
                  size="small"
                >
                  <AddIcon />
                </IconButton>
              )}
              <IconButton
                onClick={(e) => handleTrackMenuOpen(e, track)}
                size="small"
              >
                <MoreVert />
              </IconButton>
            </ListItem>
            {index < tracks.length - 1 && <Divider />}
          </Box>
        ))}
      </List>

      <Menu
        anchorEl={trackMenuAnchorEl}
        open={Boolean(trackMenuAnchorEl)}
        onClose={handleTrackMenuClose}
        PaperProps={{
          sx: { 
            '& .MuiList-root': {
              padding: 0
            }
          }
        }}
        MenuListProps={{
          sx: {
            padding: 0
          }
        }}
      >
        {showAlbumLink && albumId && (
          <MenuItem 
            onClick={() => {
              handleTrackMenuClose();
              navigate(`/albums/${albumId}`);
            }}
          >
            Перейти к альбому
          </MenuItem>
        )}
      </Menu>
    </>
  );
};

export default TrackList; 