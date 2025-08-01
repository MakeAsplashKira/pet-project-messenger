import './assets/main.css'

import { createApp } from 'vue'
import { createPinia } from 'pinia'
import piniaPluginPersistedstate from 'pinia-plugin-persistedstate'

import App from './App.vue'
import router from './router'
import { useAuthStore } from './stores/auth'

const app = createApp(App)

//Инициализация хранилища pinia
const pinia = createPinia()
pinia.use(piniaPluginPersistedstate)
app.use(pinia)

const authStore = useAuthStore()
authStore.checkAuth().then(authenticated => {
    if (!authenticated && router.currentRoute.value.meta.requiresAuth) {
        console.log('Надо тренероватся бро..')
    }
})

app.use(router)

app.mount('#app')
