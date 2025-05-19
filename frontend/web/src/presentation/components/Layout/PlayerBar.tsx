import { Box, IconButton, Typography, Slider } from '@mui/material';
import {
  PlayArrow,
  Pause,
  SkipNext,
  SkipPrevious,
  VolumeUp,
  VolumeOff,
} from '@mui/icons-material';
import { usePlayerContext } from '../../contexts/PlayerContext';
import { useCallback, useEffect, useState } from 'react';
import type { SyntheticEvent } from 'react';

const PlayerBar = () => {
  const {
    state: { currentTrack, isPlaying, volume, progress },
    togglePlay,
    playNext,
    playPrevious,
    setVolume,
    setProgress,
  } = usePlayerContext();

  const [isDragging, setIsDragging] = useState(false);

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
            {currentTrack?.track_number || ''}
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
      </Box>
    </Box>
  );
};

export default PlayerBar; 