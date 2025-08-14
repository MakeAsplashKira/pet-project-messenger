<template>
    <div class="header-wrapper">
        <div class="header">
        <div class="header-logo" @click="$router.push(`/`)">
            <div ref="lottieContainer" class="logo-animation"></div> 
            <div class="site-name">MediaVerse</div>
        </div>
        <MiniPlayer v-if="$route && !$route.path.includes('/music')"/>
        <div class="header-profile" @click="toggleProfileModal()">
          <img :src="avatarUrl" alt="">
        </div>
        <div class="header-profile-modal" ref="headerModal">
          <div class="modal-btn settings">
            <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 1024 1024"><path d="M512.5 390.6c-29.9 0-57.9 11.6-79.1 32.8c-21.1 21.2-32.8 49.2-32.8 79.1c0 29.9 11.7 57.9 32.8 79.1c21.2 21.1 49.2 32.8 79.1 32.8c29.9 0 57.9-11.7 79.1-32.8c21.1-21.2 32.8-49.2 32.8-79.1c0-29.9-11.7-57.9-32.8-79.1a110.96 110.96 0 0 0-79.1-32.8m412.3 235.5l-65.4-55.9c3.1-19 4.7-38.4 4.7-57.7s-1.6-38.8-4.7-57.7l65.4-55.9a32.03 32.03 0 0 0 9.3-35.2l-.9-2.6a442.5 442.5 0 0 0-79.6-137.7l-1.8-2.1a32.12 32.12 0 0 0-35.1-9.5l-81.2 28.9c-30-24.6-63.4-44-99.6-57.5l-15.7-84.9a32.05 32.05 0 0 0-25.8-25.7l-2.7-.5c-52-9.4-106.8-9.4-158.8 0l-2.7.5a32.05 32.05 0 0 0-25.8 25.7l-15.8 85.3a353.44 353.44 0 0 0-98.9 57.3l-81.8-29.1a32 32 0 0 0-35.1 9.5l-1.8 2.1a445.93 445.93 0 0 0-79.6 137.7l-.9 2.6c-4.5 12.5-.8 26.5 9.3 35.2l66.2 56.5c-3.1 18.8-4.6 38-4.6 57c0 19.2 1.5 38.4 4.6 57l-66 56.5a32.03 32.03 0 0 0-9.3 35.2l.9 2.6c18.1 50.3 44.8 96.8 79.6 137.7l1.8 2.1a32.12 32.12 0 0 0 35.1 9.5l81.8-29.1c29.8 24.5 63 43.9 98.9 57.3l15.8 85.3a32.05 32.05 0 0 0 25.8 25.7l2.7.5a448.27 448.27 0 0 0 158.8 0l2.7-.5a32.05 32.05 0 0 0 25.8-25.7l15.7-84.9c36.2-13.6 69.6-32.9 99.6-57.5l81.2 28.9a32 32 0 0 0 35.1-9.5l1.8-2.1c34.8-41.1 61.5-87.4 79.6-137.7l.9-2.6c4.3-12.4.6-26.3-9.5-35m-412.3 52.2c-97.1 0-175.8-78.7-175.8-175.8s78.7-175.8 175.8-175.8s175.8 78.7 175.8 175.8s-78.7 175.8-175.8 175.8"/></svg>
            <span>Настройки</span>
          </div>
          <div class="modal-btn logout" @click="logout">
            <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24"><path fill="none" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 4.001H5v14a2 2 0 0 0 2 2h8m1-5l3-3m0 0l-3-3m3 3H9"/></svg>
            <span>Выйти</span>
          </div>
        </div>

    </div>
  </div>

</template>

<script setup>
import { ref, onMounted, onUnmounted, computed } from 'vue';
import lottie from 'lottie-web';
import { useAuthStore } from '@/stores/auth';
import MiniPlayer from '../player/MiniPlayer.vue';

const lottieContainer = ref(null);
const headerModal = ref(null)
const headerModalOn = ref(false)

const auth = useAuthStore()

const avatarUrl = computed(()=> {
    return `http://localhost:8080/api/profile/${auth.userID}/avatar`
})

