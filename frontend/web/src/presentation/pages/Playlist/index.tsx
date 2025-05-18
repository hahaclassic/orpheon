import { useState, useEffect, useRef } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import {
  Box,
  Typography,
  Paper,
  IconButton,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  CircularProgress,
  Alert,
  Button,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  TextField,
  Tooltip,
  Menu,
  MenuItem,
  Switch,
  Divider,
} from '@mui/material';
import EditIcon from '@mui/icons-material/Edit';
import FavoriteIcon from '@mui/icons-material/Favorite';
import FavoriteBorderIcon from '@mui/icons-material/FavoriteBorder';
import MoreVertIcon from '@mui/icons-material/MoreVert';
import PhotoCameraIcon from '@mui/icons-material/PhotoCamera';
import DeleteIcon from '@mui/icons-material/Delete';
import axios from 'axios';
import { useAuthContext } from '../../contexts/AuthContext';

interface Track {
  id: string;
  title: string;
  artists: string[];
  duration: number;
  albumCover: string;
  albumName: string;
}

interface Playlist {
  id: string;
  name: string;
  description?: string;
  is_private: boolean;
  coverImage?: string;
  createdAt: string;
  updatedAt: string;
  rating: number;
  isFavorite?: boolean;
  owner_id: string;
  owner_name?: string;
}

const formatDuration = (seconds: number) => {
  const minutes = Math.floor(seconds / 60);
  const remainingSeconds = seconds % 60;
  return `${minutes}:${remainingSeconds.toString().padStart(2, '0')}`;
};

