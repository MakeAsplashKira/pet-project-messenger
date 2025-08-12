<template>
  <div class="uploader-step">
    <div class="canvas-wrap"
    @dragover.prevent="onDragOver"
    @dragleave.prevent="onDragLeave"
    @drop.prevent="onDrop"
    :class="{'drag-active' : isDragging}"
    >
      <div class="label-photo">
        <AvatarUploadSVG :class="{'label-photo-hidden':isDragging}"/>
        <UploadAvatarOnDrag :class="{'label-photo-hidden':!isDragging}"/>
      </div>
      <div class="controls">
        <div @click="zoomIn()" class="zoom-btn in">+</div>
        <div @click="zoomOut()" class="zoom-btn out">-</div>
      </div>
      <canvas
        ref="canvasRef"
        width="500"
        height="500"
        @pointerdown="onPointerDown"
        @pointermove="onPointerMove"
        @pointerup="onPointerUp"
        @pointercancel="onPointerUp"
        @wheel.prevent="onWheel"
      ></canvas>
    </div>
    <div class="input-wrapper">
      <input id="fileInput" class="input-file" type="file" accept="image/*" @change="onFileChange" />
      <label for="fileInput">
        <svg  xmlns="http://www.w3.org/2000/svg" width="28" height="28" viewBox="0 0 1024 1024"><path data-v-aca44f06="" fill="#d4d4d4" d="M779.3 196.6c-94.2-94.2-247.6-94.2-341.7 0l-261 260.8c-1.7 1.7-2.6 4-2.6 6.4s.9 4.7 2.6 6.4l36.9 36.9a9 9 0 0 0 12.7 0l261-260.8c32.4-32.4 75.5-50.2 121.3-50.2s88.9 17.8 121.2 50.2c32.4 32.4 50.2 75.5 50.2 121.2c0 45.8-17.8 88.8-50.2 121.2l-266 265.9l-43.1 43.1c-40.3 40.3-105.8 40.3-146.1 0c-19.5-19.5-30.2-45.4-30.2-73s10.7-53.5 30.2-73l263.9-263.8c6.7-6.6 15.5-10.3 24.9-10.3h.1c9.4 0 18.1 3.7 24.7 10.3c6.7 6.7 10.3 15.5 10.3 24.9c0 9.3-3.7 18.1-10.3 24.7L372.4 653c-1.7 1.7-2.6 4-2.6 6.4s.9 4.7 2.6 6.4l36.9 36.9a9 9 0 0 0 12.7 0l215.6-215.6c19.9-19.9 30.8-46.3 30.8-74.4s-11-54.6-30.8-74.4c-41.1-41.1-107.9-41-149 0L463 364L224.8 602.1A172.22 172.22 0 0 0 174 724.8c0 46.3 18.1 89.8 50.8 122.5c33.9 33.8 78.3 50.7 122.7 50.7c44.4 0 88.8-16.9 122.6-50.7l309.2-309C824.8 492.7 850 432 850 367.5c.1-64.6-25.1-125.3-70.7-170.9"></path></svg>
      </label>
      <div class="save-btn" @click="saveAvatar()" :class="{'active-save-btn':saveState.allowed}">
        <svg v-if="!saveState.saved && !saveState.proccessing" xmlns="http://www.w3.org/2000/svg" width="28" height="28" viewBox="0 0 1024 1024"><path d="M893.3 293.3L730.7 130.7c-12-12-28.3-18.7-45.3-18.7H144c-17.7 0-32 14.3-32 32v736c0 17.7 14.3 32 32 32h736c17.7 0 32-14.3 32-32V338.5c0-17-6.7-33.2-18.7-45.2M384 176h256v112H384zm128 554c-79.5 0-144-64.5-144-144s64.5-144 144-144s144 64.5 144 144s-64.5 144-144 144m0-224c-44.2 0-80 35.8-80 80s35.8 80 80 80s80-35.8 80-80s-35.8-80-80-80"/></svg>
        <svg v-if="saveState.proccessing" xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24"><path fill="#ffffff" d="M12 2A10 10 0 1 0 22 12A10 10 0 0 0 12 2Zm0 18a8 8 0 1 1 8-8A8 8 0 0 1 12 20Z" opacity=".5"/><path fill="#ffffff" d="M20 12h2A10 10 0 0 0 12 2V4A8 8 0 0 1 20 12Z"><animateTransform attributeName="transform" dur="1s" from="0 12 12" repeatCount="indefinite" to="360 12 12" type="rotate"/></path></svg>
        <svg v-if="saveState.saved && !saveState.proccessing" xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 16 16"><path fill="#35fb32" d="M16 8A8 8 0 1 1 0 8a8 8 0 0 1 16 0m-3.97-3.03a.75.75 0 0 0-1.08.022L7.477 9.417L5.384 7.323a.75.75 0 0 0-1.06 1.06L6.97 11.03a.75.75 0 0 0 1.079-.02l3.992-4.99a.75.75 0 0 0-.01-1.05z"/></svg>
      </div>
    </div>
  </div>
