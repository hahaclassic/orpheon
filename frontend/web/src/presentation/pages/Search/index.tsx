import { useState, useEffect } from "react";
import {
  Box,
  Container,
  Typography,
  TextField,
  InputAdornment,
  Grid,
  Card,
  CardContent,
  CardMedia,
  Button,
  CircularProgress,
  Alert,
} from "@mui/material";
import { useNavigate } from "react-router-dom";
import SearchIcon from "@mui/icons-material/Search";
import { apiService } from "../../../presentation/services/api";
import LoadingSpinner from "../../components/LoadingSpinner";
import ErrorBoundary from "../../components/ErrorBoundary";

interface Track {
  ID: string;
  Title: string;
  ArtistName: string;
  AlbumName: string;
  CoverImage: string;
  Duration: number;
}

const Search = () => {
  const [searchQuery, setSearchQuery] = useState("");
  const [searchResults, setSearchResults] = useState<Track[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const navigate = useNavigate();

  const handleSearch = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!searchQuery.trim()) {
      setSearchResults([]);
      return;
    }

    try {
      setLoading(true);
      setError(null);
      const response = await apiService.get(`/search?query=${encodeURIComponent(searchQuery)}`);
      const results = Array.isArray(response) ? response : [];
      setSearchResults(results);
    } catch (err) {
      console.error("Error fetching search results:", err);
      setError("Ошибка при поиске. Пожалуйста, попробуйте позже.");
      setSearchResults([]);
    } finally {
      setLoading(false);
    }
  };

  const handleTrackClick = (trackId: string) => {
    navigate(`/track/${trackId}`);
  };

  return (
    <Container maxWidth="lg" sx={{ py: 4 }}>
      <Box component="form" onSubmit={handleSearch} sx={{ mb: 4 }}>
        <TextField
          fullWidth
          variant="outlined"
          placeholder="Поиск треков..."
          value={searchQuery}
          onChange={(e) => setSearchQuery(e.target.value)}
          InputProps={{
            startAdornment: (
              <InputAdornment position="start">
                <SearchIcon />
              </InputAdornment>
            ),
            endAdornment: (
              <InputAdornment position="end">
                <Button type="submit" variant="contained" color="primary">
                  Поиск
                </Button>
              </InputAdornment>
            ),
          }}
        />
      </Box>

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
        <Grid container spacing={3}>
          {searchResults.length === 0 ? (
            <Grid item xs={12}>
              <Typography variant="h6" textAlign="center" color="text.secondary">
                {searchQuery ? "Ничего не найдено" : "Введите запрос для поиска"}
              </Typography>
            </Grid>
          ) : (
            searchResults.map((track) => (
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
      )}
    </Container>
  );
};

export default Search; 