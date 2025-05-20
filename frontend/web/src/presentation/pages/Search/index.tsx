import { useState } from "react";
import {
  Box,
  Container,
  Typography,
  TextField,
  InputAdornment,
  Grid,
  Button,
  CircularProgress,
  Alert,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
} from "@mui/material";
import { useNavigate } from "react-router-dom";
import SearchIcon from "@mui/icons-material/Search";
import { apiService } from "../../../presentation/services/api";
import LoadingSpinner from "../../components/LoadingSpinner";
import ErrorBoundary from "../../components/ErrorBoundary";
import AlbumCard from "../../components/ContentCards/AlbumCard";
import PlaylistCard from "../../components/ContentCards/PlaylistCard";
import ArtistCard from "../../components/ContentCards/ArtistCard";
import TrackList from "../../components/ContentCards/TrackList";
import { useSearch } from "../../contexts/SearchContext";

type ContentType = "track" | "album" | "playlist" | "artist";

const Search = () => {
  const { searchState, setSearchState } = useSearch();
  const [searchQuery, setSearchQuery] = useState(searchState.query);
  const [country, setCountry] = useState(searchState.country);
  const [genre, setGenre] = useState(searchState.genre);
  const [contentType, setContentType] = useState<ContentType>(searchState.contentType);
  const navigate = useNavigate();

  const handleSearch = async (e: React.FormEvent) => {
    e.preventDefault();

    try {
      setSearchState({
        ...searchState,
        loading: true,
        error: null,
      });

      const params = new URLSearchParams({
        query: searchQuery,
        type: contentType,
        ...(country && { country }),
        ...(genre && { genre }),
      });
      
      const response = await apiService.get(`/search?${params.toString()}`);
      const results = Array.isArray(response) ? response : [];
      
      setSearchState({
        query: searchQuery,
        results,
        contentType,
        country,
        genre,
        loading: false,
        error: null,
      });
    } catch (err) {
      console.error("Error fetching search results:", err);
      setSearchState({
        ...searchState,
        loading: false,
        error: "Ошибка при поиске. Пожалуйста, попробуйте позже.",
        results: [],
      });
    }
  };

  const renderResults = () => {
    if (searchState.results.length === 0) {
      return (
        <Typography variant="h6" textAlign="center" color="text.secondary">
          {searchQuery ? "Ничего не найдено" : "Введите запрос для поиска"}
        </Typography>
      );
    }

    switch (contentType) {
      case "track":
        return <TrackList tracks={searchState.results} />;

      case "album":
        return (
          <Grid container spacing={3}>
            {searchState.results.map((album) => (
              <Grid item xs={12} sm={6} md={3} key={album.id}>
                <AlbumCard
                  id={album.id}
                  title={album.title}
                  label={album.label}
                  release_date={album.release_date}
                  artists={album.artists}
                  genres={album.genres}
                />
              </Grid>
            ))}
          </Grid>
        );

      case "playlist":
        return (
          <Grid container spacing={3}>
            {searchState.results.map((playlist) => (
              <Grid item xs={12} sm={6} md={3} key={playlist.id}>
                <PlaylistCard
                  id={playlist.id}
                  name={playlist.name}
                  trackCount={playlist.track_count}
                  isFavorite={playlist.is_favorite}
                />
              </Grid>
            ))}
          </Grid>
        );

      case "artist":
        return (
          <Grid container spacing={3}>
            {searchState.results.map((artist) => (
              <Grid item xs={12} sm={6} md={3} key={artist.id}>
                <ArtistCard
                  id={artist.id}
                  name={artist.name}
                  country={artist.country}
                  genre={artist.genre}
                  albumCount={artist.albumCount}
                  coverImage={artist.coverImage}
                />
              </Grid>
            ))}
          </Grid>
        );
    }
  };

  return (
    <Container maxWidth="lg" sx={{ py: 4 }}>
      <Box component="form" onSubmit={handleSearch} sx={{ mb: 4 }}>
        <Grid container spacing={2} alignItems="center">
          <Grid item xs={12} md={4}>
            <TextField
              fullWidth
              variant="outlined"
              placeholder="Поиск..."
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              InputProps={{
                startAdornment: (
                  <InputAdornment position="start">
                    <SearchIcon />
                  </InputAdornment>
                ),
              }}
            />
          </Grid>
          
          <Grid item xs={12} md={2}>
            <FormControl fullWidth>
              <InputLabel>Тип контента</InputLabel>
              <Select
                value={contentType}
                label="Тип контента"
                onChange={(e) => setContentType(e.target.value as ContentType)}
              >
                <MenuItem value="track">Треки</MenuItem>
                <MenuItem value="album">Альбомы</MenuItem>
                <MenuItem value="playlist">Плейлисты</MenuItem>
                <MenuItem value="artist">Артисты</MenuItem>
              </Select>
            </FormControl>
          </Grid>
          
          <Grid item xs={12} md={2}>
            <FormControl fullWidth>
              <InputLabel>Страна</InputLabel>
              <Select
                value={country}
                label="Страна"
                onChange={(e) => setCountry(e.target.value)}
              >
                <MenuItem value="">Все страны</MenuItem>
                <MenuItem value="RU">Россия</MenuItem>
                <MenuItem value="US">США</MenuItem>
                <MenuItem value="GB">Великобритания</MenuItem>
                <MenuItem value="DE">Германия</MenuItem>
                <MenuItem value="FR">Франция</MenuItem>
              </Select>
            </FormControl>
          </Grid>
          
          <Grid item xs={12} md={2}>
            <FormControl fullWidth>
              <InputLabel>Жанр</InputLabel>
              <Select
                value={genre}
                label="Жанр"
                onChange={(e) => setGenre(e.target.value)}
              >
                <MenuItem value="">Все жанры</MenuItem>
                <MenuItem value="pop">Поп</MenuItem>
                <MenuItem value="rock">Рок</MenuItem>
                <MenuItem value="hiphop">Хип-хоп</MenuItem>
                <MenuItem value="electronic">Электронная</MenuItem>
                <MenuItem value="classical">Классическая</MenuItem>
                <MenuItem value="jazz">Джаз</MenuItem>
              </Select>
            </FormControl>
          </Grid>

          <Grid item xs={12} md={2}>
            <Button 
              type="submit" 
              variant="contained" 
              color="primary"
              fullWidth
              disabled={searchState.loading}
            >
              {searchState.loading ? <CircularProgress size={24} /> : "Поиск"}
            </Button>
          </Grid>
        </Grid>
      </Box>

      {searchState.error && (
        <Alert severity="error" sx={{ mb: 2 }}>
          {searchState.error}
        </Alert>
      )}

      {searchState.loading ? (
        <Box sx={{ display: "flex", justifyContent: "center", p: 3 }}>
          <CircularProgress />
        </Box>
      ) : (
        renderResults()
      )}
    </Container>
  );
};

export default Search; 