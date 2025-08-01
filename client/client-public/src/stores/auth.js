import { defineStore } from 'pinia';
import { ref, computed } from 'vue';
import { useRouter } from 'vue-router';
import api, { logoutUser, registerUser, loginUser, refreshTokensByServer} from '@/api';


export const useAuthStore = defineStore('auth', () => {
  const router = useRouter();
  
  // State
  const accessToken = ref(null);
  const user = ref(null);
  const authModalOpen = ref(false);
  const authStep = ref(1); // Для многошаговой формы
  
  // Getters
  const isAuthenticated = computed(() => !!accessToken.value);
  
  // Actions
  const setAuth = (token) => {
    accessToken.value = token;
    localStorage.setItem('accessToken', token);
    api.defaults.headers.common['Authorization'] = `Bearer ${token}`;
  };
  
  const clearAuth = () => {
    accessToken.value = null;
    user.value = null;
    localStorage.removeItem('accessToken');
    delete api.defaults.headers.common['Authorization'];
  };
  
  const openAuthModal = (step = 1) => {
    authModalOpen.value = true;
    authStep.value = step;
  };
  
  const closeAuthModal = () => {
    authModalOpen.value = false;
  };
  
  const nextStep = () => {
    authStep.value++;
  };
  
  const prevStep = () => {
    authStep.value--;
  };
  
const register = async (userData) => {
  try {
    const {data} = await registerUser(userData);
    
    if (data?.status === 'success') {
        setAuth(data.access_token)
        return {success: true}
    }
    return {success: false, error: data?.error}

  } catch (error) {
    return {
        success: false,
        error: error.response?.data?.error || error.message
    }
    }
};
  const login = async (credentials) => {
    try {
      const { data } = await loginUser(credentials);
      setAuth(data.accessToken);
      user.value = data.user;
      closeAuthModal();
      return { success: true };
    } catch (error) {
      return {
        success: false,
        error: error.response?.data?.message || 'Login failed'
      };
    }
  };
  
  const logout = async () => {
    try {
      await logoutUser();
    } finally {
      clearAuth();
      router.push('/');
    }
  };

  const refreshTokens = async () => {
    try {
        const {data} = await refreshTokensByServer()
        setAuth(data.access_token)
        closeAuthModal()
        return true
    }
    catch(error) {
        clearAuth()
        return false
    }
  }
  
  const isTokenExpired = (token) => {
    const payload = JSON.parse(atob(token.split('.')[1]))
    return payload.exp * 1000 < Date.now()
  }

  const checkAuth = async () => {
   if (accessToken.value) {
        if(!isTokenExpired) {
            return true
        }
   }
   return await refreshTokens()
  };

  return {
    accessToken,
    user,
    authModalOpen,
    authStep,
    isAuthenticated,
    refreshTokens,
    clearAuth,
    checkAuth,
    openAuthModal,
    closeAuthModal,
    nextStep,
    prevStep,
    register,
    login,
    logout,
  };
}, {
    persist: {
        key: 'auth-store',
        paths: ['accessToken', 'user'],
        storage: localStorage,
        beforeHydrate: (context) => {
            //логика
        },
        afterHydrate: (context) => {
            if (context.store.accessToken) {
                api.defaults.headers.common['Authorization'] = `Bearer ${context.store.accessToken}`
            }
        }
    }
});