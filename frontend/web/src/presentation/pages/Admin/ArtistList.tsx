import { useState, useEffect } from 'react';
import {
  Box,
  Typography,
  Paper,
  Button,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  IconButton,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  TextField,
  Alert,
  CircularProgress,
} from '@mui/material';
import { Edit as EditIcon, Delete as DeleteIcon } from '@mui/icons-material';
import { api } from '../../../core/infrastructure/services/api';

interface Artist {
  ID: string;
  Name: string;
  Description: string;
  Country: string;
}

const ArtistList = () => {
  const [artists, setArtists] = useState<Artist[]>([]);
  const [open, setOpen] = useState(false);
  const [editingArtist, setEditingArtist] = useState<Artist | null>(null);
  const [formData, setFormData] = useState({
    name: '',
    description: '',
    country: '',
  });
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(true);

  const fetchArtists = async () => {
    try {
      setLoading(true);
      const response = await api.get('/artists');
      const artistsData = Array.isArray(response.data) ? response.data : [];
      setArtists(artistsData);
      setError(null);
    } catch (err) {
      setError('Ошибка при загрузке артистов');
      console.error('Error fetching artists:', err);
      setArtists([]);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchArtists();
  }, []);

  const handleOpen = (artist?: Artist) => {
    if (artist) {
      setEditingArtist(artist);
      setFormData({
        name: artist.Name,
        description: artist.Description,
        country: artist.Country,
      });
    } else {
      setEditingArtist(null);
      setFormData({
        name: '',
        description: '',
        country: '',
      });
    }
    setOpen(true);
  };

  const handleClose = () => {
    setOpen(false);
    setEditingArtist(null);
    setFormData({
      name: '',
      description: '',
      country: '',
    });
    setError(null);
  };

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setFormData({ ...formData, [e.target.name]: e.target.value });
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      const trimmedData = {
        Name: formData.name.trim(),
        Description: formData.description.trim(),
        Country: formData.country.trim(),
      };

      if (editingArtist) {
        await api.put(`/artists/${editingArtist.ID}`, trimmedData);
      } else {
        await api.post('/artists', trimmedData);
      }
      handleClose();
      fetchArtists();
    } catch (err) {
      setError('Ошибка при сохранении артиста');
      console.error('Error saving artist:', err);
    }
  };

  const handleDelete = async (id: string) => {
    if (window.confirm('Вы уверены, что хотите удалить этого артиста?')) {
      try {
        await api.delete(`/artists/${id}`);
        fetchArtists();
      } catch (err) {
        setError('Ошибка при удалении артиста');
        console.error('Error deleting artist:', err);
      }
    }
  };

  return (
    <Box sx={{ p: 4, maxWidth: 1200, mx: 'auto' }}>
      <Paper sx={{ p: 4 }} elevation={3}>
        <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 3, gap: 2 }}>
          <Typography variant="h5" fontWeight={700}>
            Управление артистами
          </Typography>
          <Button variant="contained" color="primary" onClick={() => handleOpen()}>
            Добавить артиста
          </Button>
        </Box>

        {error && (
          <Alert severity="error" sx={{ mb: 2 }}>
            {error}
          </Alert>
        )}

        {loading ? (
          <Box sx={{ display: 'flex', justifyContent: 'center', p: 3 }}>
            <CircularProgress />
          </Box>
        ) : (
          <TableContainer>
            <Table>
              <TableHead>
                <TableRow>
                  <TableCell>Имя</TableCell>
                  <TableCell>Страна</TableCell>
                  <TableCell>Описание</TableCell>
                  <TableCell align="right">Действия</TableCell>
                </TableRow>
              </TableHead>
              <TableBody>
                {artists.length === 0 ? (
                  <TableRow>
                    <TableCell colSpan={4} align="center">
                      Нет доступных артистов
                    </TableCell>
                  </TableRow>
                ) : (
                  artists.map((artist) => (
                    <TableRow key={artist.ID}>
                      <TableCell>{artist.Name}</TableCell>
                      <TableCell>{artist.Country}</TableCell>
                      <TableCell>{artist.Description}</TableCell>
                      <TableCell align="right">
                        <IconButton onClick={() => handleOpen(artist)} color="primary">
                          <EditIcon />
                        </IconButton>
                        <IconButton onClick={() => handleDelete(artist.ID)} color="error">
                          <DeleteIcon />
                        </IconButton>
                      </TableCell>
                    </TableRow>
                  ))
                )}
              </TableBody>
            </Table>
          </TableContainer>
        )}
      </Paper>

      <Dialog open={open} onClose={handleClose} maxWidth="sm" fullWidth>
        <DialogTitle>
          {editingArtist ? 'Редактировать артиста' : 'Добавить артиста'}
        </DialogTitle>
        <form onSubmit={handleSubmit}>
          <DialogContent>
            <TextField
              autoFocus
              margin="dense"
              label="Имя"
              fullWidth
              value={formData.name}
              onChange={handleChange}
              name="name"
              required
              sx={{ mb: 2 }}
            />
            <TextField
              margin="dense"
              label="Страна"
              fullWidth
              value={formData.country}
              onChange={handleChange}
              name="country"
              required
              sx={{ mb: 2 }}
            />
            <TextField
              margin="dense"
              label="Описание"
              fullWidth
              value={formData.description}
              onChange={handleChange}
              name="description"
              required
              multiline
              rows={4}
            />
          </DialogContent>
          <DialogActions>
            <Button onClick={handleClose}>Отмена</Button>
            <Button type="submit" variant="contained" color="primary">
              {editingArtist ? 'Сохранить' : 'Добавить'}
            </Button>
          </DialogActions>
        </form>
      </Dialog>
    </Box>
  );
};

export default ArtistList; 