import { Box, Typography, Card, CardContent, CardMedia, styled } from '@mui/material';
import { useEffect } from 'react';
import { useApi } from '../../hooks/useApi';
import { apiService } from '../../../core/infrastructure/services/api';
import LoadingSpinner from '../../components/LoadingSpinner';
import ErrorBoundary from '../../components/ErrorBoundary';
import type { Playlist, Track } from '../../../core/infrastructure/services/api';

const Section = styled(Box)({
  marginBottom: '40px',
});

const StyledCard = styled(Card)({
  backgroundColor: 'background.paper',
  transition: 'transform 0.2s ease-in-out',
  '&:hover': {
    transform: 'scale(1.02)',
    cursor: 'pointer',
  },
});

const Home = () => {
  const {
    data: featuredPlaylists,
    loading: playlistsLoading,
    error: playlistsError,
    execute: fetchPlaylists,
  } = useApi<Playlist[]>();

  const {
    data: recentlyPlayed,
    loading: tracksLoading,
    error: tracksError,
    execute: fetchRecentTracks,
  } = useApi<Track[]>();

  useEffect(() => {
    fetchPlaylists(apiService.getPlaylists());
    fetchRecentTracks(apiService.getTracks());
  }, [fetchPlaylists, fetchRecentTracks]);

  if (playlistsLoading || tracksLoading) {
    return <LoadingSpinner />;
  }

  if (playlistsError || tracksError) {
    return (
      <Box sx={{ p: 3, textAlign: 'center' }}>
        <Typography color="error">
          {playlistsError?.message || tracksError?.message || 'Failed to load content'}
        </Typography>
      </Box>
    );
  }

  return (
    <ErrorBoundary>
      <Box>
        <Section>
          <Typography variant="h4" sx={{ mb: 3 }}>
            Good Evening
          </Typography>
          <Box sx={{ display: 'grid', gridTemplateColumns: { xs: '1fr', sm: 'repeat(2, 1fr)', md: 'repeat(4, 1fr)' }, gap: 3 }}>
            {featuredPlaylists?.map((playlist) => (
              <Box key={playlist.id}>
                <StyledCard>
                  <CardMedia
                    component="img"
                    height="200"
                    image={playlist.coverUrl}
                    alt={playlist.title}
                  />
                  <CardContent>
                    <Typography variant="h6" noWrap>
                      {playlist.title}
                    </Typography>
                    <Typography variant="body2" color="text.secondary" noWrap>
                      {playlist.description}
                    </Typography>
                  </CardContent>
                </StyledCard>
              </Box>
            ))}
          </Box>
        </Section>

        <Section>
          <Typography variant="h5" sx={{ mb: 3 }}>
            Recently Played
          </Typography>
          <Box sx={{ display: 'grid', gap: 2 }}>
            {recentlyPlayed?.map((track) => (
              <Box key={track.id}>
                <StyledCard>
                  <Box sx={{ display: 'flex', alignItems: 'center', p: 2 }}>
                    <CardMedia
                      component="img"
                      sx={{ width: 60, height: 60, borderRadius: 1 }}
                      image={track.coverUrl}
                      alt={track.title}
                    />
                    <Box sx={{ ml: 2 }}>
                      <Typography variant="subtitle1">{track.title}</Typography>
                      <Typography variant="body2" color="text.secondary">
                        {track.artist}
                      </Typography>
                    </Box>
                  </Box>
                </StyledCard>
              </Box>
            ))}
          </Box>
        </Section>
      </Box>
    </ErrorBoundary>
  );
};

export default Home; 