<script setup>
import { usePlayerStore } from '@/stores/player'
import { computed } from 'vue'

const player = usePlayerStore()

const progressPercent = computed(() => {
  if (player.duration === 0) return 0
  return (player.currentTime / player.duration) * 100
})

function onClickProgressBar(event) {
  const rect = event.currentTarget.getBoundingClientRect()
  const clickX = event.clientX - rect.left
  const seekTime = (clickX / rect.width) * player.duration
  player.seek(seekTime)
}
</script>

<template>
  <div class="player-progress-bar-wrapper" @click="onClickProgressBar">
    <div class="player-progress-bar-background">
      <div class="player-progress-buffered-bar" :style="{width: player.bufferedPercent + '%'}"></div>
      <div class="player-progress-bar-foreground" :style="{ width: progressPercent + '%' }"></div>
    </div>
  </div>
</template>

<style scoped>
.player-progress-bar-wrapper {
  padding: 3px 0;
  cursor: pointer;
  position: relative;
}

.player-progress-bar-background {
  margin-top: 10px;
  width: 100%;
  height: 4px;
  background-color: rgba(255, 255, 255, 0.356);
  border-radius: 10px;
  position: relative;
  overflow: hidden;
}

.player-progress-bar-foreground {
  height: 100%;
  position: relative;
  background-color: white;
  transition: width 0.1s linear;
  border-radius: 10px;
  z-index: 3;
}

.player-progress-buffered-bar {
    height: 100%;
    position: absolute;
    left: 0;
    height: 100%;
    transition: width .1s linear;
    border-radius: 10px;
    background-color:rgb(160, 160, 160);
    z-index: 2;
}
</style>
