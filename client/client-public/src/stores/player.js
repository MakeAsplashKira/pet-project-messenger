import { defineStore } from 'pinia'
import { ref, reactive } from 'vue'

export const usePlayerStore = defineStore('player', () => {
  const audio = ref(null)
  const isPlaying = ref(false)
  const currentTime = ref(0)
  const duration = ref(0)
  const bufferedPercent = ref(0)

  function initAudio(src) {
    if (audio.value) {
      audio.value.pause()
      audio.value.src = ''
      removeListeners()
    }
    audio.value = new Audio(src)
    audio.value.addEventListener('timeupdate', onTimeUpdate)
    audio.value.addEventListener('loadedmetadata', onLoadedMetadata)
    audio.value.addEventListener('ended', onEnded)
    audio.value.addEventListener('progress', onProgress)
  }

  function onTimeUpdate() {
    currentTime.value = audio.value.currentTime
  }

  function onLoadedMetadata() {
    duration.value = audio.value.duration
  }

  function onEnded() {
    isPlaying.value = false
  }

  function onProgress() {
    if(!audio.value) return
    const buffered = audio.value.buffered
    const dur = audio.value.duration
    if (buffered.length && dur > 0) {
        const bufferedEnd = buffered.end(buffered.length-1)
        bufferedPercent.value = Math.min(100, (bufferedEnd / dur) * 100)
    }
  }

  function play() {
    if (!audio.value) return
    audio.value.play()
    isPlaying.value = true
  }

  function pause() {
    if (!audio.value) return
    audio.value.pause()
    isPlaying.value = false
  }

  function togglePlay() {
    if (!audio.value) return
    if (isPlaying.value) pause()
    else play()
  }

  function seek(time) {
    if (!audio.value) return
    if (time < 0) time = 0
    if (time > duration.value) time = duration.value
    audio.value.currentTime = time
  }

  function removeListeners() {
    if (!audio.value) return
    audio.value.removeEventListener('timeupdate', onTimeUpdate)
    audio.value.removeEventListener('loadedmetadata', onLoadedMetadata)
    audio.value.removeEventListener('ended', onEnded)
  }


  return {
    audio,
    isPlaying,
    currentTime,
    duration,
    bufferedPercent,
    initAudio,
    play,
    pause,
    togglePlay,
    seek,
  }
})