function toggleProfileModal() {
  if(headerModal.value) {
    if(!headerModalOn.value) {
      headerModal.value.style.transform = 'translateX(0%)'
      headerModal.value.style.opacity = '1'
      headerModal.value.style.pointerEvents = 'all'
      headerModalOn.value = true
    }
    else {
      headerModal.value.style.transform = 'translateX(200%)'
      headerModal.value.style.opacity = '0'
      headerModal.value.style.pointerEvents = 'none'
      headerModalOn.value = false
    }
  }
}

const handleKeyDown = (event) => {
  if (event.key === 'Escape') {
    toggleProfileModal()
  }
}

onMounted(async () => {
  window.addEventListener('keydown', handleKeyDown)
  const animationData = await import('@/assets/logo-animation.json');
  lottie.loadAnimation({
    container: lottieContainer.value,
    renderer: 'svg',
    animationData: animationData.default,
    loop: true,
  });
});

onUnmounted(()=> {
  window.removeEventListener('keydown', handleKeyDown)
})

const logout = async() => {
  toggleProfileModal()
  const data = await auth.logout()
}
</script>



<style scoped>

.header-wrapper {
  position: fixed;
  top: 20px;
  width: 100%;
  height: 80px;
  display: flex;
  background-color: #0e1621;
  justify-content: center;
}



.header {
  padding: 0 20px;
  display: flex;
  position: relative;
  justify-content: space-between;
  max-width: 1200px;
  min-width: 600px;
  width: 100%;
  height: 100%;
  background-color: #17212b;
  align-items: center;
  box-shadow: 0px 1px 8px black;
  border-radius: 8px;
  margin: 0 20px;
}

.header-logo {
  cursor: pointer;
  overflow: hidden;
}

.logo-animation {
  position: absolute;
  top: -90%;
  width: 300px; /* Размер логотипа */
  height: 300px;
  left: -90px;
  pointer-events: none;
}

.site-name {
  margin-left: 70px;
  font-size: 32px;
  user-select: none;
  font-family: "Lobster Two", sans-serif;
}

.header-profile {
  position: relative;
  width: 45px;
  height: 45px;
  border-radius: 50%;
  cursor: pointer;
  box-shadow: 2px 3px 6px black;
  background-color: #0e1621;
  z-index: 10;
}
.header-profile > img {
  width: 45px;
  height: 45px;
}

.header-profile::after {
  content: "";
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  width: 100%; 
  height: 100%;
  border-radius: 50%;
  background: transparent;
  box-shadow: 0 0 100px 10px rgba(1, 255, 255, 0.408); 
  opacity: 0; 
  transition: all 0.6s ease-in-out; 
  z-index: -1; 
}

.header-profile:hover::after {
  opacity: 1;
}


.header-profile-modal {
  position: absolute;
  width: 200px;
  height: 300px;
  background-color: #0e1621;
  top: 0;
  right: 0;
  z-index: 1;
  border-radius: 8px;
  box-shadow: 0px 2px 6px black;
  transition: all .3s ease;
  padding: 15px;
  transform: translateX(200%);
  opacity: 0;
}

.modal-btn {
  display: flex;
  align-items: center;
  justify-content: left;
  gap: 5px;
  background-color: #17212b;
  border-radius: 8px;
  cursor: pointer;
  padding: 8px 16px;
  box-shadow: 1px 1px 6px black;
  transition: all .3s ease-in-out;
}

.modal-btn:first-child {
  margin-top: 70px;
}

.settings:hover > svg {
  animation: rotating 4s linear infinite;
}

@keyframes rotating {
  0% {
    transform: rotate(0deg);
  }
  100% {
    transform: rotate(360deg);
  }
}

.modal-btn:not(:first-child) {
  margin-top: 10px;
}
.modal-btn:hover {
  transform: translateY(-2px);
  box-shadow: 1px 4px 6px black;
}

.modal-btn > svg {
  fill: rgba(255, 255, 255, .6);
  stroke: rgba(255, 255, 255, .6);
  transition: all .3s ease-in-out;
}
.modal-btn:hover > svg {
  fill: rgba(255, 255, 255, .9);
  stroke: rgba(255, 255, 255, .9);
}

.modal-btn > span {
  color: rgba(255,255,255,.6);
  transition: all .3s ease-in-out;
}
.modal-btn:hover > span {
  color: rgba(255,255,255,.9);
}


</style>