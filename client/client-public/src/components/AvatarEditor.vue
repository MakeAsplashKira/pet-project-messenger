<template>
  <div class="uploader-step">
    <input type="file" accept="image/*" @change="onFileChange" />

    <div class="canvas-wrap">
      <canvas
        ref="canvasRef"
        width="400"
        height="400"
        @pointerdown="onPointerDown"
        @pointermove="onPointerMove"
        @pointerup="onPointerUp"
        @pointercancel="onPointerUp"
        @wheel.prevent="onWheel"
      ></canvas>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'

const MIN_SCALE = 0.5
const MAX_SCALE = 2
const canvasRef = ref(null)
const img = new Image()
let imgLoaded = false

// состояние изображения (позиция и масштаб)
const state = {
  x: 0,
  y: 0,
  scale: 1,
  imgWidth: 0,
  imgHeight: 0
}

let dragging = false
const lastPointer = { x: 0, y: 0 }

function onFileChange(e) {
  const file = e.target.files?.[0]
  if (!file) return

  imgLoaded = false
  const url = URL.createObjectURL(file)
  img.onload = () => {
    URL.revokeObjectURL(url)
    state.imgWidth = img.width
    state.imgHeight = img.height

    const canvas = canvasRef.value
    const cw = canvas.width
    const ch = canvas.height

    // scale так, чтобы изображение покрывало canvas (cover)
    const scale = Math.max(cw / img.width, ch / img.height)
    state.scale = scale
    state.x = (cw - img.width * scale) / 2
    state.y = (ch - img.height * scale) / 2

    imgLoaded = true
    draw()
  }
  img.src = url
}

function drawImageWithMask(ctx, canvas) {
  ctx.clearRect(0, 0, canvas.width, canvas.height);
  
  // 1. Сначала рисуем изображение
  ctx.drawImage(
    img,
    state.x, state.y,
    img.width * state.scale,
    img.height * state.scale
  );
  
  // 2. Сохраняем состояние контекста
  ctx.save();
  
  // 3. Создаем круглую маску
  ctx.beginPath();
  ctx.arc(
    canvas.width / 2,
    canvas.height / 2,
    canvas.width / 2.5,
    0,
    Math.PI * 2
  );
  ctx.clip(); // Все последующие рисунки будут обрезаны по кругу
  
  // 4. Рисуем изображение еще раз (только внутри круга)
  ctx.drawImage(
    img,
    state.x, state.y,
    img.width * state.scale,
    img.height * state.scale
  );
  
  // 5. Восстанавливаем состояние
  ctx.restore();
  
  // 6. Рисуем затемнение вокруг круга
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
  
  // 7. Рисуем обводку круга
  ctx.beginPath();
  ctx.arc(
    canvas.width / 2,
    canvas.height / 2,
    canvas.width / 2.5,
    0,
    Math.PI * 2
  );
  ctx.strokeStyle = 'white';
  ctx.lineWidth = 2;
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
  
  const canvas = canvasRef.value;
  const dx = e.clientX - lastPointer.x;
  const dy = e.clientY - lastPointer.y;
  
  // Рассчитываем новые координаты
  let newX = state.x + dx;
  let newY = state.y + dy;
  
  // Рассчитываем границы перемещения
  const imgWidth = img.width * state.scale;
  const imgHeight = img.height * state.scale;
  const centerX = canvas.width / 2;
  const centerY = canvas.height / 2;
  const radius = canvas.width / 2.5; // Радиус видимой области
  
  // Максимальные смещения (чтобы центр изображения не выходил за круг)
  const maxX = centerX + radius - imgWidth/2;
  const minX = centerX - radius - imgWidth/2;
  const maxY = centerY + radius - imgHeight/2;
  const minY = centerY - radius - imgHeight/2;
  
  // Ограничиваем координаты
  newX = Math.max(minX, Math.min(maxX, newX));
  newY = Math.max(minY, Math.min(maxY, newY));
  
  // Обновляем состояние
  state.x = newX;
  state.y = newY;
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
  if (!imgLoaded) return
  e.preventDefault()

  const canvas = canvasRef.value
  const rect = canvas.getBoundingClientRect()

  // координаты курсора относительно canvas
  const mouseX = e.clientX - rect.left
  const mouseY = e.clientY - rect.top

  // текущий масштаб и позиция
  const prevScale = state.scale
  let newScale = prevScale * (e.deltaY < 0 ? 1.1 : 0.9)

  newScale = Math.min(MAX_SCALE, Math.max(MIN_SCALE, newScale))

  // чтобы зум был относительно курсора, подгоняем x и y
  state.x = mouseX - ((mouseX - state.x) * newScale) / prevScale
  state.y = mouseY - ((mouseY - state.y) * newScale) / prevScale
  state.scale = newScale

  draw()
}
</script>

<style scoped>
.canvas-wrap {
    display: flex;
    align-items: center;
    justify-content: center;
     width: 100%;
    height: 200px;
    margin-top: 80px;
}
canvas {
  width: 200px;
  height: 200px;
  background-color: #0e1621;
  touch-action: none; 
  cursor: grab;
  border-radius: 8px;
  display: block;
  transition: all .3s ease-in-out;
  box-shadow: 2px 2px 8px black;
}
canvas:active {
  cursor: grabbing;
}
</style>
