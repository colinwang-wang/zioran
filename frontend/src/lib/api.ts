import axios from 'axios';

const isServer = typeof window === 'undefined';
const baseURL = isServer ? 'http://127.0.0.1:8080/api/v1' : '/api/v1';

const api = axios.create({ baseURL, timeout: 15000 });

api.interceptors.request.use((config) => {
  if (!isServer) {
    const token = localStorage.getItem('token');
    if (token) config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

api.interceptors.response.use(
  (res) => {
    // 后端统一返回 {code, message, data}，自动解包
    const body = res.data;
    if (body && typeof body === 'object' && 'code' in body) {
      if (body.code !== 0) {
        return Promise.reject({ response: { data: body, status: res.status } });
      }
      // 将res.data替换为body.data（实际业务数据）
      res.data = body.data;
    }
    return res;
  },
  (error) => {
    if (error.response?.status === 401 && !isServer) {
      localStorage.removeItem('token');
      localStorage.removeItem('user');
      window.location.href = '/login';
    }
    return Promise.reject(error);
  }
);

export default api;
