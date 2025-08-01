<script setup>
import { ref, onMounted } from 'vue';
import lottie from 'lottie-web';
import { logoutUser } from '@/api';

const lottieContainer = ref(null);

onMounted(async () => {
  const animationData = await import('@/assets/logo-animation.json');
  lottie.loadAnimation({
    container: lottieContainer.value,
    renderer: 'svg',
    animationData: animationData.default,
    loop: true,
  });
});
</script>

<template>
    <div class="header-wrapper">
        <div class="header">
        <div class="header-logo" @click="$router.push(`/`)">
            <div ref="lottieContainer" class="logo-animation"></div> 
            <div class="site-name">MediaVerse</div>
        </div>
        <div class="header-main">тут что-то будет</div>
        <div @click="logoutUser" class="vihjod" style="cursor: pointer;">Выход</div>
        <div class="header-profile" @click="$router.push('/profile')"></div>
    </div>
  </div>

</template>

<style scoped>

.header-wrapper {
  position: fixed;
  top: 0%;
  width: 100%;
  height: 80px;
  background-color: #17212b;
  display: flex;
  justify-content: center;
  -webkit-box-shadow: 0px 6px 20px -3px rgb(0, 0, 0);
  -moz-box-shadow: 0px 6px 20px -3px rgb(0, 0, 0);
  box-shadow: 0 0 20px rgba(0, 0, 0, 0.25);
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
  align-items: center;
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
  width: 45px;
  height: 45px;
  border-radius: 50%;
  cursor: pointer;
  box-shadow: 1px 2px 5px black;
  background-color: #0e1621;
  position:relative; 
  z-index: 1;
}

.header-profile::after {
  content: "";
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  width: 100%; /* Увеличиваем размер свечения */
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



</style>