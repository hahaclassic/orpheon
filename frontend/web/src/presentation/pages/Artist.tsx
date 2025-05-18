import { useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';
import { Box, Container, Typography, Grid, Card, CardContent, CardMedia, List, ListItem, ListItemText, ListItemAvatar, Avatar } from '@mui/material';
import { MusicNote, Album } from '@mui/icons-material';

interface Artist {
  id: string;
  name: string;
  description: string;
  imageUrl: string;
  tracks: Track[];
  albums: Album[];
}

interface Track {
  id: string;
  title: string;
  duration: string;
  albumId: string;
  albumName: string;
}

interface Album {
  id: string;
  title: string;
  releaseDate: string;
  imageUrl: string;
}

const ArtistPage = () => {
  const { id } = useParams<{ id: string }>();
  const [artist, setArtist] = useState<Artist | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    // TODO: Fetch artist data from API
    // This is a placeholder for the actual API call
    const fetchArtistData = async () => {
      try {
        // const response = await api.getArtist(id);
        // setArtist(response.data);
        setLoading(false);
      } catch (error) {
        console.error('Error fetching artist data:', error);
        setLoading(false);
      }
    };

    fetchArtistData();
  }, [id]);

  if (loading) {
    return <Typography>Loading...</Typography>;
  }

  if (!artist) {
    return <Typography>Artist not found</Typography>;
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
              image={artist.imageUrl}
              alt={artist.name}
            />
            <Box>
              <Typography variant="h3" component="h1" gutterBottom>
                {artist.name}
              </Typography>
              <Typography variant="body1" color="text.secondary">
                {artist.description}
              </Typography>
            </Box>
          </Box>
        </Grid>

        {/* Tracks Section */}
        <Grid item xs={12} md={6}>
          <Typography variant="h5" gutterBottom sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
            <MusicNote /> Tracks
          </Typography>
          <List>
            {artist.tracks.map((track) => (
              <ListItem key={track.id} button>
                <ListItemAvatar>
                  <Avatar>
                    <MusicNote />
                  </Avatar>
                </ListItemAvatar>
                <ListItemText
                  primary={track.title}
                  secondary={`${track.albumName} • ${track.duration}`}
                />
              </ListItem>
            ))}
          </List>
        </Grid>

        {/* Albums Section */}
        <Grid item xs={12} md={6}>
          <Typography variant="h5" gutterBottom sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
            <Album /> Albums
          </Typography>
          <Grid container spacing={2}>
            {artist.albums.map((album) => (
              <Grid item xs={12} sm={6} key={album.id}>
                <Card>
                  <CardMedia
                    component="img"
                    height="140"
                    image={album.imageUrl}
                    alt={album.title}
                  />
                  <CardContent>
                    <Typography variant="h6" noWrap>
                      {album.title}
                    </Typography>
                    <Typography variant="body2" color="text.secondary">
                      {new Date(album.releaseDate).toLocaleDateString()}
                    </Typography>
                  </CardContent>
                </Card>
              </Grid>
            ))}
          </Grid>
        </Grid>
      </Grid>
    </Container>
  );
};

export default ArtistPage; 