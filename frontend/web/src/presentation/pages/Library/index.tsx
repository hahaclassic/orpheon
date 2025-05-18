import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import {
  Box,
  Typography,
  Card,
  CardContent,
  CardMedia,
  Tabs,
  Tab,
  Grid,
  CircularProgress,
  Alert,
} from '@mui/material';
import { apiService } from '../../../presentation/services/api';

interface Track {
  ID: string;
  Title: string;
  ArtistName: string;
  AlbumName: string;
  CoverImage: string;
  Duration: number;
}

interface TabPanelProps {
  children?: React.ReactNode;
  index: number;
  value: number;
}

function TabPanel(props: TabPanelProps) {
  const { children, value, index, ...other } = props;

  return (
    <div
      role="tabpanel"
      hidden={value !== index}
      id={`simple-tabpanel-${index}`}
      aria-labelledby={`simple-tab-${index}`}
      {...other}
    >
      {value === index && <Box sx={{ p: 3 }}>{children}</Box>}
    </div>
  );
}

const Library = () => {
  const [value, setValue] = useState(0);
  const [tracks, setTracks] = useState<Track[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const navigate = useNavigate();

  const fetchTracks = async () => {
    try {
      setLoading(true);
      setError(null);
      const response = await apiService.get('/tracks');
      const tracksData = Array.isArray(response) ? response : [];
      setTracks(tracksData);
    } catch (err) {
      console.error("Error fetching tracks:", err);
      setError("Ошибка при загрузке треков. Пожалуйста, попробуйте позже.");
      setTracks([]);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchTracks();
  }, []);

  const handleChange = (event: React.SyntheticEvent, newValue: number) => {
    setValue(newValue);
  };

  const handleTrackClick = (trackId: string) => {
    navigate(`/track/${trackId}`);
  };

  return (
    <Box sx={{ width: '100%', p: 4 }}>
      <Box sx={{ borderBottom: 1, borderColor: 'divider' }}>
        <Tabs value={value} onChange={handleChange} aria-label="library tabs">
          <Tab label="Треки" />
          <Tab label="Альбомы" />
          <Tab label="Артисты" />
        </Tabs>
      </Box>

      {error && (
        <Alert severity="error" sx={{ mt: 2 }}>
          {error}
        </Alert>
      )}

      <TabPanel value={value} index={0}>
        {loading ? (
          <Box sx={{ display: 'flex', justifyContent: 'center', p: 3 }}>
            <CircularProgress />
          </Box>
        ) : (
          <Grid container spacing={3}>
            {tracks.length === 0 ? (
              <Grid item xs={12}>
                <Typography variant="h6" textAlign="center" color="text.secondary">
                  Нет доступных треков
                </Typography>
              </Grid>
            ) : (
              tracks.map((track) => (
                <Grid item xs={12} sm={6} md={4} key={track.ID}>
                  <Card
                    sx={{
                      height: '100%',
                      display: 'flex',
                      flexDirection: 'column',
                      cursor: 'pointer',
                      '&:hover': {
                        transform: 'scale(1.02)',
                        transition: 'transform 0.2s ease-in-out',
                      },
                    }}
                    onClick={() => handleTrackClick(track.ID)}
                  >
                    <CardMedia
                      component="img"
                      height="200"
                      image={track.CoverImage || '/default-cover.jpg'}
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
      </TabPanel>
      <TabPanel value={value} index={1}>
        <Typography>Альбомы (в разработке)</Typography>
      </TabPanel>
      <TabPanel value={value} index={2}>
        <Typography>Артисты (в разработке)</Typography>
      </TabPanel>
    </Box>
  );
};

export default Library; 