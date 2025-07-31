import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import Messages from '@/components/MainPage/Messages.vue'
import UserProfile from '@/components/MainPage/UserProfile.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
    },
    {
      path: '/profile',
      name: 'profile',
      component: UserProfile,
      meta: {requiresAuth: true}
      
    },
    {
      path: '/messages',
      name: 'messages',
      component: Messages,
      meta: {requiresAuth: true}
    },
  ],
})
router.beforeEach((to, from, next) => {
  const authStore = useAuthStore()
  
  if (to.meta.requiresAuth && !authStore.isAuth) {
    authStore.openAuthModal() // Открываем модалку
    next(true) // Отменяем текущую навигацию
  } else {
    next() // Продолжаем навигацию
  }
})
export default router
