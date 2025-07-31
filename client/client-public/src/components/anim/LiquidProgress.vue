<template>
  
    <div class="container">
      <div class="circle" :style="circleStyle">
        <div class="bg-field" :style="bgFieldStyle">
          <div class="wave-whitespace-1"></div>
        </div>
        <p :style="textStyle" class="text" id="anim-auth-text">{{ text }}</p>
      </div>
    </div>
</template>

<script setup>
import { computed, onMounted, watch } from 'vue'
const props = defineProps({
    currentStage: {
      type: Boolean,
      default: false
    },
    fillPercent: {
        type: Number,
        default: 40
    },
    cSize: {
        type: Number,
        default: 300
    },
    text: {
        type: String,
        default: "1"
    }
})



const updateFill = (value) => {
  document.documentElement.style.setProperty('--fill', value)
}
const updateCSize =()=> {
    document.documentElement.style.setProperty('--circle-size', props.cSize + 'px')
}



watch (()=> props.fillPercent, (value) => {
    updateFill(value)
    //стилизация
})

const circleStyle = computed(()=> ({
  width: `${props.cSize}px`,
  height: `${props.cSize}px`
}))

const bgFieldStyle = computed(()=>({
  bottom: `${props.cSize * props.fillPercent / 100}px`,
  height: `${props.cSize * 2}px`,
  width: `${props.cSize * 2}px`
}))

const textStyle = computed(()=>({
  color: props.currentStage ? '#00ff0d' : '#fff',
  fontSize: props.text.length > 1 ? '22px' : '40px'
}))

onMounted(()=> {
  updateCSize()
  updateFill(props.fillPercent)
})



</script>

<style scoped>

.text {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  font-size: 40px;
  font-weight: bold;
  text-align: center;
  margin: 0;
  mix-blend-mode: difference;
  z-index: 2;
  transition: all .1s ease-in-out;
}

.circle {
  position: relative;
  clip-path: circle();
  background: linear-gradient(125deg, rgb(0, 255, 21) 0%, rgb(3, 255, 37) 100%);
  box-shadow: inset 0 0 40px 20px rgba(0, 0, 0, 0.2);
  display: grid;
  place-items: center;
}

.bg-field, .bg-field-2 {
  position: absolute;
  transition: bottom 350ms;
}

.bg-field-2::before, .bg-field-2::after {
  content: "";
  position: absolute;
  inset: 0;
  background-color: white;
  clip-path: inset(0 0 0 0 round 45%);
  animation: spin linear infinite;
  animation-duration: 5s;
}

.bg-field-2::before, .bg-field-2::after {
  clip-path: inset(0 0 0 0 round 44%);
  background-color: rgba(255, 255, 255, 0.2);
  transform: translateX(50px);
  animation-duration: 8s;
}

.bg-field::before, .bg-field-2::before {
  z-index: -1;
}

.bg-field::after, .bg-field-2::after {
  z-index: 1;
}

.wave-whitespace-1, .wave-whitespace-2 {
  position: relative;
  width: 100%;
  height: 100%;
  background: #17212b;
  clip-path: inset(0 0 0 0 round 45%);
  animation: spin linear infinite;
  animation-duration: 5s;
}

.wave-whitespace-2 {
  position: absolute;
  background: none;
  animation: spin linear infinite;
  animation-duration: 8s;
}

@keyframes spin {
  0% {
    transform: rotate(0deg);
  }
  100% {
    transform: rotate(360deg);
  }
}
</style>