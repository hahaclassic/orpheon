import { useEffect, useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { Box, Container, Typography, Grid, Card, CardContent, CardMedia, List, ListItem, ListItemText, ListItemAvatar, Avatar, Chip, Button } from '@mui/material';
import { MusicNote, Person, ArrowBack } from '@mui/icons-material';
import { albumService } from '../../core/infrastructure/services/albumService';
import type { Album, Artist, Genre } from '../../core/infrastructure/services/albumService';

const AlbumPage = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const [album, setAlbum] = useState<Album | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchAlbumData = async () => {
      if (!id) {
        setError('Album ID is missing');
        setLoading(false);
        return;
      }
      
      try {
        setLoading(true);
        const data = await albumService.getAlbum(id);
        setAlbum(data);
        setError(null);
      } catch (error) {
        console.error('Error fetching album data:', error);
        setError('Failed to load album data');
        setAlbum(null);
      } finally {
        setLoading(false);
      }
    };

    fetchAlbumData();
  }, [id]);

  if (loading) {
    return (
      <Container maxWidth="lg" sx={{ py: 4 }}>
        <Typography>Loading...</Typography>
      </Container>
    );
  }

  if (error || !album) {
    return (
      <Container maxWidth="lg" sx={{ py: 4 }}>
        <Box sx={{ display: 'flex', flexDirection: 'column', alignItems: 'center', gap: 2 }}>
          <Typography color="error">{error || 'Album not found'}</Typography>
          <Button
            startIcon={<ArrowBack />}
            onClick={() => navigate(-1)}
            variant="outlined"
          >
            Go Back
          </Button>
        </Box>
      </Container>
    );
  }

  return (
    <Container maxWidth="lg" sx={{ py: 4 }}>
      <Grid container spacing={4}>
        {/* Album Header */}
        <Grid item xs={12} md={4}>
          <Card>
            <CardMedia
              component="img"
              image={album.imageUrl || '/default-album.png'}
              alt={album.title}
              sx={{ aspectRatio: '1/1' }}
            />
          </Card>
        </Grid>
        <Grid item xs={12} md={8}>
          <Box sx={{ mb: 2 }}>
            <Typography variant="h3" component="h1" gutterBottom>
              {album.title}
            </Typography>
            
            {/* Artists */}
            <Box sx={{ mb: 2 }}>
              {album.artists.map((artist) => (
                <Typography
                  key={artist.id}
                  variant="h6"
                  component="a"
                  href={`/artists/${artist.id}`}
                  sx={{ 
                    display: 'flex', 
                    alignItems: 'center', 
                    gap: 1, 
                    textDecoration: 'none', 
                    color: 'inherit',
                    mb: 1 
                  }}
                >
                  <Person /> {artist.name}
                </Typography>
              ))}
            </Box>

            {/* Genres */}
            {album.genres && album.genres.length > 0 && (
              <Box sx={{ mt: 2, display: 'flex', gap: 1, flexWrap: 'wrap' }}>
                {album.genres.map((genre) => (
                  <Chip 
                    key={genre.id} 
                    label={genre.title} 
                    size="small" 
                  />
                ))}
              </Box>
            )}

            {/* Release Date */}
            {album.release_date && (
              <Typography variant="body2" color="text.secondary" sx={{ mt: 2 }}>
                Released: {new Date(album.release_date).toLocaleDateString()}
              </Typography>
            )}

            {/* Label */}
            {album.label && (
              <Typography variant="body2" color="text.secondary" sx={{ mt: 1 }}>
                Label: {album.label}
              </Typography>
            )}
          </Box>
        </Grid>
      </Grid>
    </Container>
  );
};

export default AlbumPage; 