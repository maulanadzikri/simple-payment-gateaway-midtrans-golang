import axios from 'axios';

// Try localhost first, if doesn't work use the IP address
const API_BASE_URL = process.env.REACT_APP_API_URL || 'http://localhost:8080/api/v1';
// Alternative: 'http://192.168.1.8:8080/api/v1'

const api = axios.create({
  baseURL: API_BASE_URL,
});

api.interceptors.request.use((config) => {
  const token = localStorage.getItem('token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

export const authAPI = {
  register: (data) => api.post('/auth/register', data),
  login: (data) => api.post('/auth/login', data),
  logout: () => api.post('/auth/logout'),
};

export const userAPI = {
  getProfile: () => api.get('/profile'),
};

export const paymentAPI = {
  create: (data) => api.post('/payments/create', data),
  createQris: (data) => api.post('/payments/qris', data),
  getStatus: (orderID) => api.get(`/payments/status/${orderID}`),
  getHistory: () => api.get('/payments/history'),
};

export default api;