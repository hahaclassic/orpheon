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
import axios from 'axios';

interface Genre {
  id: number;
  name: string;
  description: string;
}

const GenreList = () => {
  const [genres, setGenres] = useState<Genre[]>([]);
  const [open, setOpen] = useState(false);
  const [editingGenre, setEditingGenre] = useState<Genre | null>(null);
  const [formData, setFormData] = useState({ name: '', description: '' });
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(true);

  const fetchGenres = async () => {
    try {
      setLoading(true);
      const response = await axios.get('/api/v1/genres');
      // Убедимся, что данные являются массивом
      const genresData = Array.isArray(response.data) ? response.data : [];
      setGenres(genresData);
      setError(null);
    } catch (err) {
      setError('Ошибка при загрузке жанров');
      console.error('Error fetching genres:', err);
      setGenres([]); // Устанавливаем пустой массив в случае ошибки
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchGenres();
  }, []);

  const handleOpen = (genre?: Genre) => {
    if (genre) {
      setEditingGenre(genre);
      setFormData({ name: genre.name, description: genre.description });
    } else {
      setEditingGenre(null);
      setFormData({ name: '', description: '' });
    }
    setOpen(true);
  };

  const handleClose = () => {
    setOpen(false);
    setEditingGenre(null);
    setFormData({ name: '', description: '' });
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
      };

      if (editingGenre) {
        await axios.put(`/api/v1/genres/${editingGenre.id}`, trimmedData);
      } else {
        await axios.post('/api/v1/genres', trimmedData);
      }
      handleClose();
      fetchGenres();
    } catch (err) {
      setError('Ошибка при сохранении жанра');
      console.error('Error saving genre:', err);
    }
  };

  const handleDelete = async (id: number) => {
    if (window.confirm('Вы уверены, что хотите удалить этот жанр?')) {
      try {
        await axios.delete(`/api/v1/genres/${id}`);
        fetchGenres();
      } catch (err) {
        setError('Ошибка при удалении жанра');
        console.error('Error deleting genre:', err);
      }
    }
  };

  return (
    <Box sx={{ p: 4, maxWidth: 1200, mx: 'auto' }}>
      <Paper sx={{ p: 4 }} elevation={3}>
        <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 3 }}>
          <Typography variant="h5" fontWeight={700}>
            Управление жанрами
          </Typography>
          <Button variant="contained" color="primary" onClick={() => handleOpen()}>
            Добавить жанр
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
                  <TableCell>Название</TableCell>
                  <TableCell>Описание</TableCell>
                  <TableCell align="right">Действия</TableCell>
                </TableRow>
              </TableHead>
              <TableBody>
                {genres.length === 0 ? (
                  <TableRow>
                    <TableCell colSpan={3} align="center">
                      Нет доступных жанров
                    </TableCell>
                  </TableRow>
                ) : (
                  genres.map((genre) => (
                    <TableRow key={genre.id}>
                      <TableCell>{genre.name}</TableCell>
                      <TableCell>{genre.description}</TableCell>
                      <TableCell align="right">
                        <IconButton onClick={() => handleOpen(genre)} color="primary">
                          <EditIcon />
                        </IconButton>
                        <IconButton onClick={() => handleDelete(genre.id)} color="error">
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
          {editingGenre ? 'Редактировать жанр' : 'Добавить жанр'}
        </DialogTitle>
        <form onSubmit={handleSubmit}>
          <DialogContent>
            <TextField
              autoFocus
              margin="dense"
              label="Название"
              fullWidth
              value={formData.name}
              onChange={handleChange}
              name="name"
              required
            />
            <TextField
              margin="dense"
              label="Описание"
              fullWidth
              multiline
              rows={3}
              value={formData.description}
              onChange={handleChange}
              name="description"
            />
          </DialogContent>
          <DialogActions>
            <Button onClick={handleClose}>Отмена</Button>
            <Button type="submit" variant="contained" color="primary">
              {editingGenre ? 'Сохранить' : 'Добавить'}
            </Button>
          </DialogActions>
        </form>
      </Dialog>
    </Box>
  );
};

export default GenreList; 