import axios from 'axios';
import type { AxiosResponse, AxiosInstance, AxiosRequestConfig } from 'axios';

const API_URL = import.meta.env.VITE_API_URL;

interface User {
  id: string;
  name: string;
  registration_date: string;
  birth_date: string;
  access_lvl: number;
}

interface Genre {
  id: string;
  title: string;
}

interface License {
  id: string;
  title: string;
  description: string;
}

interface Artist {
  id: string;
  name: string;
  description: string;
  country: string;
}

interface Playlist {
  id: string;
  title: string;
  // Add any other necessary properties for a playlist
}

// Создаем тип для apiService, который возвращает данные напрямую
type ApiService = {
  get: <T = any>(url: string, config?: AxiosRequestConfig) => Promise<T>;
  post: <T = any>(url: string, data?: any, config?: AxiosRequestConfig) => Promise<T>;
  put: <T = any>(url: string, data?: any, config?: AxiosRequestConfig) => Promise<T>;
  delete: <T = any>(url: string, config?: AxiosRequestConfig) => Promise<T>;
  interceptors: {
    request: {
      use: (fulfilled: (config: AxiosRequestConfig) => AxiosRequestConfig, rejected?: (error: any) => any) => number;
    };
    response: {
      use: (fulfilled: (response: AxiosResponse) => any, rejected?: (error: any) => any) => number;
    };
  };
  (config: AxiosRequestConfig): Promise<any>;
};

const apiService = axios.create({
  baseURL: API_URL,
  withCredentials: true,
  headers: {
    'Content-Type': 'application/json',
  },
}) as unknown as ApiService;

// Request interceptor
apiService.interceptors.request.use(
  (config) => {
    console.log('[API] Request:', {
      url: config.url,
      method: config.method,
      headers: config.headers,
      withCredentials: config.withCredentials,
    });
    return config;
  },
  (error) => {
    console.error('[API] Request Error:', error);
    return Promise.reject(error);
  }
);

// Response interceptor
apiService.interceptors.response.use(
  (response) => {
    console.log('[API] Response:', {
      url: response.config.url,
      status: response.status,
      data: response.data,
      headers: response.headers,
    });
    return response.data;
  },
  async (error) => {
    console.error('[API] Error:', {
      url: error.config?.url,
      status: error.response?.status,
      data: error.response?.data,
      headers: error.response?.headers,
    });

    if (error.response?.status === 401) {
      try {
        const refreshResponse = await axios.post(
          `${API_URL}/api/v1/auth/refresh`,
          {},
          { withCredentials: true }
        );
        if (refreshResponse.data) {
          const originalRequest = error.config;
          return apiService(originalRequest);
        }
      } catch (refreshError) {
        console.error('[API] Token refresh failed:', refreshError);
      }
    }

    return Promise.reject(error);
  }
);

// API methods
export const api = {
  register: (login: string, password: string): Promise<User> =>
    apiService.post('/api/v1/auth/register', { login, password }),

  login: (login: string, password: string): Promise<User> =>
    apiService.post('/api/v1/auth/login', { login, password }),

  logout: (): Promise<void> => apiService.post('/api/v1/auth/logout'),

  changePassword: (oldPassword: string, newPassword: string): Promise<void> =>
    apiService.post('/api/v1/auth/change-password', { oldPassword, newPassword }),

  getMe: (): Promise<User> => apiService.get('/api/v1/me'),

  // User endpoints
  getUser: (id: string): Promise<User> => apiService.get(`/api/v1/users/${id}`),
  getUserPlaylists: (id: string): Promise<Playlist[]> => apiService.get(`/api/v1/users/${id}/playlists`),
  getUserFavorites: (id: string): Promise<Playlist[]> => apiService.get(`/api/v1/users/${id}/favorites`),

  // Tracks
  getTracks: () => apiService.get('/api/v1/tracks'),
  getTrack: (id: string) => apiService.get(`/api/v1/tracks/${id}`),
  createTrack: (data: FormData) =>
    apiService.post('/api/v1/tracks', data, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    }),
  updateTrack: (id: string, data: FormData) =>
    apiService.put(`/api/v1/tracks/${id}`, data, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    }),
  deleteTrack: (id: string) => apiService.delete(`/api/v1/tracks/${id}`),

  // Albums
  getAlbums: () => apiService.get('/api/v1/albums'),
  getAlbum: (id: string) => apiService.get(`/api/v1/albums/${id}`),
  createAlbum: (data: FormData) =>
    apiService.post('/api/v1/albums', data, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    }),
  updateAlbum: (id: string, data: FormData) =>
    apiService.put(`/api/v1/albums/${id}`, data, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    }),
  deleteAlbum: (id: string) => apiService.delete(`/api/v1/albums/${id}`),

  // Artists
  getArtists: () => apiService.get<Artist[]>('/api/v1/artists'),
  getArtist: (id: string) => apiService.get(`/api/v1/artists/${id}`),
  createArtist: (data: Omit<Artist, 'id'>) => apiService.post<Artist>('/api/v1/artists', data),
  updateArtist: (id: string, data: Omit<Artist, 'id'>) => apiService.put<Artist>(`/api/v1/artists/${id}`, data),
  deleteArtist: (id: string) => apiService.delete(`/api/v1/artists/${id}`),
  getArtistTracks: (id: string) => apiService.get(`/api/v1/artists/${id}/tracks`),
  getArtistAlbums: (id: string) => apiService.get(`/api/v1/artists/${id}/albums`),

  // Genres
  getGenres: (): Promise<Genre[]> => apiService.get('/api/v1/genres'),
  getGenre: (id: string): Promise<Genre> => apiService.get(`/api/v1/genres/${id}`),
  createGenre: (data: { title: string }): Promise<Genre> =>
    apiService.post('/api/v1/genres', data),
  updateGenre: (id: string, data: { title: string }): Promise<Genre> =>
    apiService.put(`/api/v1/genres/${id}`, data),
  deleteGenre: (id: string): Promise<void> => apiService.delete(`/api/v1/genres/${id}`),

  // Licenses
  getLicenses: (): Promise<License[]> => apiService.get('/api/v1/licenses'),
  getLicense: (id: string): Promise<License> => apiService.get(`/api/v1/licenses/${id}`),
  createLicense: (data: { title: string; description: string }): Promise<License> =>
    apiService.post('/api/v1/licenses', data),
  updateLicense: (id: string, data: { title: string; description: string }): Promise<License> =>
    apiService.put(`/api/v1/licenses/${id}`, data),
  deleteLicense: (id: string): Promise<void> => apiService.delete(`/api/v1/licenses/${id}`),
};

export { apiService }; 