</template>

<script setup>
import { watch, onUnmounted, reactive, ref } from 'vue'
import AvatarUploadSVG from '../../public/AvatarUploadSVG.vue'
import UploadAvatarOnDrag from '../../public/UploadAvatarOnDrag.vue'
import useNotify from '@/composable/useNotify'
import { useRegAuth } from '@/stores/regAuth'
const {notify} = useNotify()
const regAuth = useRegAuth()

let MIN_SCALE = 0.5
const MAX_SCALE = 2
const canvasRef = ref(null)
const img = new Image()
let imgLoaded = false
const isDragging = ref(false)
const saveState = reactive({
  allowed: false,
  saved: false,
  proccessing: false
})
const state = {
  x: 0,
  y: 0,
  scale: 1,
  imgWidth: 0,
  imgHeight: 0
}

function onDragOver(e) {
  e.preventDefault()
  isDragging.value = true
}

function onDragLeave(e) {
  e.preventDefault()
  isDragging.value = false
}

function onDrop(e) {
  e.preventDefault()
  isDragging.value = false
  const file = e.dataTransfer.files[0];
  if (!file || !file.type.match('image.*')){
    notify.warning('Аватар', 'Пожалуйста, вставтье изображение с соотвествующим расширением.')
    return
  }
  const event = {target: {files: [file]}}
  onFileChange(event)
}

function clampPosition() {
 const canvas = canvasRef.value;
  const imgWidth = img.width * state.scale;
  const imgHeight = img.height * state.scale;
  const centerX = canvas.width / 2;
  const centerY = canvas.height / 2;
  const radius = canvas.width / 2.5; // радиус видимой области

  // Горизонтальная ось
  if (imgWidth <= radius * 2) {
    // если картинка уже круга — центрируем
    state.x = centerX - imgWidth / 2;
  } else {
    const minX = centerX - radius;
    const maxX = centerX + radius - imgWidth;
    state.x = Math.min(minX, Math.max(maxX, state.x));
  }

  // Вертикальная ось
  if (imgHeight <= radius * 2) {
    state.y = centerY - imgHeight / 2;
  } else {
    const minY = centerY - radius;
    const maxY = centerY + radius - imgHeight;
    state.y = Math.min(minY, Math.max(maxY, state.y));
  }
}

let dragging = false
const lastPointer = { x: 0, y: 0 }

