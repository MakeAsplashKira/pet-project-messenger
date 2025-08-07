<template>
    <div class="header-wrapper">
        <div class="header">
        <div class="header-logo" @click="$router.push(`/`)">
            <div ref="lottieContainer" class="logo-animation"></div> 
            <div class="site-name">MediaVerse</div>
        </div>
        <MiniPlayer/>
        <!-- <div @click="logout" class="vihjod" style="cursor: pointer;">Выход</div> -->
        <div class="header-profile" @click="toggleProfileModal()">
        </div>
        <div class="header-profile-modal" ref="headerModal">
          <div class="logout-btn">
            <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24"><path fill="none" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 4.001H5v14a2 2 0 0 0 2 2h8m1-5l3-3m0 0l-3-3m3 3H9"/></svg>
            <span>Выйти</span>
          </div>
        </div>

    </div>
  </div>

</template>

<script setup>
import { ref, onMounted } from 'vue';
import lottie from 'lottie-web';
import { useAuthStore } from '@/stores/auth';
import MiniPlayer from '../player/MiniPlayer.vue';

const lottieContainer = ref(null);
const headerModal = ref(null)
const headerModalOn = ref(false)

const auth = useAuthStore()



function toggleProfileModal(event) {
  if(headerModal.value) {
    if(!headerModalOn.value) {
      headerModal.value.style.transform = 'translateX(0%)'
      headerModalOn.value = true
    }
    else {
      headerModal.value.style.transform = 'translateX(200%)'
      headerModalOn.value = false
    }
  }
  console.log(event)
}

onMounted(async () => {
  const animationData = await import('@/assets/logo-animation.json');
  lottie.loadAnimation({
    container: lottieContainer.value,
    renderer: 'svg',
    animationData: animationData.default,
    loop: true,
  });
});

const logout = async() => {
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
  box-shadow: 1px 2px 5px black;
  background-color: #0e1621;
  z-index: 10;
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
}

.logout-btn {
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
.logout-btn:hover {
  transform: translateY(-2px);
  box-shadow: 1px 4px 6px black;
}

.logout-btn > svg {
  stroke: rgba(255, 255, 255, .6);
  transition: all .3s ease-in-out;
}
.logout-btn:hover > svg {
  stroke: rgba(255, 255, 255, .9);
}

.logout-btn > span {
  color: rgba(255,255,255,.6);
  transition: all .3s ease-in-out;
}
.logout-btn:hover > span {
  color: rgba(255,255,255,.9);
}


</style>