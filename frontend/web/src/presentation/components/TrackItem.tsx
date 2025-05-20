import { useState } from 'react';
import { 
  Box, 
  Typography, 
  IconButton,
  Menu,
  MenuItem,
  Tooltip,
  Link,
} from '@mui/material';
import { PlayArrow, Pause, Add as AddIcon, MoreVert, Headphones } from '@mui/icons-material';
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
  coverUrl?: string;
  license?: {
    title: string;
    description: string;
    url: string;
  };
  total_streams?: number;
}

interface TrackItemProps {
  track: Track;
  index: number;
  currentTrackId?: string;
  isPlaying?: boolean;
  onTrackClick: (trackId: string) => void;
  onAddToPlaylist?: (event: React.MouseEvent<HTMLElement>, track: Track) => void;
  showTrackNumber?: boolean;
  showAlbumLink?: boolean;
}

const formatDuration = (seconds: number) => {
  const minutes = Math.floor(seconds / 60);
  const remainingSeconds = seconds % 60;
  return `${minutes}:${remainingSeconds.toString().padStart(2, '0')}`;
};

const formatNumber = (num: number): string => {
  if (num >= 1000000) {
    return `${(num / 1000000).toFixed(1)}M`;
  }
  if (num >= 1000) {
    return `${(num / 1000).toFixed(1)}K`;
  }
  return num.toString();
};

const TrackItem = ({
  track,
  index,
  currentTrackId,
  isPlaying,
  onTrackClick,
  onAddToPlaylist,
  showTrackNumber = true,
  showAlbumLink = true,
}: TrackItemProps) => {
  const navigate = useNavigate();
  const [menuAnchorEl, setMenuAnchorEl] = useState<null | HTMLElement>(null);

  const handleMenuOpen = (event: React.MouseEvent<HTMLElement>) => {
    event.stopPropagation();
    setMenuAnchorEl(event.currentTarget);
  };

  const handleMenuClose = () => {
    setMenuAnchorEl(null);
  };

  return (
    <Box
      sx={{
        display: 'flex',
        alignItems: 'center',
        p: 1,
        cursor: 'pointer',
        '&:hover': {
          backgroundColor: 'action.hover',
        },
      }}
      onClick={() => onTrackClick(track.id)}
    >
      <Box
        sx={{
          width: 56,
          height: 56,
          position: 'relative',
          mr: 2,
          flexShrink: 0,
          '&:hover .play-icon': {
            opacity: 1,
          },
        }}
      >
        {track.coverUrl ? (
          <Box
            component="img"
            src={track.coverUrl}
            alt={track.name}
            sx={{
              width: '100%',
              height: '100%',
              borderRadius: 1,
              objectFit: 'cover',
            }}
            onError={(e) => {
              console.error('Error loading track cover:', track.coverUrl);
              e.currentTarget.style.display = 'none';
            }}
          />
        ) : (
          <Box
            sx={{
              width: '100%',
              height: '100%',
              borderRadius: 1,
              bgcolor: 'action.hover',
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'center',
            }}
          >
            <Typography variant="body2" color="text.secondary">
              {track.track_number || index + 1}
            </Typography>
          </Box>
        )}
        <IconButton
          size="small"
          className="play-icon"
          sx={{
            position: 'absolute',
            top: '50%',
            left: '50%',
            transform: 'translate(-50%, -50%)',
            opacity: currentTrackId === track.id ? 1 : 0,
            transition: 'opacity 0.2s',
            bgcolor: 'rgba(0, 0, 0, 0.5)',
            '&:hover': {
              bgcolor: 'rgba(0, 0, 0, 0.7)',
            },
          }}
        >
          {currentTrackId === track.id && isPlaying ? <Pause /> : <PlayArrow />}
        </IconButton>
      </Box>

      <Box sx={{ flex: 1, minWidth: 0 }}>
        <Typography
          variant="body1"
          sx={{
            fontWeight: 500,
            whiteSpace: 'nowrap',
            overflow: 'hidden',
            textOverflow: 'ellipsis',
          }}
        >
          {track.name}
        </Typography>
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
          {track.license && (
            <>
              <Typography variant="body2" color="text.secondary">
                •
              </Typography>
              <Tooltip title={track.license.description} arrow>
                <Link
                  href={track.license.url}
                  target="_blank"
                  rel="noopener noreferrer"
                  onClick={(e) => e.stopPropagation()}
                  sx={{
                    color: 'text.secondary',
                    textDecoration: 'none',
                    '&:hover': {
                      textDecoration: 'underline',
                    },
                  }}
                >
                  <Typography variant="body2" color="text.secondary">
                    {track.license.title}
                  </Typography>
                </Link>
              </Tooltip>
            </>
          )}
        </Box>
      </Box>

      {track.total_streams !== undefined && (
        <Box
          sx={{
            display: 'flex',
            alignItems: 'center',
            gap: 0.5,
            ml: 2,
            mr: 1,
            color: 'text.secondary',
          }}
        >
          <Headphones sx={{ fontSize: 16 }} />
          <Typography variant="body2" color="text.secondary">
            {formatNumber(track.total_streams)}
          </Typography>
        </Box>
      )}

      {onAddToPlaylist && (
        <IconButton
          onClick={(e) => onAddToPlaylist(e, track)}
          size="small"
          sx={{ ml: 1 }}
        >
          <AddIcon />
        </IconButton>
      )}

      <IconButton
        onClick={handleMenuOpen}
        size="small"
        sx={{ ml: 1 }}
      >
        <MoreVert />
      </IconButton>

      <Menu
        anchorEl={menuAnchorEl}
        open={Boolean(menuAnchorEl)}
        onClose={handleMenuClose}
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
        {showAlbumLink && (
          <MenuItem 
            onClick={() => {
              handleMenuClose();
              navigate(`/albums/${track.album_id}`);
            }}
          >
            Перейти к альбому
          </MenuItem>
        )}
      </Menu>
    </Box>
  );
};

export default TrackItem; 