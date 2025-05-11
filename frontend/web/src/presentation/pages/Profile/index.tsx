import { useState } from 'react';
import {
  Box,
  Typography,
  Paper,
  Avatar,
  Button,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  TextField,
  Alert,
} from '@mui/material';
import LockIcon from '@mui/icons-material/Lock';
import axios from 'axios';

const MOCK_USER = {
  username: 'Demo User',
  email: 'demo@orpheon.app',
  avatar: '',
};

const Profile = () => {
  const [passwordDialogOpen, setPasswordDialogOpen] = useState(false);
  const [passwordData, setPasswordData] = useState({
    old: '',
    new: '',
    confirm: '',
  });
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState<string | null>(null);

  const handlePasswordChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    setPasswordData((prev) => ({ ...prev, [name]: value }));
  };

  const handlePasswordSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);
    setSuccess(null);
    if (passwordData.new !== passwordData.confirm) {
      setError('Новые пароли не совпадают');
      return;
    }
    setLoading(true);
    try {
      const token = localStorage.getItem('access_token');
      await axios.post(
        'http://localhost:8080/api/v1/auth/password/update',
        { old: passwordData.old, new: passwordData.new },
        {
          headers: {
            Authorization: `Bearer ${token}`,
            'Content-Type': 'application/json',
          },
        }
      );
      setSuccess('Пароль успешно изменён!');
      setPasswordData({ old: '', new: '', confirm: '' });
      setPasswordDialogOpen(false);
    } catch (err: any) {
      setError(err?.response?.data?.message || 'Ошибка при смене пароля');
    } finally {
      setLoading(false);
    }
  };

  return (
    <Box sx={{ p: 3, maxWidth: 500, mx: 'auto' }}>
      <Paper sx={{ p: 4, display: 'flex', alignItems: 'center', gap: 3, mb: 4 }}>
        <Avatar sx={{ width: 80, height: 80, fontSize: 40 }} src={MOCK_USER.avatar}>
          {MOCK_USER.username[0]}
        </Avatar>
        <Box>
          <Typography variant="h5" fontWeight={700} gutterBottom>
            {MOCK_USER.username}
          </Typography>
          <Typography color="text.secondary">{MOCK_USER.email}</Typography>
        </Box>
      </Paper>

      <Paper sx={{ p: 3, mb: 4 }}>
        <Typography variant="h6" gutterBottom>
          Настройки аккаунта
        </Typography>
        <Button
          variant="outlined"
          startIcon={<LockIcon />}
          onClick={() => setPasswordDialogOpen(true)}
        >
          Сменить пароль
        </Button>
      </Paper>

      <Dialog open={passwordDialogOpen} onClose={() => setPasswordDialogOpen(false)} maxWidth="xs" fullWidth>
        <form onSubmit={handlePasswordSubmit}>
          <DialogTitle>Смена пароля</DialogTitle>
          <DialogContent sx={{ display: 'flex', flexDirection: 'column', gap: 2, pt: 2 }}>
            {error && <Alert severity="error">{error}</Alert>}
            {success && <Alert severity="success">{success}</Alert>}
            <TextField
              label="Текущий пароль"
              name="old"
              type="password"
              value={passwordData.old}
              onChange={handlePasswordChange}
              required
              fullWidth
            />
            <TextField
              label="Новый пароль"
              name="new"
              type="password"
              value={passwordData.new}
              onChange={handlePasswordChange}
              required
              fullWidth
            />
            <TextField
              label="Подтвердите новый пароль"
              name="confirm"
              type="password"
              value={passwordData.confirm}
              onChange={handlePasswordChange}
              required
              fullWidth
            />
          </DialogContent>
          <DialogActions>
            <Button onClick={() => setPasswordDialogOpen(false)} disabled={loading}>
              Отмена
            </Button>
            <Button type="submit" variant="contained" disabled={loading}>
              Сменить
            </Button>
          </DialogActions>
        </form>
      </Dialog>
    </Box>
  );
};

export default Profile; 