async function onFileChange(e) {
  try {
    // 1. Получаем файл
    const file = e.target.files?.[0] || e.dataTransfer?.files?.[0];
    if (!file || !file.type.match('image.*')) {
      notify.warning('Аватар', 'Пожалуйста выберите фотографию с нужным типом файла.')
      return;
    }

    // 2. Сброс состояний
    saveState.allowed = true
    saveState.saved = false
    regAuth.setOriginalImage(null)
    imgLoaded = false

    // 3. Сохраняем оригинал (blob)
    regAuth.setOriginalImage(file)

    // 4. Преобразуем файл в Data URL
    const imageData = await blobToDataUrl(file)

    // 5. Загружаем в canvas
    await loadImageToCanvas(imageData)
  } catch(error) {
    notify.error('Аватар', error.message)
  }
}
function blobToDataUrl(blob) {
  return new Promise((resolve, reject) => {
    const reader = new FileReader();
    reader.onload = e => resolve(e.target.result);
    reader.onerror = () => reject(new Error('Ошибка чтения файла'));
    reader.readAsDataURL(blob);
  });
}
async function loadImageToCanvas(imageData) {
  return new Promise((resolve, reject) => {
    img.src = imageData
    img.onload = () => {
      try {
        // Настройки canvas
        prepareCanvas()
        const canvas = canvasRef.value;
        const radius = canvas.width / 2.5;
        
        // Рассчет масштаба
        const coverScale = Math.max(
          (radius * 2) / img.width, 
          (radius * 2) / img.height
        );
        
        state.scale = coverScale;
        MIN_SCALE = Math.min(
          (radius * 2) / img.width,
          (radius * 2) / img.height
        );
        
        // Центрирование
        state.x = (canvas.width - img.width * state.scale) / 2;
        state.y = (canvas.height - img.height * state.scale) / 2;
        
        // Сохраняем размеры
        state.imgWidth = img.width;
        state.imgHeight = img.height;
        
        // Рисуем
        imgLoaded = true;
        draw();
        resolve();
      } catch(error) {
        reject(new Error(`Ошибка инициализации: ${error.message}`))
      }
    }
    img.onerror = () => reject(new Error('Невозможно загрузить изображение'))
  })
}



function prepareCanvas() {
  const svg = document.querySelector('.label-photo')
  const controls = document.querySelector('.controls')
  controls.classList.add('controls-visible')
  canvasRef.value.style.borderRadius = '8px'
  svg.style.opacity = 0
}

function drawImageWithMask(ctx, canvas) {
  ctx.clearRect(0, 0, canvas.width, canvas.height);
  
  ctx.drawImage(
    img,
    state.x, state.y,
    img.width * state.scale,
    img.height * state.scale
  );
  
  ctx.save();
  
  ctx.beginPath();
  ctx.arc(
    canvas.width / 2,
    canvas.height / 2,
    canvas.width / 2.5,
    0,
    Math.PI * 2
  );
  ctx.clip(); 
  
  ctx.drawImage(
    img,
    state.x, state.y,
    img.width * state.scale,
    img.height * state.scale
  );
  
  ctx.restore();
  
  ctx.fillStyle = 'rgba(0, 0, 0, 0.5)';
  ctx.beginPath();
  ctx.rect(0, 0, canvas.width, canvas.height);
  ctx.arc(
    canvas.width / 2,
    canvas.height / 2,
    canvas.width / 2.5,
    0,
    Math.PI * 2
  );
  ctx.fill('evenodd'); // Вырезаем круг из затемнения
  
  ctx.beginPath();
  ctx.arc(
    canvas.width / 2,
    canvas.height / 2,
    canvas.width / 2.5,
    0,
    Math.PI * 2
  );
  ctx.strokeStyle = 'white';
  ctx.lineWidth = 4;
  ctx.stroke();
}




function draw() {
  const canvas = canvasRef.value
  if (!canvas) return
  const ctx = canvas.getContext('2d')

  if (!imgLoaded) {
    ctx.clearRect(0, 0, canvas.width, canvas.height)
    ctx.fillStyle = '#0e1621'
    ctx.fillRect(0, 0, canvas.width, canvas.height)
    return
  }

  ctx.save()
  drawImageWithMask(ctx, canvas)
  ctx.restore()
}


function onPointerDown(e) {
  e.preventDefault()
  const canvas = canvasRef.value
  canvas.setPointerCapture(e.pointerId)
  dragging = true
  lastPointer.x = e.clientX
  lastPointer.y = e.clientY
}

