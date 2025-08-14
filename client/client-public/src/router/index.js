import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import Messages from '@/components/MainPage/Messages.vue'
import UserProfile from '@/components/MainPage/UserProfile.vue'
import Music from '@/components/Music.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: UserProfile
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
    {
      path: '/music',
      name: 'music',
      component : Music,
      meta: {requiresAuth: true}
    }
  ],
})
router.beforeEach((to, from, next) => {
  const authStore = useAuthStore()
  
  if (to.meta.requiresAuth && !authStore.isAuthenticated) {
    authStore.closeAuthModal()
    
    if(!authStore.authModalOpen) {
      authStore.openAuthModal()
    }

    return next(from.path === '/' ? false : from.path)
  }
  next()
})
export default router
