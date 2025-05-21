import { Box, IconButton, Typography, Slider, Menu, MenuItem } from '@mui/material';
import {
  PlayArrow,
  Pause,
  SkipNext,
  SkipPrevious,
  VolumeUp,
  VolumeOff,
  Add as AddIcon,
  MoreVert,
  Check,
} from '@mui/icons-material';
import { usePlayerContext } from '../../contexts/PlayerContext';
import { useCallback, useEffect, useState } from 'react';
import type { SyntheticEvent } from 'react';
import { useNavigate } from 'react-router-dom';
import { apiService } from '../../services/api';

const PlayerBar = () => {
  const navigate = useNavigate();
  const {
    state: { currentTrack, isPlaying, volume, progress, duration },
    togglePlay,
    playNext,
    playPrevious,
    setVolume,
    setProgress,
  } = usePlayerContext();

  const [isDragging, setIsDragging] = useState(false);
  const [menuAnchorEl, setMenuAnchorEl] = useState<null | HTMLElement>(null);
  const [playlists, setPlaylists] = useState<Array<{ id: string; name: string; tracks: any[] }>>([]);

  useEffect(() => {
    const fetchPlaylists = async () => {
      try {
        const response = await apiService.get('/me/playlists');
        setPlaylists(response);
      } catch (err) {
        console.error('Error fetching playlists:', err);
      }
    };

    fetchPlaylists();
  }, []);

  const handleProgressChange = useCallback(
    (_: Event | SyntheticEvent, newValue: number | number[]) => {
      if (!isDragging) {
        setProgress(newValue as number);
      }
    },
    [isDragging, setProgress]
  );

  const handleProgressChangeCommitted = useCallback(
    (_: Event | SyntheticEvent, newValue: number | number[]) => {
      setIsDragging(false);
      setProgress(newValue as number);
    },
    [setProgress]
  );

  const handleVolumeChange = useCallback(
    (_: Event | SyntheticEvent, newValue: number | number[]) => {
      setVolume(newValue as number);
    },
    [setVolume]
  );

  const toggleMute = useCallback(() => {
    setVolume(volume === 0 ? 1 : 0);
  }, [volume, setVolume]);

  const handleMenuOpen = (event: React.MouseEvent<HTMLElement>) => {
    setMenuAnchorEl(event.currentTarget);
  };

  const handleMenuClose = () => {
    setMenuAnchorEl(null);
  };

  const handleGoToAlbum = () => {
    if (currentTrack?.album_id) {
      navigate(`/albums/${currentTrack.album_id}`);
      handleMenuClose();
    }
  };

  const handlePlaylistSelect = async (playlistId: string) => {
    if (currentTrack) {
      try {
        const playlist = playlists.find(p => p.id === playlistId);
        const isInPlaylist = playlist && isTrackInPlaylist(playlist, currentTrack.id);

        if (isInPlaylist) {
          await apiService.delete(`/playlists/${playlistId}/tracks/${currentTrack.id}`);
          setPlaylists(playlists.map(p => {
            if (p.id === playlistId) {
              return {
                ...p,
                tracks: (p.tracks || []).filter(t => t.id !== currentTrack.id)
              };
            }
            return p;
          }));
        } else {
          await apiService.post(`/playlists/${playlistId}/tracks`, { track_id: currentTrack.id });
          setPlaylists(playlists.map(p => {
            if (p.id === playlistId) {
              return {
                ...p,
                tracks: [...(p.tracks || []), currentTrack]
              };
            }
            return p;
          }));
        }
      } catch (err: any) {
        console.error('Error managing track in playlist:', {
          error: err,
          response: err.response?.data,
          status: err.response?.status,
          playlistId,
          trackId: currentTrack.id
        });
      }
    }
    handleMenuClose();
  };

  const isTrackInPlaylist = (playlist: { tracks: any[] }, trackId: string) => {
    return playlist.tracks?.some(track => track.id === trackId) || false;
  };

  useEffect(() => {
    const handleKeyPress = (event: KeyboardEvent) => {
      if (event.code === 'Space' && !event.repeat) {
        event.preventDefault();
        togglePlay();
      }
    };

    window.addEventListener('keydown', handleKeyPress);
    return () => window.removeEventListener('keydown', handleKeyPress);
  }, [togglePlay]);

  // Auto-play next track when current track ends
  useEffect(() => {
    if (progress >= duration && duration > 0) {
      playNext();
    }
  }, [progress, duration, playNext]);

  return (
    <Box
      sx={{
        height: { xs: '120px', sm: '100px', md: '90px' },
        backgroundColor: 'background.paper',
        borderTop: '1px solid rgba(255, 255, 255, 0.1)',
        p: 2,
        display: 'flex',
        alignItems: 'center',
        gap: 2,
        width: '100%',
        minWidth: 0,
        borderRadius: '12px 12px 0 0',
        boxShadow: '0 -4px 6px -1px rgba(0, 0, 0, 0.1)',
      }}
    >
      <Box
        sx={{
          display: 'flex',
          minWidth: 0,
          flex: { xs: '1 1 100%', sm: '1 1 30%', md: '1 1 25%' },
          flexDirection: { xs: 'column', sm: 'row' },
          alignItems: { xs: 'flex-start', sm: 'center' },
          gap: { xs: 1, sm: 2 },
          maxWidth: { sm: '30%', md: '25%' },
        }}
      >
        {currentTrack?.coverUrl && (
          <Box
            component="img"
            src={currentTrack.coverUrl}
            alt={currentTrack.name}
            sx={{
              width: { xs: 40, sm: 56 },
              height: { xs: 40, sm: 56 },
              borderRadius: 1,
              objectFit: 'cover',
              flexShrink: 0,
            }}
          />
        )}
        <Box sx={{ minWidth: 0, display: 'flex', flexDirection: 'column', flex: 1 }}>
          <Typography
            variant="subtitle1"
            sx={{
              fontWeight: 500,
              whiteSpace: 'nowrap',
              overflow: 'hidden',
              textOverflow: 'ellipsis',
            }}
          >
            {currentTrack?.name || 'Нет трека'}
          </Typography>
          <Typography
            variant="body2"
            sx={{
              color: 'text.secondary',
              fontSize: '0.875rem',
              whiteSpace: 'nowrap',
              overflow: 'hidden',
              textOverflow: 'ellipsis',
            }}
          >
            {currentTrack?.artists?.map(artist => artist.name).join(', ') || ''}
          </Typography>
        </Box>
      </Box>

      <Box
        sx={{
          display: 'flex',
          alignItems: 'center',
          gap: 2,
          flex: { xs: '1 1 100%', sm: '1 1 40%', md: '1 1 50%' },
          justifyContent: 'center',
          order: { xs: 3, sm: 'unset' },
          width: { xs: '100%', sm: 'auto' },
          minWidth: 0,
        }}
      >
        <IconButton onClick={playPrevious} size="small" disabled={!currentTrack}>
          <SkipPrevious />
        </IconButton>
        <IconButton onClick={togglePlay} size="large" disabled={!currentTrack}>
          {isPlaying ? <Pause /> : <PlayArrow />}
        </IconButton>
        <IconButton onClick={playNext} size="small" disabled={!currentTrack}>
          <SkipNext />
        </IconButton>
        <Box sx={{ width: '100%', maxWidth: 400, mx: 2, minWidth: 0 }}>
          <Slider
            value={progress}
            onChange={handleProgressChange}
            onChangeCommitted={handleProgressChangeCommitted}
            onMouseDown={() => setIsDragging(true)}
            aria-label="track progress"
            disabled={!currentTrack}
            max={duration}
            sx={{
              color: 'primary.main',
              height: 4,
              '& .MuiSlider-track': {
                border: 'none',
              },
              '& .MuiSlider-thumb': {
                width: 8,
                height: 8,
                '&:before': {
                  boxShadow: '0 2px 12px 0 rgba(0,0,0,0.4)',
                },
                '&:hover, &.Mui-focusVisible': {
                  boxShadow: '0px 0px 0px 8px rgb(177 156 217 / 16%)',
                },
                '&.Mui-active': {
                  width: 12,
                  height: 12,
                },
              },
              '& .MuiSlider-rail': {
                opacity: 0.28,
              },
            }}
          />
        </Box>
      </Box>

      <Box
        sx={{
          display: 'flex',
          alignItems: 'center',
          gap: 2,
          flex: { xs: '1 1 100%', sm: '1 1 30%', md: '1 1 25%' },
          justifyContent: 'flex-end',
          order: { xs: 2, sm: 'unset' },
          width: { xs: '100%', sm: 'auto' },
          minWidth: 0,
          maxWidth: { sm: '30%', md: '25%' },
        }}
      >
        <IconButton onClick={toggleMute} size="small">
          {volume === 0 ? <VolumeOff /> : <VolumeUp />}
        </IconButton>
        <Box sx={{ width: 100, display: { xs: 'none', sm: 'block' }, minWidth: 0 }}>
          <Slider
            value={volume}
            onChange={handleVolumeChange}
            min={0}
            max={1}
            step={0.01}
            aria-label="volume"
            sx={{
              color: 'primary.main',
              height: 4,
              '& .MuiSlider-track': {
                border: 'none',
              },
              '& .MuiSlider-thumb': {
                width: 8,
                height: 8,
                '&:before': {
                  boxShadow: '0 2px 12px 0 rgba(0,0,0,0.4)',
                },
                '&:hover, &.Mui-focusVisible': {
                  boxShadow: '0px 0px 0px 8px rgb(177 156 217 / 16%)',
                },
                '&.Mui-active': {
                  width: 12,
                  height: 12,
                },
              },
              '& .MuiSlider-rail': {
                opacity: 0.28,
              },
            }}
          />
        </Box>
        {currentTrack && (
          <>
            <IconButton size="small" onClick={handleMenuOpen}>
              <AddIcon />
            </IconButton>
            <IconButton size="small" onClick={handleGoToAlbum}>
              <MoreVert />
            </IconButton>
            <Menu
              anchorEl={menuAnchorEl}
              open={Boolean(menuAnchorEl)}
              onClose={handleMenuClose}
              PaperProps={{
                sx: { 
                  maxHeight: '280px',
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
              {playlists.map((playlist) => (
                <MenuItem 
                  key={playlist.id} 
                  onClick={() => handlePlaylistSelect(playlist.id)}
                  sx={{ 
                    minHeight: '40px',
                    '&:hover': {
                      backgroundColor: 'action.hover'
                    },
                    display: 'flex',
                    justifyContent: 'space-between',
                    alignItems: 'center'
                  }}
                >
                  <span>{playlist.name}</span>
                  {currentTrack && isTrackInPlaylist(playlist, currentTrack.id) && (
                    <Check sx={{ ml: 1, color: 'white' }} />
                  )}
                </MenuItem>
              ))}
            </Menu>
          </>
        )}
      </Box>
    </Box>
  );
};

export default PlayerBar; 