function onPointerMove(e) {
  if (!dragging) return;
  saveState.saved = false
  regAuth.setAvatarImage(null)

  const dx = e.clientX - lastPointer.x;
  const dy = e.clientY - lastPointer.y;

  state.x += dx;
  state.y += dy;

  clampPosition();

  lastPointer.x = e.clientX;
  lastPointer.y = e.clientY;
  
  draw();
}

function onPointerUp(e) {
  dragging = false
  try {
    canvasRef.value.releasePointerCapture(e.pointerId)
  } catch (err) { /* ignore */ }
}

function onWheel(e) {
   if (!imgLoaded) return;
   saveState.saved = false
   regAuth.setAvatarImage(null)
  e.preventDefault();

  const canvas = canvasRef.value;
  const rect = canvas.getBoundingClientRect();

  const mouseX = e.clientX - rect.left;
  const mouseY = e.clientY - rect.top;

  const prevScale = state.scale;
  let newScale = prevScale * (e.deltaY < 0 ? 1.1 : 0.9);
  newScale = Math.min(MAX_SCALE, Math.max(MIN_SCALE, newScale));

  state.x = mouseX - ((mouseX - state.x) * newScale) / prevScale;
  state.y = mouseY - ((mouseY - state.y) * newScale) / prevScale;
  state.scale = newScale;

  clampPosition(); 

  draw();
}
function zoomIn() {
  if (!imgLoaded) return;
  saveState.saved =false
  regAuth.setAvatarImage(null)
  
  const canvas = canvasRef.value;
  const rect = canvas.getBoundingClientRect();
  const centerX = rect.width / 2;
  const centerY = rect.height / 2;
  
  zoomAtPosition(centerX, centerY, 1.1);
}

function zoomOut() {
  if (!imgLoaded) return;
  saveState.saved =false
  regAuth.setAvatarImage(null)

  const canvas = canvasRef.value;
  const rect = canvas.getBoundingClientRect();
  const centerX = rect.width / 2;
  const centerY = rect.height / 2;
  
  zoomAtPosition(centerX, centerY, 0.9); //
}
function zoomAtPosition(x, y, factor) {
  const prevScale = state.scale;
  let newScale = prevScale * factor;
  newScale = Math.min(MAX_SCALE, Math.max(MIN_SCALE, newScale));
  
  state.x = x - ((x - state.x) * newScale) / prevScale;
  state.y = y - ((y - state.y) * newScale) / prevScale;
  state.scale = newScale;
  
  clampPosition();
  draw();
}

const handleKeyDown = (e) => {
  if (e.ctrlKey && e.key === '+') {
      zoomIn();
      e.preventDefault();
    } else if (e.ctrlKey && e.key === '-') {
      zoomOut();
      e.preventDefault();
    }
}

window.addEventListener('keydown', handleKeyDown)
onUnmounted(()=> {
  window.removeEventListener('keydown', handleKeyDown)
})

function saveAvatar() {
  if (!imgLoaded) {
    notify.warning("Аватар", "Сначала загрузите изображение!");
    return;
  }

  saveState.proccessing = true;

  const canvas = canvasRef.value;
  const centerX = canvas.width / 2;
  const centerY = canvas.height / 2;
  const radius = canvas.width / 2.5; // Должно совпадать с радиусом из drawImageWithMask

  // Создаем временный canvas
  const tempCanvas = document.createElement('canvas');
  tempCanvas.width = radius * 2;
  tempCanvas.height = radius * 2;
  const tempCtx = tempCanvas.getContext('2d');

  // Вычисляем видимую область в исходных координатах изображения
  const imgCenterX = (centerX - state.x) / state.scale;
  const imgCenterY = (centerY - state.y) / state.scale;
  const imgRadius = radius / state.scale;

  // Вычисляем координаты левого верхнего угла области для вырезки
  const sx = imgCenterX - imgRadius;
  const sy = imgCenterY - imgRadius;
  const sw = imgRadius * 2;
  const sh = imgRadius * 2;

  // Рисуем вырезанную область на временном canvas
  tempCtx.beginPath();
  tempCtx.arc(radius, radius, radius, 0, Math.PI * 2);
  tempCtx.clip();

  tempCtx.drawImage(
    img,
    sx, sy, sw, sh,  // Какая область исходного изображения берется
    0, 0,            // Куда помещаем на tempCanvas
    radius * 2,      // Ширина на tempCanvas
    radius * 2       // Высота на tempCanvas
  );

  // Сохраняем
  tempCanvas.toBlob((blob) => {
    if (!blob) {
      notify.error("Ошибка", "Не удалось создать файл!");
      saveState.proccessing = false;
      return;
    }

    //Сохранем в pinia store вырезанное изображение
    regAuth.setAvatarImage(blob)

    saveState.saved = true;
    saveState.proccessing = false;
  }, 'image/png', 0.9);
}

