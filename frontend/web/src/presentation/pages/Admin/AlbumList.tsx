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
  Select,
  MenuItem,
  FormControl,
  InputLabel,
  Chip,
  OutlinedInput,
} from '@mui/material';
import { Edit as EditIcon, Delete as DeleteIcon } from '@mui/icons-material';
import { DatePicker } from '@mui/x-date-pickers/DatePicker';
import { api } from '../../../core/infrastructure/services/api';

interface Album {
  ID: string;
  Title: string;
  Description: string;
  ReleaseDate: string;
  ArtistID: string;
  ArtistName: string;
  LicenseID: string;
  LicenseName: string;
}

interface Artist {
  ID: string;
  Name: string;
}

interface License {
  ID: string;
  Title: string;
}

const ITEM_HEIGHT = 48;
const ITEM_PADDING_TOP = 8;
const MenuProps = {
  PaperProps: {
    style: {
      maxHeight: ITEM_HEIGHT * 4.5 + ITEM_PADDING_TOP,
      width: 250,
    },
  },
};

const AlbumList = () => {
  const [albums, setAlbums] = useState<Album[]>([]);
  const [artists, setArtists] = useState<Artist[]>([]);
  const [licenses, setLicenses] = useState<License[]>([]);
  const [open, setOpen] = useState(false);
  const [editingAlbum, setEditingAlbum] = useState<Album | null>(null);
  const [formData, setFormData] = useState({
    title: '',
    description: '',
    releaseDate: '',
    artistId: '',
    licenseId: '',
    selectedArtists: [] as string[],
  });
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(true);

  const fetchAlbums = async () => {
    try {
      setLoading(true);
      const response = await api.get('/albums');
      const albumsData = Array.isArray(response.data) ? response.data : [];
      setAlbums(albumsData);
      setError(null);
    } catch (err) {
      setError('Ошибка при загрузке альбомов');
      console.error('Error fetching albums:', err);
      setAlbums([]);
    } finally {
      setLoading(false);
    }
  };

  const fetchArtists = async () => {
    try {
      const response = await api.get('/artists');
      const artistsData = Array.isArray(response.data) ? response.data : [];
      setArtists(artistsData);
    } catch (err) {
      console.error('Error fetching artists:', err);
    }
  };

  const fetchLicenses = async () => {
    try {
      const response = await api.get('/licenses');
      const licensesData = Array.isArray(response.data) ? response.data : [];
      setLicenses(licensesData);
    } catch (err) {
      console.error('Error fetching licenses:', err);
    }
  };

  useEffect(() => {
    fetchAlbums();
    fetchArtists();
    fetchLicenses();
  }, []);

  const handleOpen = (album?: Album) => {
    if (album) {
      setEditingAlbum(album);
      setFormData({
        title: album.Title,
        description: album.Description,
        releaseDate: album.ReleaseDate,
        artistId: album.ArtistID,
        licenseId: album.LicenseID,
        selectedArtists: [album.ArtistID],
      });
    } else {
      setEditingAlbum(null);
      setFormData({
        title: '',
        description: '',
        releaseDate: '',
        artistId: '',
        licenseId: '',
        selectedArtists: [],
      });
    }
    setOpen(true);
  };

  const handleClose = () => {
    setOpen(false);
    setEditingAlbum(null);
    setFormData({
      title: '',
      description: '',
      releaseDate: '',
      artistId: '',
      licenseId: '',
      selectedArtists: [],
    });
    setError(null);
  };

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setFormData({ ...formData, [e.target.name]: e.target.value });
  };

  const handleSelectChange = (e: any) => {
    setFormData({ ...formData, [e.target.name]: e.target.value });
  };

  const handleArtistsChange = (event: any) => {
    const {
      target: { value },
    } = event;
    setFormData({
      ...formData,
      selectedArtists: typeof value === 'string' ? value.split(',') : value,
    });
  };

  const handleDateChange = (date: Date | null) => {
    if (date) {
      setFormData({
        ...formData,
        releaseDate: date.toISOString().split('T')[0],
      });
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      const trimmedData = {
        Title: formData.title.trim(),
        Description: formData.description.trim(),
        ReleaseDate: formData.releaseDate.trim(),
        ArtistID: formData.selectedArtists[0], // Берем первого выбранного артиста как основного
        LicenseID: formData.licenseId.trim(),
        ArtistIDs: formData.selectedArtists, // Отправляем все выбранные ID артистов
      };

      if (editingAlbum) {
        await api.put(`/albums/${editingAlbum.ID}`, trimmedData);
      } else {
        await api.post('/albums', trimmedData);
      }
      handleClose();
      fetchAlbums();
    } catch (err) {
      setError('Ошибка при сохранении альбома');
      console.error('Error saving album:', err);
    }
  };

  const handleDelete = async (id: string) => {
    if (window.confirm('Вы уверены, что хотите удалить этот альбом?')) {
      try {
        await api.delete(`/albums/${id}`);
        fetchAlbums();
      } catch (err) {
        setError('Ошибка при удалении альбома');
        console.error('Error deleting album:', err);
      }
    }
  };

  return (
    <Box sx={{ p: 4, maxWidth: 1200, mx: 'auto' }}>
      <Paper sx={{ p: 4 }} elevation={3}>
        <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 3, gap: 2 }}>
          <Typography variant="h5" fontWeight={700}>
            Управление альбомами
          </Typography>
          <Button variant="contained" color="primary" onClick={() => handleOpen()}>
            Добавить альбом
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
                  <TableCell>Исполнитель</TableCell>
                  <TableCell>Лицензия</TableCell>
                  <TableCell>Дата выпуска</TableCell>
                  <TableCell>Описание</TableCell>
                  <TableCell align="right">Действия</TableCell>
                </TableRow>
              </TableHead>
              <TableBody>
                {albums.length === 0 ? (
                  <TableRow>
                    <TableCell colSpan={6} align="center">
                      Нет доступных альбомов
                    </TableCell>
                  </TableRow>
                ) : (
                  albums.map((album) => (
                    <TableRow key={album.ID}>
                      <TableCell>{album.Title}</TableCell>
                      <TableCell>{album.ArtistName}</TableCell>
                      <TableCell>{album.LicenseName}</TableCell>
                      <TableCell>{new Date(album.ReleaseDate).toLocaleDateString()}</TableCell>
                      <TableCell>{album.Description}</TableCell>
                      <TableCell align="right">
                        <IconButton onClick={() => handleOpen(album)} color="primary">
                          <EditIcon />
                        </IconButton>
                        <IconButton onClick={() => handleDelete(album.ID)} color="error">
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
          {editingAlbum ? 'Редактировать альбом' : 'Добавить альбом'}
        </DialogTitle>
        <form onSubmit={handleSubmit}>
          <DialogContent>
            <TextField
              autoFocus
              margin="dense"
              label="Название"
              fullWidth
              value={formData.title}
              onChange={handleChange}
              name="title"
              required
              sx={{ mb: 2 }}
            />
            <FormControl fullWidth sx={{ mb: 2 }}>
              <InputLabel>Артисты</InputLabel>
              <Select
                multiple
                value={formData.selectedArtists}
                onChange={handleArtistsChange}
                input={<OutlinedInput label="Артисты" />}
                renderValue={(selected) => (
                  <Box sx={{ display: 'flex', flexWrap: 'wrap', gap: 0.5 }}>
                    {selected.map((value) => (
                      <Chip
                        key={value}
                        label={artists.find(artist => artist.ID === value)?.Name || value}
                      />
                    ))}
                  </Box>
                )}
                MenuProps={MenuProps}
              >
                {artists.map((artist) => (
                  <MenuItem key={artist.ID} value={artist.ID}>
                    {artist.Name}
                  </MenuItem>
                ))}
              </Select>
            </FormControl>
            <FormControl fullWidth sx={{ mb: 2 }}>
              <InputLabel>Лицензия</InputLabel>
              <Select
                value={formData.licenseId}
                onChange={handleSelectChange}
                name="licenseId"
                label="Лицензия"
                required
              >
                {licenses.map((license) => (
                  <MenuItem key={license.ID} value={license.ID}>
                    {license.Title}
                  </MenuItem>
                ))}
              </Select>
            </FormControl>
            <DatePicker
              label="Дата выпуска"
              value={formData.releaseDate ? new Date(formData.releaseDate) : null}
              onChange={handleDateChange}
              slotProps={{
                textField: {
                  fullWidth: true,
                  margin: 'dense',
                  required: true,
                  sx: { mb: 2 }
                }
              }}
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
              {editingAlbum ? 'Сохранить' : 'Добавить'}
            </Button>
          </DialogActions>
        </form>
      </Dialog>
    </Box>
  );
};

export default AlbumList; 