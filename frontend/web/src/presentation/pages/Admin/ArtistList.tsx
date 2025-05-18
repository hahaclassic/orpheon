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
import { api } from '../../../presentation/services/api';

interface Artist {
  id: string;
  name: string;
  description: string;
  country: string;
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
      const response = await api.getArtists();
      console.log('Raw API response:', response);
      const artistsData = Array.isArray(response) ? response : [];
      console.log('Processed artists data:', artistsData);
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
        name: artist.name,
        description: artist.description,
        country: artist.country,
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
        name: formData.name.trim(),
        description: formData.description.trim(),
        country: formData.country.trim(),
      };

      if (editingArtist) {
        await api.updateArtist(editingArtist.id, trimmedData);
      } else {
        await api.createArtist(trimmedData);
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
        await api.deleteArtist(id);
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
                  <TableCell>ID</TableCell>
                  <TableCell>Имя</TableCell>
                  <TableCell>Страна</TableCell>
                  <TableCell>Описание</TableCell>
                  <TableCell align="right">Действия</TableCell>
                </TableRow>
              </TableHead>
              <TableBody>
                {artists.length === 0 ? (
                  <TableRow>
                    <TableCell colSpan={5} align="center">
                      Нет доступных артистов
                    </TableCell>
                  </TableRow>
                ) : (
                  artists.map((artist) => (
                    <TableRow key={artist.id}>
                      <TableCell>{artist.id}</TableCell>
                      <TableCell>{artist.name}</TableCell>
                      <TableCell>{artist.country}</TableCell>
                      <TableCell>{artist.description}</TableCell>
                      <TableCell align="right">
                        <IconButton onClick={() => handleOpen(artist)} color="primary">
                          <EditIcon />
                        </IconButton>
                        <IconButton onClick={() => handleDelete(artist.id)} color="error">
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
              name="name"
              label="Имя"
              value={formData.name}
              onChange={handleChange}
              fullWidth
              margin="normal"
              required
            />
            <TextField
              name="country"
              label="Страна"
              value={formData.country}
              onChange={handleChange}
              fullWidth
              margin="normal"
              required
            />
            <TextField
              name="description"
              label="Описание"
              value={formData.description}
              onChange={handleChange}
              fullWidth
              margin="normal"
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