watch(() => regAuth.currentStep,
(newVal)=> {
  if(newVal === 2) {
    notify.info('Аватар', 'Вы можете перетащить фотографию прямо в область для загрузки фотографии, вы увидите соотвествующие подсказки')
  }
})

</script>

<style scoped>
.canvas-wrap {
  position: relative;
    display: flex;
    align-items: center;
    justify-content: center;
     width: 100%;
    height: 180px;
    margin-top: 75px;
}

.label-photo {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-48%, -50%);
  user-select: none;
  pointer-events: none;
  transition: all .3s ease-in-out;
}
.label-photo-hidden {
  opacity: 0;
}
.label-photo > svg {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-48%, -50%);
  transition: all .15s ease-in-out;
  stroke: rgba(255,255,255, .7);
}
.save-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: #0e1621;
  box-shadow: 3px 3px 6px black;
  padding: 5px;
  border-radius: 8px;
  margin-left: 10px;
  transition: all .3s ease-in-out;
  cursor: not-allowed;
  user-select: none;
  width: 40px;
  height: 40px;
}
.save-btn > svg {
  fill: rgba(255,255,255,.5);
  transition: all .3s ease-in-out;
}

.active-save-btn:hover {
  transform: translateY(-2px);
  box-shadow: 3px 6px 6px black;
}
.active-save-btn:hover > svg {
  fill: rgba(255,255,255,.9);
}
.active-save-btn {
  cursor: pointer;
}

canvas {
  width: 180px;
  height: 180px;
  background-color: #0e1621;
  touch-action: none; 
  cursor: grab;
  border-radius: 50%;
  display: block;
  box-shadow: 2px 2px 8px black;
  transition: all .3s ease-in-out;
  cursor: auto;
}
canvas:active {
  cursor: grabbing;
}

.input-wrapper {
  margin-top: 20px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.input-file {
  display: none;
}

.input-wrapper > label {
  position: relative;
  background-color: #0e1621;
  cursor: pointer;
  border-radius: 8px;
  box-shadow: 3px 3px 6px black;
  padding: 5px;
  width: 40px;
  height: 40px;
  transition: all .3s ease-in-out;
}
.input-wrapper > label > svg {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  opacity: .7;
  transition: all .3s ease-in-out;
}

.input-wrapper > label:hover > svg {
  opacity: 1;
  color: rgba(255,255,255,.9);
}

.input-wrapper > label:hover {
  transform: translateY(-2px);
}

.controls {
  transition: all .3s ease-in-out;
  user-select: none;
  opacity: 0;
}

.controls-visible {
  opacity: 1;
}

.zoom-btn {
  position: absolute;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: bold;
  bottom: 0;
  right: 70px;
  font-size: 24px;
  background-color: #0e1621;
  border-radius: 8px;
  box-shadow: 3px 3px 6px black;
  padding: 5px;
  width: 25px;
  height: 25px;
  cursor: pointer;
  transition: all .2s ease-in-out;
  color: rgba(255,255,255,.5);
}

.zoom-btn:hover {
  transform: translateY(-2px);
  color: rgba(255,255,255,.9);
}

.in {
  margin-bottom: 35px;
}

</style>
