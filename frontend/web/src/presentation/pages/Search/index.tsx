import { useState } from 'react';
import {
  Box,
  TextField,
  InputAdornment,
  Typography,
  Card,
  CardContent,
  CardMedia,
  styled,
} from '@mui/material';
import SearchIcon from '@mui/icons-material/Search';
import { useApi } from '../../hooks/useApi';
import { apiService } from '../../../core/infrastructure/services/api';
import LoadingSpinner from '../../components/LoadingSpinner';
import ErrorBoundary from '../../components/ErrorBoundary';
import type { Track } from '../../../core/infrastructure/services/api';

const SearchContainer = styled(Box)({
  padding: '24px',
});

const SearchInput = styled(TextField)({
  width: '100%',
  marginBottom: '32px',
  '& .MuiOutlinedInput-root': {
    borderRadius: '8px',
  },
});

const StyledCard = styled(Card)({
  backgroundColor: 'background.paper',
  transition: 'transform 0.2s ease-in-out',
  '&:hover': {
    transform: 'scale(1.02)',
    cursor: 'pointer',
  },
});

const Search = () => {
  const [query, setQuery] = useState('');
  const {
    data: searchResults,
    loading,
    error,
    execute: search,
  } = useApi<Track[]>();

  const handleSearch = (value: string) => {
    setQuery(value);
    if (value.trim()) {
      search(apiService.searchTracks(value));
    }
  };

  return (
    <ErrorBoundary>
      <SearchContainer>
        <SearchInput
          placeholder="Search for songs, artists, or albums"
          value={query}
          onChange={(e) => handleSearch(e.target.value)}
          InputProps={{
            startAdornment: (
              <InputAdornment position="start">
                <SearchIcon />
              </InputAdornment>
            ),
          }}
        />

        {loading && <LoadingSpinner />}

        {error && (
          <Box sx={{ p: 3, textAlign: 'center' }}>
            <Typography color="error">
              {error.message || 'Failed to load search results'}
            </Typography>
          </Box>
        )}

        {searchResults && searchResults.length > 0 && (
          <Box sx={{ display: 'grid', gridTemplateColumns: { xs: '1fr', sm: 'repeat(2, 1fr)', md: 'repeat(3, 1fr)' }, gap: 3 }}>
            {searchResults.map((track) => (
              <Box key={track.id}>
                <StyledCard>
                  <CardMedia
                    component="img"
                    height="200"
                    image={track.coverUrl}
                    alt={track.title}
                  />
                  <CardContent>
                    <Typography variant="h6" noWrap>
                      {track.title}
                    </Typography>
                    <Typography variant="body2" color="text.secondary" noWrap>
                      {track.artist}
                    </Typography>
                    {track.album && (
                      <Typography variant="body2" color="text.secondary" noWrap>
                        {track.album}
                      </Typography>
                    )}
                  </CardContent>
                </StyledCard>
              </Box>
            ))}
          </Box>
        )}

        {searchResults && searchResults.length === 0 && query && (
          <Box sx={{ p: 3, textAlign: 'center' }}>
            <Typography color="text.secondary">
              No results found for "{query}"
            </Typography>
          </Box>
        )}
      </SearchContainer>
    </ErrorBoundary>
  );
};

export default Search; 