import { defineStore } from "pinia";

export const useRegAuth = defineStore('regAuth', {
    state: () => ({
        currentStep: 1,
        password: "",
        email: "",
        username: "",
        avatarImage: null,
    }),
    actions: {
        setAvatarImage(image) {
            this.avatarImage = image
        },
        setPassword(password) {
            this.password = password
        },
        setEmail(email) {
            this.email = email
        },
        setUsername(username) {
            this.username = username
        }
    } 
})