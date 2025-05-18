import { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import {
  Box,
  Container,
  Typography,
  Grid,
  Card,
  CardContent,
  CardMedia,
  CircularProgress,
  Alert,
} from "@mui/material";
import { apiService } from "../../../presentation/services/api";

interface Track {
  ID: string;
  Title: string;
  ArtistName: string;
  AlbumName: string;
  CoverImage: string;
  Duration: number;
}

interface Album {
  ID: string;
  Title: string;
  ArtistName: string;
  CoverImage: string;
  ReleaseDate: string;
}

const Home = () => {
  const [recentTracks, setRecentTracks] = useState<Track[]>([]);
  const [popularAlbums, setPopularAlbums] = useState<Album[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const navigate = useNavigate();

  const fetchRecentTracks = async () => {
    try {
      const response = await apiService.get('/tracks/recent');
      const tracksData = Array.isArray(response) ? response : [];
      setRecentTracks(tracksData);
    } catch (err) {
      console.error("Error fetching recent tracks:", err);
      setError("Ошибка при загрузке последних треков");
    }
  };

  const fetchPopularAlbums = async () => {
    try {
      const response = await apiService.get('/albums/popular');
      const albumsData = Array.isArray(response) ? response : [];
      setPopularAlbums(albumsData);
    } catch (err) {
      console.error("Error fetching popular albums:", err);
      setError("Ошибка при загрузке популярных альбомов");
    }
  };

  useEffect(() => {
    const fetchData = async () => {
      setLoading(true);
      setError(null);
      try {
        await Promise.all([fetchRecentTracks(), fetchPopularAlbums()]);
      } catch (err) {
        console.error("Error fetching data:", err);
        setError("Ошибка при загрузке данных");
      } finally {
        setLoading(false);
      }
    };

    fetchData();
  }, []);

  const handleTrackClick = (trackId: string) => {
    navigate(`/track/${trackId}`);
  };

  const handleAlbumClick = (albumId: string) => {
    navigate(`/album/${albumId}`);
  };

  return (
    <Container maxWidth="lg" sx={{ py: 4 }}>
      {error && (
        <Alert severity="error" sx={{ mb: 2 }}>
          {error}
        </Alert>
      )}

      {loading ? (
        <Box sx={{ display: "flex", justifyContent: "center", p: 3 }}>
          <CircularProgress />
        </Box>
      ) : (
        <>
          <Typography variant="h4" component="h1" gutterBottom>
            Последние треки
          </Typography>
          <Grid container spacing={3} sx={{ mb: 4 }}>
            {recentTracks.length === 0 ? (
              <Grid item xs={12}>
                <Typography variant="h6" textAlign="center" color="text.secondary">
                  Нет доступных треков
                </Typography>
              </Grid>
            ) : (
              recentTracks.map((track) => (
                <Grid item xs={12} sm={6} md={4} key={track.ID}>
                  <Card
                    sx={{
                      height: "100%",
                      display: "flex",
                      flexDirection: "column",
                      cursor: "pointer",
                      "&:hover": {
                        transform: "scale(1.02)",
                        transition: "transform 0.2s ease-in-out",
                      },
                    }}
                    onClick={() => handleTrackClick(track.ID)}
                  >
                    <CardMedia
                      component="img"
                      height="200"
                      image={track.CoverImage || "/default-cover.jpg"}
                      alt={track.Title}
                    />
                    <CardContent>
                      <Typography gutterBottom variant="h6" component="div" noWrap>
                        {track.Title}
                      </Typography>
                      <Typography variant="body2" color="text.secondary" noWrap>
                        {track.ArtistName}
                      </Typography>
                      <Typography variant="body2" color="text.secondary" noWrap>
                        {track.AlbumName}
                      </Typography>
                    </CardContent>
                  </Card>
                </Grid>
              ))
            )}
          </Grid>

          <Typography variant="h4" component="h1" gutterBottom>
            Популярные альбомы
          </Typography>
          <Grid container spacing={3}>
            {popularAlbums.length === 0 ? (
              <Grid item xs={12}>
                <Typography variant="h6" textAlign="center" color="text.secondary">
                  Нет доступных альбомов
                </Typography>
              </Grid>
            ) : (
              popularAlbums.map((album) => (
                <Grid item xs={12} sm={6} md={4} key={album.ID}>
                  <Card
                    sx={{
                      height: "100%",
                      display: "flex",
                      flexDirection: "column",
                      cursor: "pointer",
                      "&:hover": {
                        transform: "scale(1.02)",
                        transition: "transform 0.2s ease-in-out",
                      },
                    }}
                    onClick={() => handleAlbumClick(album.ID)}
                  >
                    <CardMedia
                      component="img"
                      height="200"
                      image={album.CoverImage || "/default-cover.jpg"}
                      alt={album.Title}
                    />
                    <CardContent>
                      <Typography gutterBottom variant="h6" component="div" noWrap>
                        {album.Title}
                      </Typography>
                      <Typography variant="body2" color="text.secondary" noWrap>
                        {album.ArtistName}
                      </Typography>
                      <Typography variant="body2" color="text.secondary">
                        {new Date(album.ReleaseDate).toLocaleDateString()}
                      </Typography>
                    </CardContent>
                  </Card>
                </Grid>
              ))
            )}
          </Grid>
        </>
      )}
    </Container>
  );
};

export default Home; 