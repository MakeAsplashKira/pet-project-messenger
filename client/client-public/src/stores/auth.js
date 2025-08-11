import { defineStore } from 'pinia';
import { ref, computed } from 'vue';
import api, { logoutUser, registerUser, loginUser, refreshTokensByServer} from '@/api';
import useNotify from '@/composable/useNotify';


export const useAuthStore = defineStore('auth', () => {
  const {notify} = useNotify()
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
  try {
    accessToken.value = null;
    user.value = null;
    localStorage.removeItem('accessToken');
    if (api?.defaults?.headers?.common?.Authorization) {
      delete api.defaults.headers.common['Authorization'];
    }
  } catch (error) {
    console.error('clearAuth error:', error);
  }
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
  const login = async (userData) => {
    try {
      const { data } = await loginUser(userData);
      setAuth(data.access_token);
      user.value = data.user;
      closeAuthModal();
      notify.success('Авторизация', 'Вы успешно вошли в аккаунт')
      return { success: true };
    } catch (error) {
      notify.error('Авторизация', error.response?.data?.error || 'Ошибка авторизации')
    }
  };
  
  const logout = async () => {
    try {
      const {data} = await logoutUser();
      if(data?.status == 'success') {
        console.log(data)
        clearAuth();
        notify.success('Успех', 'Вы усешно вышли из аккаунта')
      }
    }
    catch(error) {
      notify.error('Ошибка', 'Не удалось выйти из аккаунта')
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
        if(!isTokenExpired(accessToken.value)) {
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