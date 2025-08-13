import { defineStore } from "pinia";

export const useRegAuth = defineStore('regAuth', {
    state: () => ({
        currentStep: 1,
        password: "",
        email: "",
        username: "",
        firstName: "",
        lastName: "",
        avatarImage: null,
        originalImage: null,
        country: "",
        city: "",
    }),
    actions: {
        setAvatarImage(image) {
            this.avatarImage = image
        },
        setOriginalImage(image) {
            this.originalImage = image
            this.avatarImage = null
        },
        setPassword(password) {
            this.password = password
        },
        setEmail(email) {
            this.email = email
        },
        setUsername(username) {
            this.username = username
        },
        setFirstName(firstName) {
            this.firstName = firstName
        },
        setLastName(lastName) {
            this.lastName = lastName
        },
        setCountry(country) {
            this.country = country
        },
        setCity(city) {
            this.city = city
        }
    } 
})