import { Box, IconButton, Typography, Slider } from '@mui/material';
import {
  PlayArrow,
  Pause,
  SkipNext,
  SkipPrevious,
  VolumeUp,
  VolumeOff,
  MusicNote,
} from '@mui/icons-material';
import { usePlayerContext } from '../../contexts/PlayerContext';
import { useCallback, useEffect } from 'react';
import { formatDuration } from '../../utils/time';

const PlayerBar = () => {
  const { state, controls } = usePlayerContext();
  const {
    currentTrack,
    isPlaying,
    volume,
    progress,
    duration,
  } = state;
  const {
    togglePlay,
    playNext,
    playPrevious,
    setVolume,
    setProgress,
  } = controls;

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

  // Auto-play next track when current track ends
  useEffect(() => {
    if (progress >= duration && duration > 0) {
      playNext();
    }
  }, [progress, duration, playNext]);

  return (
    <Box
      sx={{
        height: '100%',
        display: 'flex',
        alignItems: 'center',
        px: 2,
        gap: 2,
        border: '1px solid rgba(255, 255, 255, 0.1)',
        borderRadius: '12px',
        boxShadow: '0 4px 6px -1px rgba(0, 0, 0, 0.1)',
      }}
    >
      <Box sx={{ display: 'flex', alignItems: 'center', gap: 1, minWidth: 200 }}>
        {currentTrack ? (
          <>
            <Box
              component="img"
              src={currentTrack.coverUrl}
              alt={currentTrack.name}
              sx={{
                width: 56,
                height: 56,
                borderRadius: 1,
                objectFit: 'cover',
              }}
            />
            <Box>
              <Typography variant="subtitle1" noWrap>
                {currentTrack.name}
              </Typography>
              <Typography variant="body2" color="text.secondary" noWrap>
                {currentTrack.artists.map(artist => artist.name).join(', ')}
              </Typography>
            </Box>
          </>
        ) : (
          <>
            <Box
              sx={{
                width: 56,
                height: 56,
                borderRadius: 1,
                bgcolor: 'action.hover',
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'center',
              }}
            >
              <MusicNote sx={{ color: 'text.secondary' }} />
            </Box>
            <Box>
              <Typography variant="subtitle1" color="text.secondary">
                Нет активного трека
              </Typography>
              <Typography variant="body2" color="text.secondary">
                Выберите трек для воспроизведения
              </Typography>
            </Box>
          </>
        )}
      </Box>

      <Box sx={{ flex: 1, display: 'flex', flexDirection: 'column', alignItems: 'center', gap: 1 }}>
        <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
          <IconButton onClick={playPrevious} disabled={!currentTrack}>
            <SkipPrevious />
          </IconButton>
          <IconButton onClick={togglePlay} disabled={!currentTrack}>
            {isPlaying ? <Pause /> : <PlayArrow />}
          </IconButton>
          <IconButton onClick={playNext} disabled={!currentTrack}>
            <SkipNext />
          </IconButton>
        </Box>

        <Box sx={{ width: '100%', display: 'flex', alignItems: 'center', gap: 1 }}>
          <Typography variant="caption" color="text.secondary">
            {formatDuration(progress)}
          </Typography>
          <Slider
            value={progress}
            max={duration}
            onChange={(_, value) => setProgress(value as number)}
            disabled={!currentTrack}
            sx={{ flex: 1 }}
          />
          <Typography variant="caption" color="text.secondary">
            {formatDuration(duration)}
          </Typography>
        </Box>
      </Box>

      <Box sx={{ display: 'flex', alignItems: 'center', gap: 1, minWidth: 200 }}>
        <IconButton onClick={toggleMute} disabled={!currentTrack}>
          {volume === 0 ? <VolumeOff /> : <VolumeUp />}
        </IconButton>
        <Slider
          value={volume}
          min={0}
          max={1}
          step={0.01}
          onChange={(_, value) => setVolume(value as number)}
          disabled={!currentTrack}
          sx={{ width: 100 }}
        />
      </Box>
    </Box>
  );
};

export default PlayerBar; 