import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import {
  Box,
  Typography,
  Paper,
  Avatar,
  Button,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  TextField,
  Alert,
  Grid,
  Container,
  Tabs,
  Tab,
  Card,
  CardContent,
  CardMedia,
  IconButton,
  CircularProgress,
} from '@mui/material';
import LockIcon from '@mui/icons-material/Lock';
import EditIcon from '@mui/icons-material/Edit';
import ImageIcon from '@mui/icons-material/Image';
import axios from 'axios';
import { apiService } from '../../services/api';
import { useAuthContext } from '../../contexts/AuthContext';

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

const Profile = () => {
  const [activeTab, setActiveTab] = useState(0);
  const [passwordDialogOpen, setPasswordDialogOpen] = useState(false);
  const [passwordData, setPasswordData] = useState({
    old: '',
    new: '',
    confirm: '',
  });
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState<string | null>(null);

  const [user, setUser] = useState<User | null>(null);
  const [myPlaylists, setMyPlaylists] = useState<Playlist[]>([]);
  const [favoritePlaylists, setFavoritePlaylists] = useState<Playlist[]>([]);
  const [editDialogOpen, setEditDialogOpen] = useState(false);
  const [editForm, setEditForm] = useState({
    name: '',
    birth_date: '',
  });

  const navigate = useNavigate();
  const { user: authUser } = useAuthContext();

  useEffect(() => {
    const fetchData = async () => {
      try {
        setLoading(true);
        const [userData, myPlaylistsData, favoritePlaylistsData] = await Promise.all([
          apiService.get('/me'),
          apiService.get('/me/playlists'),
          apiService.get('/me/favorites'),
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

        const [myPlaylistsWithCovers, favoritePlaylistsWithCovers] = await Promise.all([
          loadPlaylistCovers(myPlaylistsData),
          loadPlaylistCovers(favoritePlaylistsData)
        ]);

        setUser(userData);
        setMyPlaylists(myPlaylistsWithCovers);
        setFavoritePlaylists(favoritePlaylistsWithCovers);
      } catch (err) {
        console.error('Error fetching profile data:', err);
        setError('Ошибка при загрузке данных профиля');
      } finally {
        setLoading(false);
      }
    };

    fetchData();
  }, []);

  const handlePasswordChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    setPasswordData((prev) => ({ ...prev, [name]: value }));
  };

  const handlePasswordSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);
    setSuccess(null);
    if (passwordData.new !== passwordData.confirm) {
      setError('Новые пароли не совпадают');
      return;
    }
    try {
      await apiService.put('/me/password', {
        old: passwordData.old,
        new: passwordData.new
      });
      setSuccess('Пароль успешно изменён!');
      setPasswordData({ old: '', new: '', confirm: '' });
      setPasswordDialogOpen(false);
    } catch (err: any) {
      setError(err?.response?.data?.message || 'Ошибка при смене пароля');
    }
  };

  const handleEditClick = () => {
    if (user) {
      setEditForm({
        name: user.name,
        birth_date: user.birth_date === '0001-01-01T00:00:00Z' ? '' : user.birth_date.split('T')[0],
      });
      setEditDialogOpen(true);
    }
  };

  const handleEditSubmit = async () => {
    try {
      const updatedUser = await apiService.put('/me', {
        name: editForm.name,
        birth_date: editForm.birth_date ? new Date(editForm.birth_date).toISOString() : '0001-01-01T00:00:00Z'
      });
      setUser(updatedUser);
      setEditDialogOpen(false);
    } catch (err) {
      console.error('Error updating profile:', err);
      setError('Ошибка при обновлении профиля');
    }
  };

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
        <Alert severity="error" sx={{ mb: 2 }}>{error || 'Профиль не найден'}</Alert>
      </Container>
    );
  }

  return (
    <Container 
      maxWidth="lg" 
      sx={{ 
        py: 4,
        height: 'calc(100vh - 90px)', // Высота экрана минус высота плеера
        display: 'flex',
        flexDirection: 'column'
      }}
    >
      {/* User Info Section */}
      <Box sx={{ mb: 4 }}>
        <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 3 }}>
          <Typography variant="h4" component="h1">
            Профиль
          </Typography>
          <Box>
            <Button
              variant="outlined"
              startIcon={<EditIcon />}
              onClick={handleEditClick}
              sx={{ mr: 2 }}
            >
              Редактировать
            </Button>
            <Button
              variant="outlined"
              onClick={() => setPasswordDialogOpen(true)}
            >
              Изменить пароль
            </Button>
          </Box>
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
            <Tab label="Мои плейлисты" />
            <Tab label="Избранные плейлисты" />
          </Tabs>
        </Container>
      </Box>

      {/* Content */}
      <Box sx={{ 
        flex: 1, 
        py: 4,
        overflow: 'auto' // Добавляем прокрутку при необходимости
      }}>
        <Container maxWidth="xl">
          {activeTab === 0 && (
            <Box>
              <Typography variant="h5" gutterBottom>Мои плейлисты</Typography>
              <Grid container spacing={3}>
                {myPlaylists.length === 0 ? (
                  <Grid item xs={12}>
                    <Typography variant="h6" textAlign="center" color="text.secondary">
                      У вас пока нет плейлистов
                    </Typography>
                  </Grid>
                ) : (
                  myPlaylists.map((playlist) => (
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

      {/* Password Change Dialog */}
      <Dialog open={passwordDialogOpen} onClose={() => setPasswordDialogOpen(false)}>
        <DialogTitle>Изменение пароля</DialogTitle>
        <form onSubmit={handlePasswordSubmit}>
          <DialogContent>
            {error && <Alert severity="error" sx={{ mb: 2 }}>{error}</Alert>}
            {success && <Alert severity="success" sx={{ mb: 2 }}>{success}</Alert>}
            <TextField
              autoFocus
              margin="dense"
              name="old"
              label="Текущий пароль"
              type="password"
              fullWidth
              value={passwordData.old}
              onChange={handlePasswordChange}
              required
            />
            <TextField
              margin="dense"
              name="new"
              label="Новый пароль"
              type="password"
              fullWidth
              value={passwordData.new}
              onChange={handlePasswordChange}
              required
            />
            <TextField
              margin="dense"
              name="confirm"
              label="Подтвердите новый пароль"
              type="password"
              fullWidth
              value={passwordData.confirm}
              onChange={handlePasswordChange}
              required
            />
          </DialogContent>
          <DialogActions>
            <Button onClick={() => setPasswordDialogOpen(false)}>Отмена</Button>
            <Button type="submit" variant="contained">
              Сохранить
            </Button>
          </DialogActions>
        </form>
      </Dialog>

      {/* Edit Dialog */}
      <Dialog open={editDialogOpen} onClose={() => setEditDialogOpen(false)}>
        <DialogTitle>Редактировать профиль</DialogTitle>
        <DialogContent>
          <Box sx={{ pt: 2 }}>
            <TextField
              fullWidth
              label="Имя пользователя"
              value={editForm.name}
              onChange={(e) => setEditForm(prev => ({ ...prev, name: e.target.value }))}
              sx={{ mb: 2 }}
            />
            <TextField
              fullWidth
              label="Дата рождения"
              type="date"
              value={editForm.birth_date}
              onChange={(e) => setEditForm(prev => ({ ...prev, birth_date: e.target.value }))}
              InputLabelProps={{ shrink: true }}
            />
          </Box>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setEditDialogOpen(false)}>Отмена</Button>
          <Button onClick={handleEditSubmit} variant="contained">
            Сохранить
          </Button>
        </DialogActions>
      </Dialog>
    </Container>
  );
};

export default Profile; 