const PlaylistPage = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const { user } = useAuthContext();
  const fileInputRef = useRef<HTMLInputElement>(null);
  const [playlist, setPlaylist] = useState<Playlist | null>(null);
  const [tracks, setTracks] = useState<Track[]>([]);
  const [loading, setLoading] = useState(true);
  const [tracksLoading, setTracksLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [updatingPrivacy, setUpdatingPrivacy] = useState(false);
  const [editDialogOpen, setEditDialogOpen] = useState(false);
  const [editForm, setEditForm] = useState({
    name: '',
    description: '',
  });
  const [updatingFavorite, setUpdatingFavorite] = useState(false);
  const [uploadingCover, setUploadingCover] = useState(false);
  const [menuAnchor, setMenuAnchor] = useState<null | HTMLElement>(null);

  const isOwner = user && playlist && user.id === playlist.owner_id;

  useEffect(() => {
    const fetchPlaylist = async () => {
      if (!id) return;
      
      try {
        setLoading(true);
        const response = await axios.get(`/playlists/${id}`, {
          baseURL: 'http://localhost:8080/api/v1',
          withCredentials: true,
        });
        setPlaylist(response.data);
        setEditForm({
          name: response.data.name,
          description: response.data.description || '',
        });

        // Загружаем обложку
        try {
          const coverResponse = await axios.get(`/playlists/${id}/cover`, {
            baseURL: 'http://localhost:8080/api/v1',
            withCredentials: true,
            responseType: 'blob',
          });
          
          // Создаем URL для blob
          const coverUrl = URL.createObjectURL(coverResponse.data);
          setPlaylist(prev => prev ? {
            ...prev,
            coverImage: coverUrl
          } : null);
        } catch (err) {
          // Если обложка не найдена, оставляем coverImage как undefined
          console.log('Обложка не найдена');
        }
      } catch (err) {
        setError('Не удалось загрузить плейлист');
      } finally {
        setLoading(false);
      }
    };

    const fetchTracks = async () => {
      if (!id) return;

      try {
        setTracksLoading(true);
        const response = await axios.get(`/playlists/${id}/tracks`, {
          baseURL: 'http://localhost:8080/api/v1',
          withCredentials: true,
        });
        setTracks(response.data || []);
      } catch (err) {
        setError('Не удалось загрузить треки');
        setTracks([]);
      } finally {
        setTracksLoading(false);
      }
    };

    fetchPlaylist();
    fetchTracks();

    // Очищаем URL при размонтировании компонента
    return () => {
      if (playlist?.coverImage?.startsWith('blob:')) {
        URL.revokeObjectURL(playlist.coverImage);
      }
    };
  }, [id]);

  const handlePrivacyChange = async (event: React.ChangeEvent<HTMLInputElement>) => {
    if (!playlist) return;

    try {
      setUpdatingPrivacy(true);
      await axios.patch(
        `/playlists/${id}/privacy`,
        { is_private: event.target.checked },
        {
          baseURL: 'http://localhost:8080/api/v1',
          withCredentials: true,
        }
      );
      setPlaylist(prev => prev ? { ...prev, is_private: event.target.checked } : null);
    } catch (err) {
      setError('Не удалось изменить настройки приватности');
    } finally {
      setUpdatingPrivacy(false);
    }
  };

  const handleEditSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!playlist) return;

    try {
      const response = await axios.put(
        `/playlists/${id}`,
        editForm,
        {
          baseURL: 'http://localhost:8080/api/v1',
          withCredentials: true,
        }
      );
      setPlaylist(response.data);
      setEditDialogOpen(false);
    } catch (err) {
      setError('Не удалось обновить информацию о плейлисте');
    }
  };

  const handleFavoriteClick = async () => {
    if (!playlist || updatingFavorite) return;

    try {
      setUpdatingFavorite(true);
      if (playlist.isFavorite) {
        await axios.delete(`/me/favorites/${id}`, {
          baseURL: 'http://localhost:8080/api/v1',
          withCredentials: true,
        });
        setPlaylist(prev => prev ? {
          ...prev,
          isFavorite: false,
          rating: prev.rating - 1
        } : null);
      } else {
        await axios.post(`/me/favorites/${id}`, {}, {
          baseURL: 'http://localhost:8080/api/v1',
          withCredentials: true,
        });
        setPlaylist(prev => prev ? {
          ...prev,
          isFavorite: true,
          rating: prev.rating + 1
        } : null);
      }
    } catch (err) {
      setError('Не удалось обновить статус избранного');
    } finally {
      setUpdatingFavorite(false);
    }
  };

  const handleCoverClick = () => {
    fileInputRef.current?.click();
  };

  const handleCoverChange = async (event: React.ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0];
    if (!file || !id) return;

    try {
      setUploadingCover(true);
      const formData = new FormData();
      formData.append('cover', file);

      await axios.post(`/playlists/${id}/cover`, formData, {
        baseURL: 'http://localhost:8080/api/v1',
        withCredentials: true,
        headers: {
          'Content-Type': 'multipart/form-data',
        },
      });

      // Загружаем новую обложку
      const coverResponse = await axios.get(`/playlists/${id}/cover`, {
        baseURL: 'http://localhost:8080/api/v1',
        withCredentials: true,
        responseType: 'blob',
      });

      // Освобождаем старый URL если он был
      if (playlist?.coverImage?.startsWith('blob:')) {
        URL.revokeObjectURL(playlist.coverImage);
      }

      // Создаем новый URL для blob
      const coverUrl = URL.createObjectURL(coverResponse.data);
      setPlaylist(prev => prev ? {
        ...prev,
        coverImage: coverUrl
      } : null);
    } catch (err) {
      setError('Не удалось загрузить обложку');
    } finally {
      setUploadingCover(false);
      event.target.value = '';
    }
  };

  const handleDeleteCover = async () => {
    if (!id) return;

    try {
      setUploadingCover(true);
      await axios.delete(`/playlists/${id}/cover`, {
        baseURL: 'http://localhost:8080/api/v1',
        withCredentials: true,
      });

      // Освобождаем URL если он был
      if (playlist?.coverImage?.startsWith('blob:')) {
        URL.revokeObjectURL(playlist.coverImage);
      }

      setPlaylist(prev => prev ? {
        ...prev,
        coverImage: undefined
      } : null);
    } catch (err) {
      setError('Не удалось удалить обложку');
    } finally {
      setUploadingCover(false);
    }
  };

  // Меню три точки
  const handleMenuOpen = (e: React.MouseEvent<HTMLElement>) => setMenuAnchor(e.currentTarget);
  const handleMenuClose = () => setMenuAnchor(null);

  if (loading) {
    return (
      <Box sx={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: '100vh' }}>
        <CircularProgress />
      </Box>
    );
  }

  if (error || !playlist) {
    return (
      <Box sx={{ p: 4 }}>
        <Alert severity="error">{error || 'Плейлист не найден'}</Alert>
      </Box>
    );
  }

  return (
    <Box sx={{ display: 'flex', flexDirection: 'column', minHeight: '100vh', width: '100%' }}>
      {/* Header Section */}
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
        <Box sx={{ display: 'flex', alignItems: 'center', gap: 3, pl: { xs: 2, sm: 4, md: 6 } }}>
          <Box
            sx={{
              width: 200,
              height: 200,
              bgcolor: 'primary.main',
              borderRadius: 1,
              overflow: 'hidden',
              position: 'relative',
              '&:hover .cover-actions': {
                opacity: 1,
              },
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
            {isOwner && (
              <Box
                className="cover-actions"
                sx={{
                  position: 'absolute',
                  top: 0,
                  left: 0,
                  right: 0,
                  bottom: 0,
                  bgcolor: 'rgba(0, 0, 0, 0.5)',
                  display: 'flex',
                  alignItems: 'center',
                  justifyContent: 'center',
                  gap: 1,
                  opacity: 0,
                  transition: 'opacity 0.2s',
                }}
              >
                <Tooltip title="Изменить обложку">
                  <IconButton
                    onClick={handleCoverClick}
                    disabled={uploadingCover}
                    sx={{ color: 'white' }}
                  >
                    <PhotoCameraIcon />
                  </IconButton>
                </Tooltip>
                {playlist.coverImage && (
                  <Tooltip title="Удалить обложку">
                    <IconButton
                      onClick={handleDeleteCover}
                      disabled={uploadingCover}
                      sx={{ color: 'white' }}
                    >
                      <DeleteIcon />
                    </IconButton>
                  </Tooltip>
                )}
              </Box>
            )}
            <input
              type="file"
              ref={fileInputRef}
              onChange={handleCoverChange}
              accept="image/*"
              style={{ display: 'none' }}
            />
          </Box>
          <Box sx={{ flex: 1 }}>
            <Box sx={{ display: 'flex', alignItems: 'center', gap: 2, mb: 1 }}>
              <Typography variant="h4" fontWeight={700}>
                {playlist.name}
              </Typography>
              <IconButton
                onClick={() => setEditDialogOpen(true)}
                sx={{ color: 'text.secondary' }}
              >
                <EditIcon />
              </IconButton>
              <Tooltip title={playlist?.isFavorite ? "Удалить из избранного" : "Добавить в избранное"}>
                <IconButton
                  onClick={handleFavoriteClick}
                  disabled={updatingFavorite}
                  sx={{ color: playlist?.isFavorite ? 'background.default' : 'text.secondary' }}
                >
                  {playlist?.isFavorite ? <FavoriteIcon /> : <FavoriteBorderIcon />}
                </IconButton>
              </Tooltip>
              <Typography variant="body2" color="text.secondary">
                {playlist.rating}
              </Typography>
            </Box>
            {playlist.description && (
              <Typography variant="body1" color="text.secondary" sx={{ mb: 2 }}>
                {playlist.description}
              </Typography>
            )}
            <Typography variant="body2" color="text.secondary" gutterBottom>
              Создан: {new Date(playlist.createdAt).toLocaleDateString()}
            </Typography>
            <Typography variant="body2" color="text.secondary" gutterBottom>
              Обновлен: {new Date(playlist.updatedAt).toLocaleDateString()}
            </Typography>
          </Box>
        </Box>
      </Box>

      {/* Tracks Section */}
      <Box sx={{ flex: 1, bgcolor: 'background.default', px: { xs: 2, sm: 6, md: 10 }, py: 4 }}>
        <Typography variant="h5" fontWeight={700} gutterBottom>
          Треки {!tracksLoading && `(${tracks.length})`}
        </Typography>
        {tracksLoading ? (
          <Box sx={{ display: 'flex', justifyContent: 'center', p: 4 }}>
            <CircularProgress />
          </Box>
        ) : tracks.length === 0 ? (
          <Typography color="text.secondary" sx={{ p: 4, textAlign: 'center' }}>
            В плейлисте пока нет треков
          </Typography>
        ) : (
          <TableContainer component={Paper}>
            <Table>
              <TableHead>
                <TableRow>
                  <TableCell width={50}>#</TableCell>
                  <TableCell>Название</TableCell>
                  <TableCell>Альбом</TableCell>
                  <TableCell width={100} align="right">Длительность</TableCell>
                </TableRow>
              </TableHead>
              <TableBody>
                {tracks.map((track, index) => (
                  <TableRow
                    key={track.id}
                    hover
                    sx={{ cursor: 'pointer' }}
                    onClick={() => navigate(`/tracks/${track.id}`)}
                  >
                    <TableCell>{index + 1}</TableCell>
                    <TableCell>
                      <Box sx={{ display: 'flex', alignItems: 'center', gap: 2 }}>
                        <Box
                          sx={{
                            width: 40,
                            height: 40,
                            borderRadius: 1,
                            overflow: 'hidden',
                          }}
                        >
                          <Box
                            component="img"
                            src={track.albumCover}
                            alt={track.albumName}
                            sx={{
                              width: '100%',
                              height: '100%',
                              objectFit: 'cover',
                            }}
                          />
                        </Box>
                        <Box>
                          <Typography variant="body1">{track.title}</Typography>
                          <Typography variant="body2" color="text.secondary">
                            {track.artists.join(', ')}
                          </Typography>
                        </Box>
                      </Box>
                    </TableCell>
                    <TableCell>{track.albumName}</TableCell>
                    <TableCell align="right">{formatDuration(track.duration)}</TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </TableContainer>
        )}
      </Box>

      {/* Edit Dialog */}
      <Dialog open={editDialogOpen} onClose={() => setEditDialogOpen(false)} maxWidth="sm" fullWidth>
        <form onSubmit={handleEditSubmit}>
          <DialogTitle>Редактировать плейлист</DialogTitle>
          <DialogContent sx={{ display: 'flex', flexDirection: 'column', gap: 2, pt: 2 }}>
            <TextField
              label="Название"
              value={editForm.name}
              onChange={(e) => setEditForm(prev => ({ ...prev, name: e.target.value }))}
              required
              fullWidth
            />
            <TextField
              label="Описание"
              value={editForm.description}
              onChange={(e) => setEditForm(prev => ({ ...prev, description: e.target.value }))}
              multiline
              rows={4}
              fullWidth
            />
          </DialogContent>
          <DialogActions>
            <Button onClick={() => setEditDialogOpen(false)}>Отмена</Button>
            <Button type="submit" variant="contained">Сохранить</Button>
          </DialogActions>
        </form>
      </Dialog>
    </Box>
  );
};

export default PlaylistPage; 