import axios from 'axios';

const api = axios.create({
  baseURL: 'http://localhost:8080', 
  headers: {
    'Content-Type': 'application/json',
  },
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