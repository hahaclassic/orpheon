import { useEffect, useState } from 'react';
import { useParams, Link } from 'react-router-dom';
import { Box, Container, Typography, Grid, Card, CardContent, CardMedia, List, ListItem, ListItemText, ListItemAvatar, Avatar, CircularProgress } from '@mui/material';
import { MusicNote, Album } from '@mui/icons-material';
import { api } from '../services/api';
import AlbumCard from '../components/ContentCards/AlbumCard';

interface Artist {
  id: string;
  name: string;
  description: string;
  country: string;
}

interface Track {
  id: string;
  title: string;
  duration: string;
  albumId: string;
  albumName: string;
  playCount: number;
}

interface Genre {
  id: string;
  title: string;
}

interface Album {
  id: string;
  title: string;
  label: string;
  releaseDate: string;
  artists: Artist[];
  genres: Genre[];
}

const ArtistPage = () => {
  const { id } = useParams<{ id: string }>();
  const [artist, setArtist] = useState<Artist | null>(null);
  const [tracks, setTracks] = useState<Track[]>([]);
  const [albums, setAlbums] = useState<Album[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchArtistData = async () => {
      try {
        setLoading(true);
        setError(null);

        // Fetch artist data
        const artistData = await api.getArtist(id!);
        setArtist(artistData);

        // Fetch artist's tracks
        const tracksData = await api.getArtistTracks(id!);
        setTracks(tracksData || []); // Handle null response
 
        // Fetch artist's albums
        const albumsData = await api.getArtistAlbums(id!);
        setAlbums(albumsData || []); // Handle null response

        setLoading(false);
      } catch (error) {
        console.error('Error fetching artist data:', error);
        setError('Failed to load artist data');
        setLoading(false);
      }
    };

    if (id) {
      fetchArtistData();
    }
  }, [id]);

  if (loading) {
    return (
      <Box display="flex" justifyContent="center" alignItems="center" minHeight="60vh">
        <CircularProgress />
      </Box>
    );
  }

  if (error || !artist) {
    return (
      <Container maxWidth="lg" sx={{ py: 4 }}>
        <Typography color="error" variant="h5">
          {error || 'Artist not found'}
        </Typography>
      </Container>
    );
  }

  return (
    <Container maxWidth="lg" sx={{ py: 4 }}>
      <Grid container spacing={4}>
        {/* Artist Header */}
        <Grid item xs={12}>
          <Box sx={{ display: 'flex', alignItems: 'center', gap: 3 }}>
            <CardMedia
              component="img"
              sx={{ width: 200, height: 200, borderRadius: 2 }}
              image={`/api/v1/artists/${artist.id}/avatar`}
              alt={artist.name}
            />
            <Box>
              <Typography variant="h3" component="h1" gutterBottom>
                {artist.name}
              </Typography>
              <Typography variant="body1" color="text.secondary" paragraph>
                {artist.description}
              </Typography>
              <Typography variant="body2" color="text.secondary">
                Country: {artist.country}
              </Typography>
            </Box>
          </Box>
        </Grid>

        {/* Tracks Section */}
        <Grid item xs={12}>
          <Typography variant="h5" gutterBottom sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
            <MusicNote /> Tracks
          </Typography>
          <List>
            {tracks.map((track) => (
              <ListItem
                key={track.id}
                component={Link}
                to={`/tracks/${track.id}`}
                sx={{ textDecoration: 'none', color: 'inherit' }}
              >
                <ListItemAvatar>
                  <Avatar>
                    <MusicNote />
                  </Avatar>
                </ListItemAvatar>
                <ListItemText
                  primary={track.title}
                  secondary={
                    <Link to={`/albums/${track.albumId}`} style={{ textDecoration: 'none', color: 'inherit' }}>
                      {track.albumName} • {track.duration} • {track.playCount} plays
                    </Link>
                  }
                />
              </ListItem>
            ))}
          </List>
        </Grid>

        {/* Albums Section */}
        <Grid item xs={12}>
          <Typography variant="h5" gutterBottom sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
            <Album /> Albums
          </Typography>
          <Grid container spacing={2}>
            {albums.map((album) => (
              <Grid item xs={12} sm={6} md={4} lg={3} key={album.id}>
                <AlbumCard
                  id={album.id}
                  title={album.title}
                  label={album.label}
                  release_date={album.releaseDate}
                  artists={album.artists}
                  genres={album.genres}
                />
              </Grid>
            ))}
          </Grid>
        </Grid>
      </Grid>
    </Container>
  );
};

export default ArtistPage; 