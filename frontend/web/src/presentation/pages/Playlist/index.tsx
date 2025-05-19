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
  Container,
  Grid,
  Card,
  CardMedia,
  List,
  ListItem,
  ListItemText,
} from '@mui/material';
import EditIcon from '@mui/icons-material/Edit';
import FavoriteIcon from '@mui/icons-material/Favorite';
import FavoriteBorderIcon from '@mui/icons-material/FavoriteBorder';
import MoreVertIcon from '@mui/icons-material/MoreVert';
import PhotoCameraIcon from '@mui/icons-material/PhotoCamera';
import DeleteIcon from '@mui/icons-material/Delete';
import { ArrowBack, PlayArrow, Pause } from '@mui/icons-material';
import axios from 'axios';
import { useAuthContext } from '../../contexts/AuthContext';
import { usePlayerContext } from '../../contexts/PlayerContext';
import { apiService } from '../../services/api';

const API_URL = 'http://localhost:8080/api/v1';

interface Track {
  id: string;
  name: string;
  duration: number;
  track_number: number;
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
  const { currentTrack, isPlaying, setTrack, togglePlay } = usePlayerContext();

  const isOwner = user && playlist && user.id === playlist.owner_id;

  useEffect(() => {
    const fetchPlaylistData = async () => {
      if (!id) {
        setError('Playlist ID is missing');
        setLoading(false);
        return;
      }
      
      try {
        setLoading(true);
        const data = await apiService.get(`/playlists/${id}`);
        setPlaylist(data);
        
        // Get user information
        try {
          const userData = await apiService.get(`/users/${data.owner_id}`);
          setPlaylist(prev => prev ? { ...prev, owner_name: userData.name } : null);
        } catch (err) {
          console.error('Error fetching user data:', err);
        }

        // Get playlist cover
        try {
          const coverResponse = await apiService.get(`/playlists/${id}/cover`, {
            responseType: 'blob'
          });
          const coverUrl = URL.createObjectURL(coverResponse);
          setPlaylist(prev => prev ? { ...prev, coverImage: coverUrl } : null);
        } catch (err) {
          console.error('Error fetching playlist cover:', err);
        }
        
        // Get playlist tracks
        try {
          setTracksLoading(true);
          const tracksResponse = await apiService.get(`/playlists/${id}/tracks`);
          setTracks(tracksResponse);
        } catch (err) {
          console.error('Error fetching playlist tracks:', err);
          setTracks([]);
        } finally {
          setTracksLoading(false);
        }

        setError(null);
      } catch (error) {
        console.error('Error fetching playlist data:', error);
        setError('Failed to load playlist data');
        setPlaylist(null);
      } finally {
        setLoading(false);
      }
    };

    fetchPlaylistData();
  }, [id]);

