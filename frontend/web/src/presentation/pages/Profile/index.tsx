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
} from '@mui/material';
import LockIcon from '@mui/icons-material/Lock';
import axios from 'axios';

interface UserProfile {
  username: string;
  name?: string;
  email: string;
  avatar?: string;
  birthDate?: string;
  createdAt?: string;
}

interface Playlist {
  id: string;
  name: string;
  isPublic: boolean;
  coverImage?: string;
}

const Profile = () => {
  const [passwordDialogOpen, setPasswordDialogOpen] = useState(false);
  const [passwordData, setPasswordData] = useState({
    old: '',
    new: '',
    confirm: '',
  });
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState<string | null>(null);

  const [profile, setProfile] = useState<UserProfile | null>(null);
  const [profileLoading, setProfileLoading] = useState(true);
  const [playlists, setPlaylists] = useState<Playlist[]>([]);
  const [playlistsLoading, setPlaylistsLoading] = useState(true);
  const [favorites, setFavorites] = useState<Playlist[]>([]);
  const [favoritesLoading, setFavoritesLoading] = useState(true);

  const navigate = useNavigate();

  useEffect(() => {
    setProfileLoading(true);
    axios.get('/me', {
      baseURL: 'http://localhost:8080/api/v1',
      withCredentials: true,
    })
      .then(res => setProfile(res.data))
      .catch(() => setProfile(null))
      .finally(() => setProfileLoading(false));
  }, []);

  useEffect(() => {
    setPlaylistsLoading(true);
    axios.get('/me/playlists', {
      baseURL: 'http://localhost:8080/api/v1',
      withCredentials: true,
    })
      .then(res => setPlaylists(res.data))
      .catch(() => setPlaylists([]))
      .finally(() => setPlaylistsLoading(false));
  }, []);

  useEffect(() => {
    setFavoritesLoading(true);
    axios.get('/me/favorites', {
      baseURL: 'http://localhost:8080/api/v1',
      withCredentials: true,
    })
      .then(res => setFavorites(res.data))
      .catch(() => setFavorites([]))
      .finally(() => setFavoritesLoading(false));
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
    setLoading(true);
    try {
      await axios.post(
        'http://localhost:8080/api/v1/auth/password/update',
        { old: passwordData.old, new: passwordData.new },
        {
          withCredentials: true,
          headers: {
            'Content-Type': 'application/json',
          },
        }
      );
      setSuccess('Пароль успешно изменён!');
      setPasswordData({ old: '', new: '', confirm: '' });
      setPasswordDialogOpen(false);
    } catch (err: any) {
      setError(err?.response?.data?.message || 'Ошибка при смене пароля');
    } finally {
      setLoading(false);
    }
  };

  const handlePlaylistClick = (playlistId: string) => {
    navigate(`/playlists/${playlistId}`);
  };

  return (
    <Box sx={{ width: '100%', height: '100%', minHeight: '100vh', bgcolor: 'background.default' }}>
      {/* Header */}
      <Box
        sx={{
          bgcolor: 'primary.dark',
          py: 6,
          borderBottom: '1px solid',
          borderColor: 'divider',
          zIndex: 10,
          width: '100%',
        }}
      >
        <Box sx={{ display: 'flex', alignItems: 'center', gap: 3, px: 6 }}>
          <Avatar
            sx={{
              width: 120,
              height: 120,
              fontSize: 48,
              border: '4px solid',
              borderColor: 'primary.main',
              bgcolor: 'primary.main',
            }}
            src={profile?.avatar}
          >
            {profile?.username?.[0]}
          </Avatar>
          <Box>
            <Typography variant="h4" fontWeight={700} gutterBottom>
              {profileLoading ? 'Загрузка...' : profile?.name || profile?.username || '—'}
            </Typography>
            <Typography variant="subtitle1" color="text.secondary">
              {profileLoading ? '' : profile?.email || ''}
            </Typography>
            {profile?.birthDate && (
              <Typography variant="body2" color="text.secondary">
                Дата рождения: {new Date(profile.birthDate).toLocaleDateString()}
              </Typography>
            )}
            {profile?.createdAt && (
              <Typography variant="body2" color="text.secondary">
                Дата регистрации: {new Date(profile.createdAt).toLocaleDateString()}
              </Typography>
            )}
          </Box>
        </Box>
      </Box>

      {/* Контент */}
      <Box sx={{ width: '100%', py: 4 }}>
        <Typography variant="h5" fontWeight={700} gutterBottom sx={{ px: 6 }}>Мои плейлисты</Typography>
        <Grid container spacing={3} sx={{ width: '100%', m: 0, px: 6 }}>
          {playlistsLoading ? (
            <Typography color="text.secondary" sx={{ ml: 2 }}>Загрузка...</Typography>
          ) : playlists?.length === 0 ? (
            <Typography color="text.secondary" sx={{ ml: 2 }}>Нет плейлистов</Typography>
          ) : (
            (playlists || []).map((playlist) => (
              <Grid item xs={12} sm={6} md={4} key={playlist.id}>
                <Paper 
                  sx={{ 
                    p: 3,
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
                  <Box
                    sx={{
                      width: '100%',
                      aspectRatio: '1',
                      mb: 2,
                      bgcolor: 'primary.dark',
                      borderRadius: 1,
                      overflow: 'hidden',
                    }}
                  >
                    {playlist.coverImage ? (
                      <Box
                        component="img"
                        src={playlist.coverImage}
                        alt={playlist.name}
                        sx={{
                          width: '100%',
                          height: '100%',
                          objectFit: 'cover',
                        }}
                      />
                    ) : (
                      <Box
                        sx={{
                          width: '100%',
                          height: '100%',
                          bgcolor: 'primary.main',
                        }}
                      />
                    )}
                  </Box>
                  <Typography variant="h6" noWrap>{playlist.name}</Typography>
                  {!playlist.isPublic && (
                    <Typography variant="caption" color="text.secondary">
                      Приватный
                    </Typography>
                  )}
                </Paper>
              </Grid>
            ))
          )}
        </Grid>

        <Typography variant="h5" fontWeight={700} gutterBottom sx={{ px: 6, mt: 4 }}>Избранное</Typography>
        <Grid container spacing={3} sx={{ width: '100%', m: 0, px: 6 }}>
          {favoritesLoading ? (
            <Typography color="text.secondary" sx={{ ml: 2 }}>Загрузка...</Typography>
          ) : favorites?.length === 0 ? (
            <Typography color="text.secondary" sx={{ ml: 2 }}>Нет избранных плейлистов</Typography>
          ) : (
            (favorites || []).map((playlist) => (
              <Grid item xs={12} sm={6} md={4} key={playlist.id}>
                <Paper 
                  sx={{ 
                    p: 3,
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
                  <Box
                    sx={{
                      width: '100%',
                      aspectRatio: '1',
                      mb: 2,
                      bgcolor: 'primary.dark',
                      borderRadius: 1,
                      overflow: 'hidden',
                    }}
                  >
                    {playlist.coverImage ? (
                      <Box
                        component="img"
                        src={playlist.coverImage}
                        alt={playlist.name}
                        sx={{
                          width: '100%',
                          height: '100%',
                          objectFit: 'cover',
                        }}
                      />
                    ) : (
                      <Box
                        sx={{
                          width: '100%',
                          height: '100%',
                          bgcolor: 'primary.main',
                        }}
                      />
                    )}
                  </Box>
                  <Typography variant="h6" noWrap>{playlist.name}</Typography>
                </Paper>
              </Grid>
            ))
          )}
        </Grid>

        <Box sx={{ px: 6, mt: 4 }}>
          <Typography variant="h6" gutterBottom>
            Настройки аккаунта
          </Typography>
          <Button
            variant="outlined"
            startIcon={<LockIcon />}
            onClick={() => setPasswordDialogOpen(true)}
          >
            Сменить пароль
          </Button>
        </Box>
      </Box>

      {/* Диалог смены пароля — оставить как есть */}
      <Dialog open={passwordDialogOpen} onClose={() => setPasswordDialogOpen(false)} maxWidth="xs" fullWidth>
        <form onSubmit={handlePasswordSubmit}>
          <DialogTitle>Смена пароля</DialogTitle>
          <DialogContent sx={{ display: 'flex', flexDirection: 'column', gap: 2, pt: 2 }}>
            {error && <Alert severity="error">{error}</Alert>}
            {success && <Alert severity="success">{success}</Alert>}
            <TextField
              label="Текущий пароль"
              name="old"
              type="password"
              value={passwordData.old}
              onChange={handlePasswordChange}
              required
              fullWidth
            />
            <TextField
              label="Новый пароль"
              name="new"
              type="password"
              value={passwordData.new}
              onChange={handlePasswordChange}
              required
              fullWidth
            />
            <TextField
              label="Подтвердите новый пароль"
              name="confirm"
              type="password"
              value={passwordData.confirm}
              onChange={handlePasswordChange}
              required
              fullWidth
            />
          </DialogContent>
          <DialogActions>
            <Button onClick={() => setPasswordDialogOpen(false)} disabled={loading}>
              Отмена
            </Button>
            <Button type="submit" variant="contained" disabled={loading}>
              Сменить
            </Button>
          </DialogActions>
        </form>
      </Dialog>
    </Box>
  );
};

export default Profile; 