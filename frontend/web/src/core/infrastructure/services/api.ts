import axios from 'axios';

const API_BASE_URL = 'http://localhost:8080/api/v1';

export const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
  withCredentials: true,
});

// Add request interceptor for authentication
api.interceptors.request.use((config) => {
  const token = localStorage.getItem('access_token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

// Add response interceptor for token refresh
api.interceptors.response.use(
  (response) => response,
  async (error) => {
    const originalRequest = error.config;

    // Если ошибка 401 и это не запрос на обновление токена
    if (error.response?.status === 401 && !originalRequest._retry) {
      originalRequest._retry = true;

      try {
        // Запрашиваем новый access token
        // Refresh token будет автоматически отправлен в cookie
        const response = await axios.post(`${API_BASE_URL}/auth/refresh`, {}, {
          withCredentials: true // Важно для работы с cookies
        });

        const { access_token } = response.data;

        // Сохраняем новый access token
        localStorage.setItem('access_token', access_token);

        // Обновляем заголовок Authorization
        originalRequest.headers.Authorization = `Bearer ${access_token}`;

        // Повторяем оригинальный запрос
        return api(originalRequest);
      } catch (refreshError) {
        // Если не удалось обновить токен, очищаем localStorage и перенаправляем на страницу входа
        localStorage.removeItem('access_token');
        window.location.href = '/login';
        return Promise.reject(refreshError);
      }
    }

    return Promise.reject(error);
  }
);

export interface Track {
  id: number;
  title: string;
  artist: string;
  album?: string;
  duration: number;
  coverUrl: string;
}

export interface Playlist {
  id: number;
  title: string;
  description: string;
  coverUrl: string;
  tracks: Track[];
}

export const apiService = {
  // Auth
  login: async (email: string, password: string) => {
    const response = await api.post('/auth/login', { email, password });
    return response.data;
  },

  register: async (email: string, password: string, username: string) => {
    const response = await api.post('/auth/register', { email, password, username });
    return response.data;
  },

  // Tracks
  getTracks: async () => {
    const response = await api.get<Track[]>('/tracks');
    return response.data;
  },

  getTrack: async (id: number) => {
    const response = await api.get<Track>(`/tracks/${id}`);
    return response.data;
  },

  searchTracks: async (query: string) => {
    const response = await api.get<Track[]>(`/tracks/search?q=${query}`);
    return response.data;
  },

  // Playlists
  getPlaylists: async () => {
    const response = await api.get<Playlist[]>('/playlists');
    return response.data;
  },

  getPlaylist: async (id: number) => {
    const response = await api.get<Playlist>(`/playlists/${id}`);
    return response.data;
  },

  createPlaylist: async (title: string, description: string) => {
    const response = await api.post<Playlist>('/playlists', { title, description });
    return response.data;
  },

  addTrackToPlaylist: async (playlistId: number, trackId: number) => {
    const response = await api.post(`/playlists/${playlistId}/tracks`, { trackId });
    return response.data;
  },

  removeTrackFromPlaylist: async (playlistId: number, trackId: number) => {
    const response = await api.delete(`/playlists/${playlistId}/tracks/${trackId}`);
    return response.data;
  },

  // User
  getLikedTracks: async () => {
    const response = await api.get<Track[]>('/user/liked-tracks');
    return response.data;
  },

  likeTrack: async (trackId: number) => {
    const response = await api.post(`/user/liked-tracks/${trackId}`);
    return response.data;
  },

  unlikeTrack: async (trackId: number) => {
    const response = await api.delete(`/user/liked-tracks/${trackId}`);
    return response.data;
  },
};

export default apiService; 