  const handlePrivacyChange = async (event: React.ChangeEvent<HTMLInputElement>) => {
    if (!playlist) return;

    try {
      setUpdatingPrivacy(true);
      await apiService({
        method: 'patch',
        url: `/playlists/${id}/privacy`,
        data: { is_private: event.target.checked }
      });
      setPlaylist(prev => prev ? { ...prev, is_private: event.target.checked } : null);
      event.target.checked = !event.target.checked;
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
      const response = await apiService.put(`/playlists/${id}`, {
        name: editForm.name,
        description: editForm.description,
        is_private: playlist.is_private
      });
      setPlaylist(response);
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
        await apiService.delete(`/me/favorites/${id}`);
        setPlaylist(prev => prev ? {
          ...prev,
          isFavorite: false,
          rating: prev.rating - 1
        } : null);
      } else {
        await apiService.post(`/me/favorites/${id}`);
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

      await apiService.post(`/playlists/${id}/cover`, formData, {
        headers: {
          'Content-Type': 'multipart/form-data',
        },
      });

      // Reload the cover
      const coverResponse = await apiService.get(`/playlists/${id}/cover`, {
        responseType: 'blob'
      });
      const coverUrl = URL.createObjectURL(coverResponse);
      setPlaylist(prev => prev ? { ...prev, coverImage: coverUrl } : null);
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
      await apiService.delete(`/playlists/${id}/cover`);

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

  const handleTrackClick = (trackId: string) => {
    const track = tracks.find(t => t.id === trackId);
    if (track) {
      if (currentTrack?.id === trackId) {
        togglePlay();
      } else {
        setTrack(track, tracks);
      }
    }
  };

  if (loading) {
    return (
      <Container maxWidth="lg" sx={{ py: 4 }}>
        <Typography>Loading...</Typography>
      </Container>
    );
  }

  if (error || !playlist) {
    return (
      <Container maxWidth="lg" sx={{ py: 4 }}>
        <Box sx={{ display: 'flex', flexDirection: 'column', alignItems: 'center', gap: 2 }}>
          <Typography color="error">{error || 'Playlist not found'}</Typography>
          <Button
            startIcon={<ArrowBack />}
            onClick={() => navigate(-1)}
            variant="outlined"
          >
            Назад
          </Button>
        </Box>
      </Container>
    );
  }

  return (
    <Container maxWidth="lg" sx={{ py: 4 }}>
      <Button
        startIcon={<ArrowBack />}
        onClick={() => navigate(-1)}
        variant="outlined"
        sx={{ mb: 3 }}
      >
        Назад
      </Button>
      
      <Grid container spacing={4}>
        {/* Playlist Header */}
        <Grid item xs={12} md={4}>
          <Card>
            <Box sx={{ position: 'relative' }}>
              <CardMedia
                component="img"
                image={playlist.coverImage || '/default-playlist.png'}
                alt={playlist.name}
                sx={{ 
                  aspectRatio: '1/1',
                  width: '100%',
                  height: 'auto',
                  objectFit: 'cover'
                }}
              />
              {isOwner && (
                <IconButton
                  sx={{
                    position: 'absolute',
                    bottom: 8,
                    right: 8,
                    bgcolor: 'rgba(0, 0, 0, 0.6)',
                    '&:hover': {
                      bgcolor: 'rgba(0, 0, 0, 0.8)',
                    },
                  }}
                  onClick={handleCoverClick}
                  disabled={uploadingCover}
                >
                  <PhotoCameraIcon sx={{ color: 'white' }} />
                </IconButton>
              )}
            </Box>
          </Card>
          <input
            type="file"
            accept="image/*"
            hidden
            ref={fileInputRef}
            onChange={handleCoverChange}
          />
        </Grid>
        <Grid item xs={12} md={8}>
          <Box sx={{ mb: 2 }}>
            <Typography variant="caption" color="text.secondary" sx={{ display: 'block', mb: 1 }}>
              плейлист
            </Typography>
            <Box sx={{ display: 'flex', alignItems: 'center', gap: 2, mb: 1 }}>
              <Typography variant="h3" component="h1">
                {playlist.name}
              </Typography>
              {isOwner && (
                <IconButton onClick={() => {
                  setEditForm({
                    name: playlist.name,
                    description: playlist.description || ''
                  });
                  setEditDialogOpen(true);
                }}>
                  <EditIcon />
                </IconButton>
              )}
            </Box>
            
            {/* Owner */}
            <Box sx={{ mb: 2 }}>
              <Typography
                variant="h6"
                component="a"
                href={`/users/${playlist.owner_id}`}
                sx={{ 
                  display: 'inline-block',
                  textDecoration: 'none', 
                  color: 'inherit',
                  '&:hover': {
                    textDecoration: 'underline',
                  }
                }}
              >
                {playlist.owner_name || 'Unknown user'}
              </Typography>
            </Box>

            {/* Description */}
            {playlist.description && (
              <Typography variant="body1" color="text.secondary" sx={{ mb: 2 }}>
                {playlist.description}
              </Typography>
            )}

            {/* Privacy Switch for owner */}
            {isOwner && (
              <Box sx={{ display: 'flex', alignItems: 'center', gap: 1, mb: 2 }}>
                <Typography>Приватный плейлист</Typography>
                <Switch
                  checked={playlist.is_private}
                  onChange={handlePrivacyChange}
                  disabled={updatingPrivacy}
                />
              </Box>
            )}
          </Box>
        </Grid>

        {/* Tracks List */}
        <Grid item xs={12}>
          {tracksLoading ? (
            <Box sx={{ display: 'flex', justifyContent: 'center', p: 3 }}>
              <CircularProgress />
            </Box>
          ) : tracks && tracks.length > 0 ? (
            <List>
              {tracks.map((track, index) => (
                <Box key={track.id}>
                  <ListItem
                    button
                    onClick={() => handleTrackClick(track.id)}
                    sx={{
                      display: 'flex',
                      alignItems: 'center',
                      gap: 2,
                      '&:hover': {
                        backgroundColor: 'action.hover',
                      },
                    }}
                  >
                    <IconButton
                      size="small"
                      onClick={(e) => {
                        e.stopPropagation();
                        handleTrackClick(track.id);
                      }}
                    >
                      {currentTrack?.id === track.id && isPlaying ? <Pause /> : <PlayArrow />}
                    </IconButton>
                    <ListItemText
                      primary={track.name}
                      secondary={formatDuration(track.duration)}
                    />
                  </ListItem>
                  {index < tracks.length - 1 && <Divider />}
                </Box>
              ))}
            </List>
          ) : (
            <Typography color="text.secondary" sx={{ p: 3, textAlign: 'center' }}>
              Нет треков в плейлисте
            </Typography>
          )}
        </Grid>
      </Grid>

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
    </Container>
  );
};

export default PlaylistPage; 