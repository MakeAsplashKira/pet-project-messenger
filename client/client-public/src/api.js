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
  const formData = new FormData()
  formData.append("username", userData.username)
  formData.append("email", userData.email)
  formData.append("password", userData.password)
  formData.append("image_original", userData.image_original)
  formData.append("image_avatar", userData.image_avatar)
  formData.append("first_name", userData.first_name)
  formData.append("last_name", userData.last_name)
  formData.append("country", userData.country)
  formData.append("city", userData.city)
  formData.append("gender", userData.gender)

  return api.post('/api/reg-user', formData, {
    headers: {
      'Content-Type' : 'multipart/form-data'
    }
  });
};

export const sendVerificationCodeOnServer = (email) => {
    return api.post('/api/send-verif-code', {email})
};

export const VerifyCodeOnServer = (email, code) => {
    return api.post('/api/check-verif-code', {email, code})
}

export const  loginUser = (userData) => {
  return api.post('/api/login', userData)
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

    // Игнорируем запросы на refresh и logout при 401
    if (error.response?.status === 401 && 
        !originalRequest._retry && 
        !originalRequest.url.includes('/refresh-tokens') && 
        !originalRequest.url.includes('/logout')) {
      
      originalRequest._retry = true;

      try {
        const { data } = await refreshTokensByServer();
        localStorage.setItem('accessToken', data.accessToken);
        api.defaults.headers.common['Authorization'] = `Bearer ${data.accessToken}`;
        return api(originalRequest);
      } catch (refreshError) {
        // Полная очистка при неудачном обновлении
        localStorage.removeItem('accessToken');
        delete api.defaults.headers.common['Authorization'];
        
        const auth = useAuthStore();
        auth.openAuthModal();
        router.push('/');
        
        return Promise.reject(refreshError);
      }
    }
    
    // Для всех остальных 401 (включая /refresh-tokens и /logout)
    if (error.response?.status === 401) {
      const auth = useAuthStore()
      auth.clearAuth()
    }
    
    return Promise.reject(error);
  }
);

export default api;