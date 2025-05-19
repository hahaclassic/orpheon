import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import {
  Box,
  Typography,
  Card,
  CardContent,
  CardMedia,
  Grid,
  CircularProgress,
  Alert,
  Button,
  Container,
} from '@mui/material';
import AddIcon from '@mui/icons-material/Add';
import ImageIcon from '@mui/icons-material/Image';
import { apiService } from '../../../presentation/services/api';
import { useApi } from '../../../presentation/hooks/useApi';
import PlaylistDialog from '../../../presentation/components/PlaylistDialog';
import type { Track, Playlist } from '../../../presentation/types';

const Library = () => {
  const [myPlaylists, setMyPlaylists] = useState<Playlist[]>([]);
  const [favoritePlaylists, setFavoritePlaylists] = useState<Playlist[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [createDialogOpen, setCreateDialogOpen] = useState(false);
  const navigate = useNavigate();
  const { createPlaylist } = useApi();

  const fetchPlaylistCovers = async (playlists: Playlist[]) => {
    return Promise.all(
      playlists.map(async (playlist) => {
        try {
          const coverResponse = await apiService.get(`/playlists/${playlist.id}/cover`, {
            responseType: 'blob'
          });
          const coverUrl = URL.createObjectURL(coverResponse);
          return { ...playlist, coverImage: coverUrl };
        } catch (err) {
          console.error(`Error fetching cover for playlist ${playlist.id}:`, err);
          return { ...playlist, coverImage: '/default-playlist-cover.jpg' };
        }
      })
    );
  };

  const fetchData = async () => {
    try {
      setLoading(true);
      setError(null);
      const [myPlaylistsData, favoritePlaylistsData] = await Promise.all([
        apiService.get('/me/playlists'),
        apiService.get('/me/favorites'),
      ]);
      const [myPlaylistsWithCovers, favoritePlaylistsWithCovers] = await Promise.all([
        fetchPlaylistCovers(myPlaylistsData),
        fetchPlaylistCovers(favoritePlaylistsData),
      ]);
      setMyPlaylists(myPlaylistsWithCovers);
      setFavoritePlaylists(favoritePlaylistsWithCovers);
    } catch (err) {
      console.error("Error fetching data:", err);
      setError("Ошибка при загрузке данных. Пожалуйста, попробуйте позже.");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchData();
  }, []);

  const handlePlaylistClick = (playlistId: number) => {
    navigate(`/playlists/${playlistId}`);
  };

  const handleCreatePlaylist = async (name: string) => {
    try {
      await createPlaylist(name);
      fetchData();
    } catch (err) {
      console.error("Error creating playlist:", err);
      setError("Ошибка при создании плейлиста. Пожалуйста, попробуйте позже.");
    }
  };

  // Объединяем обычные и избранные плейлисты, помечая избранные
  const allPlaylists = [
    ...myPlaylists.map((p) => ({ ...p, isFavorite: false })),
    ...favoritePlaylists.map((p) => ({ ...p, isFavorite: true })),
  ];
  // Убираем дубли по id (если плейлист есть и в моих, и в избранных)
  const uniquePlaylists = allPlaylists.filter(
    (playlist, idx, arr) => arr.findIndex((p) => p.id === playlist.id) === idx
  );

  return (
    <Container maxWidth="lg" sx={{ py: 4 }}>
      {error && (
        <Alert severity="error" sx={{ mb: 2 }}>{error}</Alert>
      )}
      <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 3 }}>
        <Typography variant="h4" component="h1" gutterBottom>
          Мои плейлисты
        </Typography>
        <Button
          variant="contained"
          startIcon={<AddIcon />}
          onClick={() => setCreateDialogOpen(true)}
        >
          Создать плейлист
        </Button>
      </Box>
      {loading ? (
        <Box sx={{ display: 'flex', justifyContent: 'center', p: 3 }}>
          <CircularProgress />
        </Box>
      ) : (
        <>
          <Grid container spacing={3}>
            {uniquePlaylists.length === 0 ? (
              <Grid item xs={12}>
                <Typography variant="h6" textAlign="center" color="text.secondary">
                  У вас пока нет плейлистов
                </Typography>
              </Grid>
            ) : (
              uniquePlaylists.map((playlist) => (
                <Grid item xs={12} sm={6} md={3} key={playlist.id}>
                  <Card
                    sx={{
                      height: '100%',
                      display: 'flex',
                      flexDirection: 'column',
                      cursor: 'pointer',
                      border: playlist.isFavorite ? '2px solid #a78bfa' : undefined,
                      boxShadow: playlist.isFavorite ? '0 0 0 2px #a78bfa' : undefined,
                      '&:hover': {
                        transform: 'scale(1.02)',
                        transition: 'transform 0.2s ease-in-out',
                        boxShadow: playlist.isFavorite ? '0 0 0 4px #a78bfa' : undefined,
                      },
                    }}
                    onClick={() => handlePlaylistClick(playlist.id)}
                  >
                    {playlist.coverImage ? (
                      <CardMedia
                        component="img"
                        sx={{
                          height: 250,
                          width: '100%',
                          objectFit: 'cover',
                          aspectRatio: '1/1'
                        }}
                        image={playlist.coverImage}
                        alt={playlist.name}
                      />
                    ) : (
                      <Box
                        sx={{
                          height: 250,
                          width: '100%',
                          display: 'flex',
                          alignItems: 'center',
                          justifyContent: 'center',
                          bgcolor: 'primary.dark',
                          borderRadius: 2,
                        }}
                      >
                        <ImageIcon sx={{ fontSize: 64, color: 'primary.contrastText', opacity: 0.3 }} />
                      </Box>
                    )}
                    <CardContent>
                      <Typography gutterBottom variant="h6" component="div" noWrap>
                        {playlist.name}
                      </Typography>
                      {playlist.isFavorite && (
                        <Typography variant="body2" color="#a78bfa">
                          ★ Избранное
                        </Typography>
                      )}
                      <Typography variant="body2" color="text.secondary">
                        {playlist.trackCount} треков
                      </Typography>
                    </CardContent>
                  </Card>
                </Grid>
              ))
            )}
          </Grid>

          {/* Блок избранных плейлистов */}
          <Box sx={{ mt: 6 }}>
            <Typography variant="h4" component="h2" gutterBottom sx={{ mb: 3 }}>
              Избранные плейлисты
            </Typography>
            <Grid container spacing={3}>
              {favoritePlaylists.length === 0 ? (
                <Grid item xs={12}>
                  <Typography variant="h6" textAlign="center" color="text.secondary">
                    У вас пока нет избранных плейлистов
                  </Typography>
                </Grid>
              ) : (
                favoritePlaylists.map((playlist) => (
                  <Grid item xs={12} sm={6} md={3} key={playlist.id}>
                    <Card
                      sx={{
                        height: '100%',
                        display: 'flex',
                        flexDirection: 'column',
                        cursor: 'pointer',
                        border: '2px solid #a78bfa',
                        boxShadow: '0 0 0 2px #a78bfa',
                        '&:hover': {
                          transform: 'scale(1.02)',
                          transition: 'transform 0.2s ease-in-out',
                          boxShadow: '0 0 0 4px #a78bfa',
                        },
                      }}
                      onClick={() => handlePlaylistClick(playlist.id)}
                    >
                      {playlist.coverImage ? (
                        <CardMedia
                          component="img"
                          sx={{
                            height: 250,
                            width: '100%',
                            objectFit: 'cover',
                            aspectRatio: '1/1'
                          }}
                          image={playlist.coverImage}
                          alt={playlist.name}
                        />
                      ) : (
                        <Box
                          sx={{
                            height: 250,
                            width: '100%',
                            display: 'flex',
                            alignItems: 'center',
                            justifyContent: 'center',
                            bgcolor: 'primary.dark',
                            borderRadius: 2,
                          }}
                        >
                          <ImageIcon sx={{ fontSize: 64, color: 'primary.contrastText', opacity: 0.3 }} />
                        </Box>
                      )}
                      <CardContent>
                        <Typography gutterBottom variant="h6" component="div" noWrap>
                          {playlist.name}
                        </Typography>
                        <Typography variant="body2" color="#a78bfa">
                          ★ Избранное
                        </Typography>
                        <Typography variant="body2" color="text.secondary">
                          {playlist.trackCount} треков
                        </Typography>
                      </CardContent>
                    </Card>
                  </Grid>
                ))
              )}
            </Grid>
          </Box>
        </>
      )}
      <PlaylistDialog
        open={createDialogOpen}
        onClose={() => setCreateDialogOpen(false)}
        onSubmit={handleCreatePlaylist}
        title="Создать новый плейлист"
      />
    </Container>
  );
};

export default Library; 