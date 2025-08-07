<template>
  <div class="play-btn-wrapper" :class="{ 'is-playing': player.isPlaying }" @click="togglePlay">
    <div class="play-btn-anim">
      <div class="playing-btn-anim a1"></div>
      <div class="playing-btn-anim a2"></div>
      <div class="playing-btn-anim a3"></div>
      <div class="playing-btn-anim a4"></div>
    </div>
    <div class="play-btn">
      <svg v-if="!player.isPlaying" xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 16 16"><path fill="#000" d="m11.596 8.697l-6.363 3.692c-.54.313-1.233-.066-1.233-.697V4.308c0-.63.692-1.01 1.233-.696l6.363 3.692a.802.802 0 0 1 0 1.393"/></svg>
      <svg v-else xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24"><path fill="#000" d="M17.276 5.47c.435.16.724.575.724 1.039V17.49c0 .464-.29.879-.724 1.039a3.69 3.69 0 0 1-2.552 0A1.107 1.107 0 0 1 14 17.491V6.51c0-.464.29-.879.724-1.04a3.69 3.69 0 0 1 2.552 0m-8 0c.435.16.724.575.724 1.039V17.49c0 .464-.29.879-.724 1.039a3.69 3.69 0 0 1-2.552 0A1.107 1.107 0 0 1 6 17.491V6.51c0-.464.29-.879.724-1.04a3.69 3.69 0 0 1 2.552 0"/></svg>
    </div>
  </div>
</template>

<script setup>
import { watch, nextTick } from 'vue'
import { usePlayerStore } from '@/stores/player'

const player = usePlayerStore()

// При первом клике инициализируем аудио, если еще не инициализировано
function togglePlay() {
  if (!player.audio) {
    // Вставь сюда нужный URL трека
    player.initAudio('http://localhost:8080/api/stream-track/7')
    player.play()
  } else {
    player.togglePlay()
  }
}

// Управляем анимацией кружков в зависимости от состояния isPlaying
function showCirclesSeq() {
  const delays = [0, 200, 400, 600]
  delays.forEach((delay, i) => {
    setTimeout(() => {
      const circle = document.querySelector(`.a${i + 1}`)
      if (circle) {
        circle.style.width = '44px'
        circle.style.height = '44px'
        circle.style.boxShadow = '3px 3px 6px black'
        circle.style.opacity = '1'
      }
    }, delay)
  })
}

function hideCircleSeq() {
  const delays = [0, 200, 400, 600]
  delays.forEach((delay, i) => {
    setTimeout(() => {
      const circle = document.querySelector(`.a${i + 1}`)
      if (circle) {
        circle.style.width = '10px'
        circle.style.height = '10px'
        circle.style.boxShadow = 'none'
        circle.style.opacity = '0'
      }
    }, delay)
  })
}

// Следим за состоянием проигрывания и запускаем/останавливаем анимацию
watch(
  () => player.isPlaying,
  (newVal) => {
    // Чтобы быть уверенным, что DOM уже обновился, используем nextTick
    nextTick(() => {
      if (newVal) {
        showCirclesSeq()
      } else {
        hideCircleSeq()
      }
    })
  },
  { immediate: true }
)
</script>

<style scoped>
/* твои стили без изменений */
.play-btn-wrapper {
  position: relative;
  width: 50px;
  height: 50px;
  cursor: pointer;
}

.play-btn {
  position: relative;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  background-color: rgba(255, 255, 255, 0.7);
  width: 40px;
  height: 40px;
  text-align: center;
  color: black;
  z-index: 5;
  transition: all 0.3s ease-in-out;
}

.play-btn > svg {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-45%, -50%);
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
  background-color: rgba(255, 255, 255, 0.7);
  opacity: 0;
  transform: scale(0.3);
  transition: all 0.5s ease-in-out;
}

/* Расположение кружков */
.a1 {
  top: 0;
  left: 50%;
  transform: translateX(-50%);
}
.a2 {
  right: 0;
  top: 50%;
  transform: translateY(-50%);
}
.a3 {
  bottom: 0;
  left: 50%;
  transform: translateX(-50%);
}
.a4 {
  left: 0;
  top: 50%;
  transform: translateY(-50%);
}

.a1,
.a2,
.a3,
.a4 {
  transition: all 0.3s ease-in-out;
}

.a1.shrink,
.a2.shrink,
.a3.shrink,
.a4.shrink {
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%) scale(0.3);
}

/* Анимация вращения */
@keyframes rotate-container {
  from {
    transform: translate(-50%, -50%) rotate(0deg);
  }
  to {
    transform: translate(-50%, -50%) rotate(360deg);
  }
}
</style>
