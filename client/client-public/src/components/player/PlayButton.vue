<template>
  <div class="play-btn-wrapper" :class="{ 'is-playing': isPlaying}" @click="togglePlay">
    <div class="play-btn-anim">
      <div class="playing-btn-anim a1"></div>
      <div class="playing-btn-anim a2"></div>
      <div class="playing-btn-anim a3"></div>
      <div class="playing-btn-anim a4"></div>
    </div>
    <div class="play-btn" v-if="!isPlaying">▶</div>
    <div class="play-btn" v-if="isPlaying">||</div>
  </div>
</template>

<script setup>
import { ref } from 'vue'

const isPlaying = ref(false)

function togglePlay() {
  if (isPlaying.value) {
    isPlaying.value = false
    hideCircleSeq()
  } else {
    isPlaying.value = true
    showCirclesSeq()
  }
}

function showCirclesSeq() {
    const delays = [0, 200, 400, 600]
    delays.forEach((delay, i) => {
        setTimeout(()=> {
            const circle = document.querySelector(`.a${i+1}`)
            if (circle) {
                circle.style.width = '44px'
                circle.style.height = '44px'
                circle.style.boxShadow = '3px 3px 6px black'
                circle.style.opacity = '1'
            }
        },delay)
    })
}

function hideCircleSeq() {
    const delays = [0, 200, 400, 600]
    delays.forEach((delay, i) => {
        setTimeout(()=> {
            const circle = document.querySelector(`.a${i+1}`)
            if(circle) {
                circle.style.width = '10px'
                circle.style.height = '10px'
                circle.style.boxShadow = 'none'
                circle.style.opacity = '0'
            }
        }, delay)
    })
}
</script>


<style scoped>
.play-btn-wrapper {
  position: relative;
  width: 50px;
  height: 50px;
  cursor: pointer;
}

.play-btn {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  background-color: rgba(255, 255, 255, .7);
  width: 40px;
  height: 40px;
  text-align: center;
  color: black;
  z-index: 5;
  transition: all .3s ease-in-out;
}

.play-btn:hover {
    transform: translate(-50%, -50%) scale(1.05);
   background-color: rgba(255, 255, 255, 1);
    
}

.play-btn-anim {
  position: absolute;
  top: 50%;
  left: 50%;
  width: 40px;
  height: 40px;
  transform: translate(-50%, -50%);
  z-index: 1;
  animation: rotate-container 4s linear infinite;
  pointer-events: none;
}

.playing-btn-anim {
  position: absolute;
  width: 10px;
  height: 10px;
  box-shadow: 3px 3px 6px black;
  border-radius: 50%;
  background-color: rgb(255, 255, 255, .7);
  opacity: 0;
  transform: scale(0.3);
  transition: all .5s ease-in-out;
}

/* Расположение кружков */
.a1 { top: 0; left: 50%; transform: translateX(-50%); }
.a2 { right: 0; top: 50%; transform: translateY(-50%); }
.a3 { bottom: 0; left: 50%; transform: translateX(-50%); }
.a4 { left: 0; top: 50%; transform: translateY(-50%); }

.a1, .a2, .a3, .a4 {
    transition: all .3s ease-in-out
}

.a1.shrink, .a2.shrink, .a3.shrink .a4.shrink {
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%) scale(.3);
}


/* Анимация вращения */
@keyframes rotate-container {
  from { transform: translate(-50%, -50%) rotate(0deg); }
  to   { transform: translate(-50%, -50%) rotate(360deg); }
}


</style>
