<script setup>
import { watch } from 'vue';
import { useAuthStore } from './stores/auth';
import MainPageHeader from './components/MainPage/MainPageHeader.vue';
import AuthModal from './components/AuthModal.vue';
import Notifications from './components/Notifications.vue';

const auth = useAuthStore()

watch(()=>auth.isAuthenticated, (isAuthenticated) => {
  if(!isAuthenticated) {
    auth.openAuthModal()
  }
}, {immediate: true})

</script>

<template>
  <Notifications/>
  <MainPageHeader/>
  <AuthModal v-if="auth.authModalOpen" />
  <div class="content-wrapper">
      <div class="content">
        <div class="content-menu">
          <div class="content-menu-row" @click="$router.push('/profile')">
            <div class="content-menu-row-icon"></div>
            <div class="content-menu-row-text">Профиль</div>
          </div>
          <div class="content-menu-row" @click="$router.push('/messages')">
            <div class="content-menu-row-icon"></div>
            <div class="content-menu-row-text">Сообщения</div>
          </div>
          <div class="content-menu-row" @click="$router.push('/music')">
            <div class="content-menu-row-icon"></div>
            <div class="content-menu-row-text">Музыка</div>
          </div>
        </div>
        <RouterView />
      </div>
  </div>
  
</template>

<style scoped>
.content-wrapper {
  display: flex;
  justify-content: center;
  width: 100%;
  min-height: calc(100vh - 100px);
  margin-top: 100px;
  background-color: #0e1621;
  padding: 0 20px;
  padding-top: 20px;
  padding-bottom: 20px;
}

.content {
  display: flex;
  max-width: 1200px;
  width: 100%;
  height: 100%;
}

.content-menu {
  width: 200px;
}

.content-menu-row{
  padding: 10px;
  border-radius: 4px;
  transition: all .4s ease-in-out;
  cursor: pointer;
 
}

.content-menu-row:hover {
  background-color: #17212b;
}

.content-menu-row-text {
   color: #ffffff88;
   transition: all .4s ease-in-out;
}

.content-menu-row:hover > .content-menu-row-text {
  color: #ffffff;
}


</style>
