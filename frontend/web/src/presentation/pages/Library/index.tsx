import { useState, useEffect } from 'react';
import { Box, Typography, Card, CardContent, CardMedia, Tabs, Tab, styled } from '@mui/material';
import { useApi } from '../../hooks/useApi';
import { apiService } from '../../../core/infrastructure/services/api';
import LoadingSpinner from '../../components/LoadingSpinner';
import ErrorBoundary from '../../components/ErrorBoundary';
import type { Playlist, Track } from '../../../core/infrastructure/services/api';

const LibraryContainer = styled(Box)({
  padding: '24px',
});

const StyledCard = styled(Card)({
  backgroundColor: 'background.paper',
  transition: 'transform 0.2s ease-in-out',
  '&:hover': {
    transform: 'scale(1.02)',
    cursor: 'pointer',
  },
});

interface TabPanelProps {
  children?: React.ReactNode;
  index: number;
  value: number;
}

const TabPanel = (props: TabPanelProps) => {
  const { children, value, index, ...other } = props;

  return (
    <div
      role="tabpanel"
      hidden={value !== index}
      id={`library-tabpanel-${index}`}
      aria-labelledby={`library-tab-${index}`}
      {...other}
    >
      {value === index && <Box sx={{ py: 3 }}>{children}</Box>}
    </div>
  );
};

const Library = () => {
  const [activeTab, setActiveTab] = useState(0);

  const {
    data: playlists,
    loading: playlistsLoading,
    error: playlistsError,
    execute: fetchPlaylists,
  } = useApi<Playlist[]>();

  const {
    data: likedTracks,
    loading: tracksLoading,
    error: tracksError,
    execute: fetchLikedTracks,
  } = useApi<Track[]>();

  useEffect(() => {
    fetchPlaylists(apiService.getPlaylists());
    fetchLikedTracks(apiService.getLikedTracks());
  }, [fetchPlaylists, fetchLikedTracks]);

  const handleTabChange = (_: React.SyntheticEvent, newValue: number) => {
    setActiveTab(newValue);
  };

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
      <LibraryContainer>
        <Typography variant="h4" sx={{ mb: 3 }}>
          Your Library
        </Typography>

        <Box sx={{ borderBottom: 1, borderColor: 'divider' }}>
          <Tabs value={activeTab} onChange={handleTabChange}>
            <Tab label="Playlists" />
            <Tab label="Liked Songs" />
          </Tabs>
        </Box>

        <TabPanel value={activeTab} index={0}>
          <Box sx={{ display: 'grid', gridTemplateColumns: { xs: '1fr', sm: 'repeat(2, 1fr)', md: 'repeat(3, 1fr)' }, gap: 3 }}>
            {playlists?.map((playlist) => (
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
                    <Typography variant="body2" color="text.secondary">
                      {playlist.tracks.length} tracks
                    </Typography>
                  </CardContent>
                </StyledCard>
              </Box>
            ))}
          </Box>
        </TabPanel>

        <TabPanel value={activeTab} index={1}>
          <Box sx={{ display: 'grid', gap: 2 }}>
            {likedTracks?.map((track) => (
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
                      {track.album && (
                        <Typography variant="body2" color="text.secondary">
                          {track.album}
                        </Typography>
                      )}
                    </Box>
                  </Box>
                </StyledCard>
              </Box>
            ))}
          </Box>
        </TabPanel>
      </LibraryContainer>
    </ErrorBoundary>
  );
};

export default Library; 