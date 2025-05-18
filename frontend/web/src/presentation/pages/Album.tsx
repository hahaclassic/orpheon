import { useEffect, useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { 
  Box, 
  Container, 
  Typography, 
  Grid, 
  Card, 
  CardMedia, 
  Chip, 
  Button,
  List,
  ListItem,
  ListItemText,
  Divider,
  IconButton,
} from '@mui/material';
import { ArrowBack, PlayArrow, Pause } from '@mui/icons-material';
import { albumService } from '../../core/infrastructure/services/albumService';
import type { Album, Artist, Genre } from '../../core/infrastructure/services/albumService';
import { apiService } from '../services/api';

interface Track {
  id: string;
  name: string;
  duration: number;
  track_number: number;
}

const AlbumPage = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const [album, setAlbum] = useState<Album | null>(null);
  const [tracks, setTracks] = useState<Track[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [coverUrl, setCoverUrl] = useState<string | null>(null);
  const [currentTrack, setCurrentTrack] = useState<string | null>(null);
  const [isPlaying, setIsPlaying] = useState(false);

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
        
        // Получаем обложку альбома
        try {
          const coverResponse = await apiService.get(`/albums/${id}/cover`, {
            responseType: 'blob'
          });
          const coverUrl = URL.createObjectURL(coverResponse);
          setCoverUrl(coverUrl);
        } catch (err) {
          console.error('Error fetching album cover:', err);
          setCoverUrl('/default-album.png');
        }

        // Получаем треки альбома
        try {
          const tracksResponse = await apiService.get(`/albums/${id}/tracks`);
          setTracks(tracksResponse);
        } catch (err) {
          console.error('Error fetching album tracks:', err);
          setTracks([]);
        }

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

  const formatDuration = (seconds: number) => {
    const minutes = Math.floor(seconds / 60);
    const remainingSeconds = seconds % 60;
    return `${minutes}:${remainingSeconds.toString().padStart(2, '0')}`;
  };

  const handleTrackClick = (trackId: string) => {
    if (currentTrack === trackId) {
      setIsPlaying(!isPlaying);
    } else {
      setCurrentTrack(trackId);
      setIsPlaying(true);
    }
  };

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
        {/* Album Header */}
        <Grid item xs={12} md={4}>
          <Card>
            <CardMedia
              component="img"
              image={coverUrl || '/default-album.png'}
              alt={album.title}
              sx={{ 
                aspectRatio: '1/1',
                width: '100%',
                height: 'auto',
                objectFit: 'cover'
              }}
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
              {album.artists.map((artist, index) => (
                <Typography
                  key={artist.id}
                  variant="h6"
                  component="a"
                  href={`/artists/${artist.id}`}
                  sx={{ 
                    display: 'inline-block',
                    textDecoration: 'none', 
                    color: 'inherit',
                    '&:hover': {
                      textDecoration: 'underline',
                    },
                    '&:not(:last-child)::after': {
                      content: '", "',
                      color: 'text.secondary',
                      marginRight: '4px'
                    }
                  }}
                >
                  {artist.name}
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
                    sx={{ 
                      backgroundColor: 'primary.light',
                      color: 'primary.contrastText',
                      '&:hover': {
                        backgroundColor: 'primary.main',
                      }
                    }}
                  />
                ))}
              </Box>
            )}

            {/* Release Date */}
            {album.release_date && (
              <Typography variant="body2" color="text.secondary" sx={{ mt: 2 }}>
                Дата релиза: {new Date(album.release_date).toLocaleDateString()}
              </Typography>
            )}

            {/* Label */}
            {album.label && (
              <Typography variant="body2" color="text.secondary" sx={{ mt: 1 }}>
                Лейбл: {album.label}
              </Typography>
            )}
          </Box>
        </Grid>

        {/* Tracks */}
        <Grid item xs={12}>
          <Typography variant="h5" sx={{ mb: 2 }}>
            Треки
          </Typography>
          <List>
            {tracks.map((track, index) => (
              <Box key={track.id}>
                <ListItem
                  sx={{
                    cursor: 'pointer',
                    '&:hover': {
                      backgroundColor: 'action.hover',
                    },
                  }}
                  onClick={() => handleTrackClick(track.id)}
                >
                  <Box
                    sx={{
                      width: 40,
                      height: 40,
                      display: 'flex',
                      alignItems: 'center',
                      justifyContent: 'center',
                      mr: 2,
                      position: 'relative',
                      '&:hover .play-icon': {
                        opacity: 1,
                      },
                      '&:hover .track-number': {
                        opacity: 0,
                      },
                    }}
                  >
                    <Typography
                      className="track-number"
                      sx={{
                        position: 'absolute',
                        transition: 'opacity 0.2s',
                      }}
                    >
                      {track.track_number}
                    </Typography>
                    <IconButton
                      size="small"
                      className="play-icon"
                      sx={{
                        position: 'absolute',
                        opacity: 0,
                        transition: 'opacity 0.2s',
                      }}
                    >
                      {currentTrack === track.id && isPlaying ? <Pause /> : <PlayArrow />}
                    </IconButton>
                  </Box>
                  <ListItemText
                    primary={track.name}
                    secondary={formatDuration(track.duration)}
                  />
                </ListItem>
                {index < tracks.length - 1 && <Divider />}
              </Box>
            ))}
          </List>
        </Grid>
      </Grid>
    </Container>
  );
};

export default AlbumPage; 