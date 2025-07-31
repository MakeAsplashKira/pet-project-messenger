import { createRouter, createWebHistory } from 'vue-router';
import Users from '@/views/Users.vue';

const routes = [
  { path: '/admin/users', component: Users },
  { path: '/pathMatch(.*)', redirect: '/admin/' }, // Перенаправление на /users
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

export default router;