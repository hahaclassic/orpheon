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
  IconButton,
} from '@mui/material';
import AddIcon from '@mui/icons-material/Add';
import ImageIcon from '@mui/icons-material/Image';
import FavoriteIcon from '@mui/icons-material/Favorite';
import FavoriteBorderIcon from '@mui/icons-material/FavoriteBorder';
import { apiService } from '../../../presentation/services/api';
import { useApi } from '../../../presentation/hooks/useApi';
import PlaylistDialog from '../../../presentation/components/PlaylistDialog';
import type { Track } from '../../../presentation/types';
import { useAuthContext } from '../../../presentation/contexts/AuthContext';

interface Playlist {
  id: number;
  name: string;
  coverImage?: string;
  trackCount: number;
  tracks: Track[];
  is_favorite: boolean;
  rating: number;
}

const Library = () => {
  const [myPlaylists, setMyPlaylists] = useState<Playlist[]>([]);
  const [favoritePlaylists, setFavoritePlaylists] = useState<Playlist[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [createDialogOpen, setCreateDialogOpen] = useState(false);
  const [updatingFavorite, setUpdatingFavorite] = useState(false);
  const navigate = useNavigate();
  const { createPlaylist } = useApi();
  const { user } = useAuthContext();

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
          return playlist; // Возвращаем плейлист без обложки
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
        fetchPlaylistCovers(myPlaylistsData ?? []),
        fetchPlaylistCovers(favoritePlaylistsData ?? []),
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
    if (!user) {
      navigate('/login', { 
        state: { 
          from: '/library',
          message: 'Чтобы получить доступ к своей библиотеке, необходимо войти'
        }
      });
      return;
    }
    fetchData();
  }, [user, navigate]);

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
    ...myPlaylists,
    ...favoritePlaylists,
  ];
  // Убираем дубли по id (если плейлист есть и в моих, и в избранных)
  const uniquePlaylists = allPlaylists.filter(
    (playlist, idx, arr) => arr.findIndex((p) => p.id === playlist.id) === idx
  );

  const handleFavoriteClick = async (e: React.MouseEvent, playlistId: number) => {
    e.stopPropagation();
    if (updatingFavorite) return;

    try {
      setUpdatingFavorite(true);
      const playlist = uniquePlaylists.find(p => p.id === playlistId);
      if (!playlist) return;

      if (playlist.is_favorite) {
        await apiService.delete(`/me/favorites/${playlistId}`);
        setMyPlaylists(prev => prev.map(p => 
          p.id === playlistId ? { ...p, is_favorite: false, rating: p.rating - 1 } : p
        ));
        setFavoritePlaylists(prev => prev.filter(p => p.id !== playlistId));
      } else {
        await apiService.post(`/me/favorites/${playlistId}`);
        setMyPlaylists(prev => prev.map(p => 
          p.id === playlistId ? { ...p, is_favorite: true, rating: p.rating + 1 } : p
        ));
        setFavoritePlaylists(prev => [...prev, { ...playlist, is_favorite: true, rating: playlist.rating + 1 }]);
      }
    } catch (err) {
      console.error("Error updating favorite status:", err);
      setError("Не удалось обновить статус избранного");
    } finally {
      setUpdatingFavorite(false);
    }
  };

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
                      '&:hover': {
                        transform: 'scale(1.02)',
                        transition: 'transform 0.2s ease-in-out',
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
                      <Box sx={{ display: 'flex', alignItems: 'center', gap: 0.5 }}>
                        <IconButton 
                          onClick={(e) => handleFavoriteClick(e, playlist.id)}
                          size="small"
                          sx={{ 
                            color: 'white',
                            '&:hover': { color: 'white' }
                          }}
                        >
                          {playlist.is_favorite ? (
                            <FavoriteIcon sx={{ fontSize: 20 }} />
                          ) : (
                            <FavoriteBorderIcon sx={{ fontSize: 20 }} />
                          )}
                        </IconButton>
                        <Typography variant="body2" color="white">
                          {playlist.rating || 0}
                        </Typography>
                      </Box>
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
                        '&:hover': {
                          transform: 'scale(1.02)',
                          transition: 'transform 0.2s ease-in-out',
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
                        <Box sx={{ display: 'flex', alignItems: 'center', gap: 0.5 }}>
                          <IconButton 
                            onClick={(e) => handleFavoriteClick(e, playlist.id)}
                            size="small"
                            sx={{ 
                              color: 'white',
                              '&:hover': { color: 'white' }
                            }}
                          >
                            {playlist.is_favorite ? (
                              <FavoriteIcon sx={{ fontSize: 20 }} />
                            ) : (
                              <FavoriteBorderIcon sx={{ fontSize: 20 }} />
                            )}
                          </IconButton>
                          <Typography variant="body2" color="white">
                            {playlist.rating || 0}
                          </Typography>
                        </Box>
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