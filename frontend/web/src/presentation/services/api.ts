import axios from 'axios';

const API_BASE_URL = 'http://localhost:8080/api/v1';

const api = axios.create({
  baseURL: API_BASE_URL,
  withCredentials: true,
  headers: {
    'Content-Type': 'application/json',
  },
});

api.interceptors.request.use((config) => {
  const token = localStorage.getItem('access_token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

api.interceptors.response.use(
  (response) => response,
  async (error) => {
    const originalRequest = error.config;

    // Проверяем, что это не запрос на вход или обновление токена
    const isAuthRequest = originalRequest.url?.includes('/auth/login') || 
                         originalRequest.url?.includes('/auth/refresh');

    if (error.response?.status === 401 && !originalRequest._retry && !isAuthRequest) {
      originalRequest._retry = true;

      try {
        // Проверяем наличие refresh токена в куках
        const response = await api.post('/auth/refresh');
        const { access_token } = response.data;
        
        localStorage.setItem('access_token', access_token);
        
        originalRequest.headers.Authorization = `Bearer ${access_token}`;
        return api(originalRequest);
      } catch (refreshError) {
        // Очищаем все токены и перенаправляем на страницу входа
        localStorage.removeItem('access_token');
        // Отправляем запрос на сервер для очистки refresh токена
        try {
          await api.post('/auth/logout');
        } catch (e) {
          console.error('Failed to clear refresh token:', e);
        }
        window.location.href = '/login';
        return Promise.reject(refreshError);
      }
    }

    return Promise.reject(error);
  }
);

export const apiService = {
  register: async (login: string, password: string) => {
    const response = await api.post('/auth/register', { login, password });
    const { access_token } = response.data;
    localStorage.setItem('access_token', access_token);
    return response.data;
  },

  login: async (login: string, password: string) => {
    const response = await api.post('/auth/login', { login, password });
    const { access_token } = response.data;
    localStorage.setItem('access_token', access_token);
    return response.data;
  },

  logout: async () => {
    try {
      // Сначала отправляем запрос на сервер для очистки refresh токена
      await api.post('/auth/logout');
    } catch (error) {
      console.error('Failed to logout on server:', error);
    } finally {
      // В любом случае очищаем локальный access токен
      localStorage.removeItem('access_token');
    }
  },

  changePassword: async (oldPassword: string, newPassword: string) => {
    const response = await api.post('/auth/password/update', {
      old: oldPassword,
      new: newPassword,
    });
    return response.data;
  },

  get: async (endpoint: string) => {
    const response = await api.get(endpoint);
    return response.data;
  },

  post: async (endpoint: string, data: unknown) => {
    const response = await api.post(endpoint, data);
    return response.data;
  },

  put: async (endpoint: string, data: unknown) => {
    const response = await api.put(endpoint, data);
    return response.data;
  },

  delete: async (endpoint: string) => {
    const response = await api.delete(endpoint);
    return response.data;
  },
}; 