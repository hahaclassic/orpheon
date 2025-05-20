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
import { ArrowBack, PlayArrow, Pause, Add as AddIcon, Check, MoreVert } from '@mui/icons-material';
import axios from 'axios';
import { useAuthContext } from '../../contexts/AuthContext';
import { usePlayerContext } from '../../contexts/PlayerContext';
import { apiService } from '../../services/api';
import TrackList from '../../components/TrackList';

const API_URL = 'http://localhost:8080/api/v1';

interface Track {
  id: string;
  name: string;
  duration: number;
  track_number: number;
  artists: Artist[];
  album_id: string;
}

interface Artist {
  id: string;
  name: string;
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
  tracks: Track[];
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
  const [coverUrl, setCoverUrl] = useState<string | null>(null);
  const [playlists, setPlaylists] = useState<Playlist[]>([]);
  const [menuAnchorEl, setMenuAnchorEl] = useState<null | HTMLElement>(null);
  const [selectedTrack, setSelectedTrack] = useState<Track | null>(null);
  const [trackMenuAnchorEl, setTrackMenuAnchorEl] = useState<null | HTMLElement>(null);
  const [selectedTrackForMenu, setSelectedTrackForMenu] = useState<Track | null>(null);

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
          setCoverUrl(coverUrl);
        } catch (err) {
          console.error('Error fetching playlist cover:', err);
          setCoverUrl('/default-playlist.png');
        }
        
        // Get playlist tracks
        try {
          setTracksLoading(true);
          const tracksResponse = await apiService.get(`/playlists/${id}/tracks`);
          setPlaylist(prev => prev ? { ...prev, tracks: tracksResponse } : null);
        } catch (err) {
          console.error('Error fetching playlist tracks:', err);
          setPlaylist(prev => prev ? { ...prev, tracks: [] } : null);
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

  useEffect(() => {
    const fetchPlaylists = async () => {
      try {
        const response = await apiService.get('/me/playlists');
        setPlaylists(response);
      } catch (err) {
        console.error('Error fetching playlists:', err);
      }
    };

    fetchPlaylists();
  }, []);

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
      setCoverUrl(coverUrl);
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

      setCoverUrl('/default-playlist.png');
    } catch (err) {
      setError('Не удалось удалить обложку');
    } finally {
      setUploadingCover(false);
    }
  };

  // Меню три точки
  const handleMenuOpen = (e: React.MouseEvent<HTMLElement>) => setMenuAnchor(e.currentTarget);
  const handleMenuClose = () => {
    setMenuAnchor(null);
    setSelectedTrack(null);
  };

  const handleTrackClick = (trackId: string) => {
    const track = playlist?.tracks.find(t => t.id === trackId);
    if (track) {
      if (currentTrack?.id === trackId) {
        togglePlay();
      } else {
        setTrack(track, playlist?.tracks || []);
      }
    }
  };

  const handleAddToPlaylist = (event: React.MouseEvent<HTMLElement>, track: Track) => {
    event.stopPropagation();
    setSelectedTrack(track);
    setMenuAnchorEl(event.currentTarget);
  };

  const handlePlaylistSelect = async (playlistId: string) => {
    if (selectedTrack) {
      try {
        const playlist = playlists.find(p => p.id === playlistId);
        const isInPlaylist = playlist && isTrackInPlaylist(playlist, selectedTrack.id);

        if (isInPlaylist) {
          // Удаляем трек из плейлиста
          await apiService.delete(`/playlists/${playlistId}/tracks/${selectedTrack.id}`);
          // Обновляем локальное состояние плейлиста
          setPlaylists(playlists.map(p => {
            if (p.id === playlistId) {
              return {
                ...p,
                tracks: (p.tracks || []).filter(t => t.id !== selectedTrack.id)
              };
            }
            return p;
          }));
        } else {
          // Добавляем трек в плейлист
          await apiService.post(`/playlists/${playlistId}/tracks`, { track_id: selectedTrack.id });
          // Обновляем локальное состояние плейлиста
          setPlaylists(playlists.map(p => {
            if (p.id === playlistId) {
              return {
                ...p,
                tracks: [...(p.tracks || []), selectedTrack]
              };
            }
            return p;
          }));
        }
      } catch (err: any) {
        console.error('Error managing track in playlist:', {
          error: err,
          response: err.response?.data,
          status: err.response?.status,
          playlistId,
          trackId: selectedTrack.id
        });
      } finally {
        handleMenuClose();
      }
    }
  };

  const handleTrackMenuOpen = (event: React.MouseEvent<HTMLElement>, track: Track) => {
    event.stopPropagation();
    setSelectedTrackForMenu(track);
    setTrackMenuAnchorEl(event.currentTarget);
  };

  const handleTrackMenuClose = () => {
    setTrackMenuAnchorEl(null);
    setSelectedTrackForMenu(null);
  };

  const isTrackInPlaylist = (playlist: Playlist, trackId: string) => {
    return playlist.tracks?.some(track => track.id === trackId) || false;
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
                image={coverUrl || '/default-playlist.png'}
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
          <Typography variant="h5" sx={{ mb: 2 }}>
            Треки
          </Typography>
          <TrackList
            tracks={playlist.tracks}
            currentTrackId={currentTrack?.id}
            isPlaying={isPlaying}
            onTrackClick={handleTrackClick}
            showTrackNumber={false}
            onAddToPlaylist={handleAddToPlaylist}
            showAlbumLink={true}
          />
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

      {/* Playlist Selection Menu */}
      <Menu
        anchorEl={menuAnchorEl}
        open={Boolean(menuAnchorEl)}
        onClose={handleMenuClose}
        PaperProps={{
          sx: { 
            maxHeight: '280px',
            '& .MuiList-root': {
              padding: 0
            }
          }
        }}
        MenuListProps={{
          sx: {
            padding: 0
          }
        }}
      >
        {playlists.map((playlist) => (
          <MenuItem 
            key={playlist.id} 
            onClick={() => handlePlaylistSelect(playlist.id)}
            sx={{ 
              minHeight: '40px',
              '&:hover': {
                backgroundColor: 'action.hover'
              },
              display: 'flex',
              justifyContent: 'space-between',
              alignItems: 'center'
            }}
          >
            <span>{playlist.name}</span>
            {selectedTrack && isTrackInPlaylist(playlist, selectedTrack.id) && (
              <Check sx={{ ml: 1, color: 'white' }} />
            )}
          </MenuItem>
        ))}
      </Menu>

      {/* Track Menu */}
      <Menu
        anchorEl={trackMenuAnchorEl}
        open={Boolean(trackMenuAnchorEl)}
        onClose={handleTrackMenuClose}
        PaperProps={{
          sx: { 
            '& .MuiList-root': {
              padding: 0
            }
          }
        }}
        MenuListProps={{
          sx: {
            padding: 0
          }
        }}
      >
        {selectedTrackForMenu && (
          <MenuItem 
            onClick={() => {
              handleTrackMenuClose();
              navigate(`/albums/${selectedTrackForMenu.album_id}`);
            }}
          >
            Перейти к альбому
          </MenuItem>
        )}
      </Menu>
    </Container>
  );
};

export default PlaylistPage; 