import { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import {
  Box,
  Typography,
  Paper,
  Avatar,
  Grid,
  Container,
  Tabs,
  Tab,
  Card,
  CardContent,
  CardMedia,
  CircularProgress,
  Alert,
} from '@mui/material';
import ImageIcon from '@mui/icons-material/Image';
import { api } from '../../services/api';

interface User {
  id: string;
  name: string;
  registration_date: string;
  birth_date: string;
  access_lvl: number;
}

interface Playlist {
  id: string;
  name: string;
  coverImage?: string;
  trackCount: number;
  is_favorite: boolean;
}

const UserProfile = () => {
  const { id } = useParams<{ id: string }>();
  const [activeTab, setActiveTab] = useState(0);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const [user, setUser] = useState<User | null>(null);
  const [userPlaylists, setUserPlaylists] = useState<Playlist[]>([]);
  const [favoritePlaylists, setFavoritePlaylists] = useState<Playlist[]>([]);

  const navigate = useNavigate();

  useEffect(() => {
    const fetchData = async () => {
      try {
        setLoading(true);
        const [userData, userPlaylistsData, favoritePlaylistsData] = await Promise.all([
          api.getUser(id!),
          api.getUserPlaylists(id!),
          api.getUserFavorites(id!),
        ]);

        // Загрузка обложек для плейлистов
        const loadPlaylistCovers = async (playlists: Playlist[]) => {
          return Promise.all(
            playlists.map(async (playlist) => {
              try {
                const response = await apiService.get(`/playlists/${playlist.id}/cover`, {
                  responseType: 'blob'
                });
                const imageUrl = URL.createObjectURL(response);
                return { ...playlist, coverImage: imageUrl };
              } catch (error) {
                console.error(`Error loading cover for playlist ${playlist.id}:`, error);
                return playlist;
              }
            })
          );
        };

        const [userPlaylistsWithCovers, favoritePlaylistsWithCovers] = await Promise.all([
          loadPlaylistCovers(userPlaylistsData),
          loadPlaylistCovers(favoritePlaylistsData)
        ]);

        setUser(userData);
        setUserPlaylists(userPlaylistsWithCovers);
        setFavoritePlaylists(favoritePlaylistsWithCovers);
      } catch (err) {
        console.error('Error fetching user data:', err);
        setError('Ошибка при загрузке данных пользователя');
      } finally {
        setLoading(false);
      }
    };

    if (id) {
      fetchData();
    }
  }, [id]);

  const handlePlaylistClick = (playlistId: string) => {
    navigate(`/playlists/${playlistId}`);
  };

  if (loading) {
    return (
      <Container maxWidth="lg" sx={{ py: 4 }}>
        <Box sx={{ display: 'flex', justifyContent: 'center', p: 3 }}>
          <CircularProgress />
        </Box>
      </Container>
    );
  }

  if (error || !user) {
    return (
      <Container maxWidth="lg" sx={{ py: 4 }}>
        <Alert severity="error" sx={{ mb: 2 }}>{error || 'Пользователь не найден'}</Alert>
      </Container>
    );
  }

  return (
    <Container 
      maxWidth="lg" 
      sx={{ 
        py: 4,
        height: 'calc(100vh - 90px)',
        display: 'flex',
        flexDirection: 'column'
      }}
    >
      {/* User Info Section */}
      <Box sx={{ mb: 4 }}>
        <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 3 }}>
          <Typography variant="h4" component="h1">
            Профиль пользователя
          </Typography>
        </Box>

        <Grid container spacing={3}>
          <Grid item xs={12} md={6}>
            <Card>
              <CardContent>
                <Box sx={{ display: 'flex', alignItems: 'center', mb: 3 }}>
                  <Avatar
                    sx={{
                      width: 100,
                      height: 100,
                      bgcolor: 'primary.main',
                      fontSize: '2rem',
                      mr: 2
                    }}
                  >
                    {user.name.charAt(0).toUpperCase()}
                  </Avatar>
                  <Box>
                    <Typography variant="h5" gutterBottom>
                      {user.name}
                    </Typography>
                    <Typography variant="body2" color="text.secondary">
                      {user.access_lvl === 1 ? 'Администратор' : 'Пользователь'}
                    </Typography>
                  </Box>
                </Box>
                <Box sx={{ mb: 2 }}>
                  <Typography variant="subtitle2" color="text.secondary">
                    Дата регистрации
                  </Typography>
                  <Typography variant="body1">
                    {new Date(user.registration_date).toLocaleDateString()}
                  </Typography>
                </Box>
                <Box sx={{ mb: 2 }}>
                  <Typography variant="subtitle2" color="text.secondary">
                    Дата рождения
                  </Typography>
                  <Typography variant="body1">
                    {user.birth_date === '0001-01-01T00:00:00Z' 
                      ? 'Не указана' 
                      : new Date(user.birth_date).toLocaleDateString()}
                  </Typography>
                </Box>
              </CardContent>
            </Card>
          </Grid>
        </Grid>
      </Box>

      {/* Tabs */}
      <Box sx={{ borderBottom: 1, borderColor: 'divider', bgcolor: 'background.paper' }}>
        <Container maxWidth="xl">
          <Tabs 
            value={activeTab} 
            onChange={(_, newValue) => setActiveTab(newValue)}
            sx={{ minHeight: 64 }}
          >
            <Tab label="Плейлисты пользователя" />
            <Tab label="Избранные плейлисты" />
          </Tabs>
        </Container>
      </Box>

      {/* Content */}
      <Box sx={{ 
        flex: 1, 
        py: 4,
        overflow: 'auto'
      }}>
        <Container maxWidth="xl">
          {activeTab === 0 && (
            <Box>
              <Typography variant="h5" gutterBottom>Плейлисты пользователя</Typography>
              <Grid container spacing={3}>
                {userPlaylists.length === 0 ? (
                  <Grid item xs={12}>
                    <Typography variant="h6" textAlign="center" color="text.secondary">
                      У пользователя пока нет плейлистов
                    </Typography>
                  </Grid>
                ) : (
                  userPlaylists.map((playlist) => (
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
          )}

          {activeTab === 1 && (
            <Box>
              <Typography variant="h5" gutterBottom>Избранные плейлисты</Typography>
              <Grid container spacing={3}>
                {favoritePlaylists.length === 0 ? (
                  <Grid item xs={12}>
                    <Typography variant="h6" textAlign="center" color="text.secondary">
                      У пользователя пока нет избранных плейлистов
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
          )}
        </Container>
      </Box>
    </Container>
  );
};

export default UserProfile; 