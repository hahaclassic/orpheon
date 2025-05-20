import { Card, CardContent, CardMedia, Typography, Box, CircularProgress } from "@mui/material";
import { useNavigate } from "react-router-dom";
import ImageIcon from '@mui/icons-material/Image';
import { useCoverImage } from "../../hooks/useCoverImage";

interface PlaylistCardProps {
  id: number;
  name: string;
  trackCount: number;
  isFavorite?: boolean;
}

const PlaylistCard = ({ id, name, trackCount, isFavorite }: PlaylistCardProps) => {
  const navigate = useNavigate();
  const { coverUrl, loading, error } = useCoverImage('playlist', id);

  const handlePlaylistClick = () => {
    navigate(`/playlists/${id}`);
  };

  return (
    <Card
      sx={{
        height: '100%',
        display: 'flex',
        flexDirection: 'column',
        cursor: 'pointer',
        border: isFavorite ? '2px solid #a78bfa' : undefined,
        boxShadow: isFavorite ? '0 0 0 2px #a78bfa' : undefined,
        '&:hover': {
          transform: 'scale(1.02)',
          transition: 'transform 0.2s ease-in-out',
          boxShadow: isFavorite ? '0 0 0 4px #a78bfa' : undefined,
        },
      }}
      onClick={handlePlaylistClick}
    >
      {loading ? (
        <Box
          sx={{
            height: 250,
            width: '100%',
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
            bgcolor: 'primary.dark',
            borderRadius: 2,
          }}
        >
          <CircularProgress color="inherit" />
        </Box>
      ) : coverUrl ? (
        <CardMedia
          component="img"
          sx={{
            height: 250,
            width: '100%',
            objectFit: 'cover',
            aspectRatio: '1/1'
          }}
          image={coverUrl}
          alt={name}
        />
      ) : (
        <Box
          sx={{
            height: 250,
            width: '100%',
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
            bgcolor: 'primary.dark',
            borderRadius: 2,
          }}
        >
          <ImageIcon sx={{ fontSize: 64, color: 'primary.contrastText', opacity: 0.3 }} />
        </Box>
      )}
      <CardContent>
        <Typography gutterBottom variant="h6" component="div" noWrap>
          {name}
        </Typography>
        {isFavorite && (
          <Typography variant="body2" color="#a78bfa">
            ★ Избранное
          </Typography>
        )}
        <Typography variant="body2" color="text.secondary">
          {trackCount} треков
        </Typography>
      </CardContent>
    </Card>
  );
};

export default PlaylistCard; 