import axios from 'axios';
import { useAuthStore } from './stores/auth';
import router from './router';

const api = axios.create({
  baseURL: 'http://localhost:8080', 
  headers: {
    'Content-Type': 'application/json',
  },
  withCredentials: true,
});


export const checkUsername = (username) => {
  return api.post('/api/check-username', { username });
};

export const checkEmailOnServer = (email) => {
    return api.post('/api/check-email', {email})
}

export const registerUser = (userData) => {
  return api.post('/api/reg-user', userData);
};

export const sendVerificationCodeOnServer = (email) => {
    return api.post('/api/send-verif-code', {email})
};

export const VerifyCodeOnServer = (email, code) => {
    return api.post('/api/check-verif-code', {email, code})
}

export const  loginUser = (credentials) => {
  return api.post('/api/login', credentials)
}

export const refreshTokensByServer =() => {
  return api.post('/api/refresh-tokens')
}

export const logoutUser = () => {
  return api.post('/api/logout')
}



api.interceptors.response.use(
  response => response,
  async error => {
    const originalRequest = error.config;

    if(error.response?.status === 401 && !originalRequest._retry) {
      originalRequest._retry = true;

    try {
      const {data} = await refreshTokensByServer()
      localStorage.setItem('accessToken', data.accessToken)
      api.defaults.headers.common['Authorization'] = `Bearer ${data.accessToken}`
      return api(originalRequest)
    }
    catch (refreshError) {
      localStorage.removeItem('accessToken')
      delete api.defaults.headers.common['Authorization']
      router.push('/')
      const auth = useAuthStore()
      auth.openAuthModal()
      return Promise.reject(refreshError)
      }
    }
    return Promise.reject(error)
  }
)

export default api;