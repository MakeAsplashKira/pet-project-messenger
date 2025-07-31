import { defineStore } from "pinia"

export const useAuthStore = defineStore('auth', {
    state:() => ({
        isAuth: false,
        authModalOpen: false,
    }),
    actions: {
        openAuthModal() {
            this.authModalOpen = true
        },
        nextStep(){
            this.curAuthStep++
        }
    }
})