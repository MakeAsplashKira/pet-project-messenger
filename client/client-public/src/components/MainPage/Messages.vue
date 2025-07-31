<template>
    <div class="main-content-messages w600">
        <div class="main-content-messages-header">
            <div class="message-user-profile">
                <div class="message-user-profile-photo">
                   
                </div>
                <div class="message-user-profile-data">
                    <div class="user-name">Артемов Александр</div>
                    <div class="user-last-online">Заходил <span class="hightlight">34</span> минуты назад</div>
                </div>
            </div>
            <div class="message-menu"><svg aria-hidden="true" display="block" viewBox="0 0 12 8" fill="none"><path fill="currentColor" fill-rule="evenodd" d="M2.156 2.295a.75.75 0 0 1 1.051-.137L6 4.306l2.793-2.148a.75.75 0 1 1 .914 1.189l-3.25 2.5a.75.75 0 0 1-.914 0l-3.25-2.5a.75.75 0 0 1-.137-1.052" clip-rule="evenodd"></path></svg></div>
            <div class="svg-close" @click="ChatClose"><svg aria-hidden="true" display="block" class="vkuiIcon vkuiIcon--24 vkuiIcon--w-24 vkuiIcon--h-24 vkuiIcon--cancel_24" width="24" height="24" viewBox="0 0 24 24" fill="currentColor" style="width: 24px; height: 24px;"><path d="M7.536 6.264a.9.9 0 0 0-1.272 1.272L10.727 12l-4.463 4.464a.9.9 0 0 0 1.272 1.272L12 13.273l4.464 4.463a.9.9 0 1 0 1.272-1.272L13.273 12l4.463-4.464a.9.9 0 1 0-1.272-1.272L12 10.727z"></path></svg></div>
        </div>
        <div class="message-content-wrapper">
            <div class="message-content">
                    <div class="message" v-for="(msg, index) in messages" :key="index">
                        {{ msg }}
                    </div>
            </div>
            <div class="message-input">
                <div class="add-file"><svg xmlns="http://www.w3.org/2000/svg" width="1024" height="1024" viewBox="0 0 1024 1024"><path fill="#d4d4d4" d="M779.3 196.6c-94.2-94.2-247.6-94.2-341.7 0l-261 260.8c-1.7 1.7-2.6 4-2.6 6.4s.9 4.7 2.6 6.4l36.9 36.9a9 9 0 0 0 12.7 0l261-260.8c32.4-32.4 75.5-50.2 121.3-50.2s88.9 17.8 121.2 50.2c32.4 32.4 50.2 75.5 50.2 121.2c0 45.8-17.8 88.8-50.2 121.2l-266 265.9l-43.1 43.1c-40.3 40.3-105.8 40.3-146.1 0c-19.5-19.5-30.2-45.4-30.2-73s10.7-53.5 30.2-73l263.9-263.8c6.7-6.6 15.5-10.3 24.9-10.3h.1c9.4 0 18.1 3.7 24.7 10.3c6.7 6.7 10.3 15.5 10.3 24.9c0 9.3-3.7 18.1-10.3 24.7L372.4 653c-1.7 1.7-2.6 4-2.6 6.4s.9 4.7 2.6 6.4l36.9 36.9a9 9 0 0 0 12.7 0l215.6-215.6c19.9-19.9 30.8-46.3 30.8-74.4s-11-54.6-30.8-74.4c-41.1-41.1-107.9-41-149 0L463 364L224.8 602.1A172.22 172.22 0 0 0 174 724.8c0 46.3 18.1 89.8 50.8 122.5c33.9 33.8 78.3 50.7 122.7 50.7c44.4 0 88.8-16.9 122.6-50.7l309.2-309C824.8 492.7 850 432 850 367.5c.1-64.6-25.1-125.3-70.7-170.9"/></svg></div>
                <div class="input"
                ref="inputRef"
                contenteditable="true"
                @input="updateHasText"
                @keydown.enter.prevent="sendMessage"
                @paste="handlePaste">
                </div>
                <div class="placeholder" :style="{opacity: hasText? 0 : 0.33}">Сообщение...</div>
                <div class="message-send-menu">
                    <div class="smile-btn" @click="toggleEmojiPicker">
                        <div class="emoji-panel" v-if="showEmojiPicker">
                            <span v-for="emoji in emojis" :key="emoji" @click="insertEmoji(emoji)">{{ emoji }}</span>
                        </div>
                        <svg xmlns="http://www.w3.org/2000/svg" width="32" height="32" viewBox="0 0 1024 1024"><path d="M288 421a48 48 0 1 0 96 0a48 48 0 1 0-96 0m352 0a48 48 0 1 0 96 0a48 48 0 1 0-96 0M512 64C264.6 64 64 264.6 64 512s200.6 448 448 448s448-200.6 448-448S759.4 64 512 64m263 711c-34.2 34.2-74 61-118.3 79.8C611 874.2 562.3 884 512 884c-50.3 0-99-9.8-144.8-29.2A370.4 370.4 0 0 1 248.9 775c-34.2-34.2-61-74-79.8-118.3C149.8 611 140 562.3 140 512s9.8-99 29.2-144.8A370.4 370.4 0 0 1 249 248.9c34.2-34.2 74-61 118.3-79.8C413 149.8 461.7 140 512 140c50.3 0 99 9.8 144.8 29.2A370.4 370.4 0 0 1 775.1 249c34.2 34.2 61 74 79.8 118.3C874.2 413 884 461.7 884 512s-9.8 99-29.2 144.8A368.89 368.89 0 0 1 775 775M664 533h-48.1c-4.2 0-7.8 3.2-8.1 7.4C604 589.9 562.5 629 512 629s-92.1-39.1-95.8-88.6c-.3-4.2-3.9-7.4-8.1-7.4H360a8 8 0 0 0-8 8.4c4.4 84.3 74.5 151.6 160 151.6s155.6-67.3 160-151.6a8 8 0 0 0-8-8.4"/></svg>
                    </div>
                    <div class="send-btn" @click="sendMessage"><svg xmlns="http://www.w3.org/2000/svg" width="32" height="32" viewBox="0 0 24 24"><path d="M21.243 12.437a.5.5 0 0 0 0-.874l-2.282-1.268A75.497 75.497 0 0 0 4.813 4.231l-.665-.208A.5.5 0 0 0 3.5 4.5v5.75a.5.5 0 0 0 .474.5l1.01.053a44.41 44.41 0 0 1 7.314.998l.238.053c.053.011.076.033.089.05a.163.163 0 0 1 .029.096c0 .04-.013.074-.029.096c-.013.017-.036.039-.089.05l-.238.053a44.509 44.509 0 0 1-7.315.999l-1.01.053a.5.5 0 0 0-.473.499v5.75a.5.5 0 0 0 .65.477l.664-.208a75.499 75.499 0 0 0 14.147-6.064z"/></svg></div>
                    <div class="voice-btn"></div>
                </div>
            </div>
        </div>
        
    </div>
</template>

<script setup>
import { ref, nextTick, onMounted } from 'vue';

const messages = ref(['Привет!', 'Как дела?']);
const inputRef = ref(null);
const messagesRef = ref(null);
const showEmojiPicker = ref(false);
const emojis = ['😊', '❤️', '🔥', '👍', '😂', '😍'];
const hasText = ref(false)

const updateHasText = () => {
  try {
    // Защищённая проверка
    hasText.value = inputRef.value?.textContent?.trim().length > 0;
  } catch (e) {
    hasText.value = false;
  }
};

// Инициализация contenteditable
onMounted(() => {
  if (inputRef.value) {
    inputRef.value.focus();
  }
});

// Автоскролл
const scrollToBottom = () => {
  nextTick(() => {
    if (messagesRef.value) {
      messagesRef.value.scrollTop = messagesRef.value.scrollHeight;
    }
  });
};

// Вставка эмодзи
const insertEmoji = (emoji) => {
  if (!inputRef.value) return;

  const selection = window.getSelection();
  if (!selection) return;

  // Для пустого инпута
  if (inputRef.value.childNodes.length === 0) {
    inputRef.value.textContent = emoji;
  } else {
    try {
      const range = selection.getRangeAt(0);
      range.deleteContents();
      range.insertNode(document.createTextNode(emoji));
      
      // Обновляем выделение
      const newRange = document.createRange();
      newRange.selectNodeContents(inputRef.value);
      newRange.collapse(false);
      selection.removeAllRanges();
      selection.addRange(newRange);
    } catch (e) {
      // Фолбэк: добавляем в конец
      inputRef.value.textContent += emoji;
    }
  }

  inputRef.value.focus();
  showEmojiPicker.value = false;
};

// Отправка сообщения
const sendMessage = () => {
  if (!inputRef.value) return;

  const text = inputRef.value.textContent.trim();
  if (!text) return;

  messages.value.push(text);
  inputRef.value.textContent = '';
  scrollToBottom();
};

const toggleEmojiPicker = () => {
  showEmojiPicker.value = !showEmojiPicker.value;
  inputRef.value?.focus();
};

const handleInput = () => {
  // Можно добавить логику автовысоты
};

const handlePaste = (e) => {
  e.preventDefault();
  const text = e.clipboardData.getData('text/plain');
  document.execCommand('insertText', false, text);
};
</script>

<style scoped>
.main-content-messages {
  position: relative;
  margin-left: 12px;
  min-height: calc(100vh - 120px);
  border-radius: 8px;
  background-color: #17212b;
  box-shadow: 0 4px 8px black;
  overflow: visible;
}

.main-content-messages-header {
    position: absolute;
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 0 15px;
    padding-top: 15px;
    padding-bottom: 15px;
    top: 0;
    left: 0;
    right: 0;
    height: 60px;
    border-top-left-radius: 8px;
    border-top-right-radius: 8px;
}

.message-menu {
    cursor: pointer;
    padding: 5px;
    transition: all .2s ease-in-out;
}

.message-menu:hover {
    border-radius: 4px;
    background-color: #31475e7a;
}

.message-menu > svg {
    width: 20px;
    height: 20px;
    margin-top: 1px;
    opacity: .5;
    transition: all .2s ease-in-out;
}

.message-menu:hover > svg {
    opacity: .8;
}



.main-content-messages-header::after {
    content: "";
    position: absolute;
    left: 0;
    right: 0;
    top: 60px;
    height: 2px;
    box-shadow: 0 6px 10px black;
    background-color: rgba(163, 163, 163, 0.308);
}

.message-user-profile {
    position: relative;
    display: flex;
    align-items: center;
    cursor: pointer;
    transition: all .3s ease-in-out;
    background: transparent;
}



.message-user-profile-photo {
    width: 35px;
    height: 35px;
    border-radius: 50%;
    background-color: #22303f;
    box-shadow: 2px 2px 4px rgba(0, 0, 0, 0.686);
}

.message-user-profile-data {
    margin-left: 15px;
    font-size: 2px;
}

.user-name {
    font-size: 16px;
}

.user-last-online {
    font-size: 14px;
    color: #7c7c7c;
}


.user-last-online > span {
    font-weight: bold;
    font-size: 14px;
    opacity: .7;    
    font-family: 'Segoe UI', -apple-system, BlinkMacSystemFont, sans-serif;
}



.svg-close {
    padding: 4px;
    transition: all .3s ease-in-out;
    border-radius: 4px;
    cursor: pointer;
}

.svg-close svg {
    fill: #ffffff7e;
    transition: all .3s ease-in-out;
}

.svg-close:hover {
    background-color: #31475e7a;
}

.svg-close:hover > svg {
    rotate: 90deg;
}

.svg-close:hover > svg {
    fill: whitesmoke;
}


.message-content-wrapper {
    display: flex;
    flex-direction: column;
    align-items: center;
    margin-top: 60px;
    width: 100%;
    height: calc(100% - 60px);
}

.message-content {
    flex: 1;
    padding: 0px 15px;
    width: 100%;
    overflow-y: auto;
    margin-bottom: 50px;
}


.message-input {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-bottom: 15px;
    border-radius: 20px;
    width: calc(100% - 28px);
    min-height: 42px;
    max-height: 168px;
    background-color: #0e1621;
    box-shadow: 3px 3px 8px black;
}

.input {
  flex: 1;
  margin: 0px 5px;
  max-height: 120px;
  min-height: 20px;
  overflow-y: auto;
  font-size: 16px;
  outline: none;
  padding: 5px 0;
  color: rgba(255, 255, 255, 0.866);
}

.placeholder {
    transition: all .1s ease-in-out;
    position: absolute;
    left: 65px;
}

.add-file {
    margin-left: 10px;
    cursor: pointer;
}

.add-file > svg {
    width: 32px;
    height: 32px;
    opacity: .5;
    transition: all .2s ease-in-out;
}

.add-file > svg:hover {
    opacity: 1;
}




.message-send-menu {
    display: flex;
}

.send-btn {
 margin-right: 5px;
}

.send-btn > svg {
    cursor: pointer;
    fill: whitesmoke;
    opacity: .5;
    transition: all .2s ease-in-out;
}

.send-btn:hover > svg {
    opacity: .9;
    fill: #5288c1;
}

.smile-btn {
    position: relative;
    cursor: pointer;
    margin-right: 5px;
}

.smile-btn > svg {
    width: 32px;
    height: 32px;
    fill: whitesmoke;
    opacity: .5;
    transition: all .2s ease-in-out;

}

.smile-btn:hover > svg {
    opacity: 1;
}


.emoji-panel {
  position: absolute;  
  display: flex;
  flex-wrap: wrap;
  gap: 5px;
  padding: 10px;
  background-color: #0e1621;
  box-shadow: 2px 2px 4px black;
  border-radius: 8px;
  left: -270px;
  top: -60px;
  width: 300px;
}

.emoji-panel span {
  cursor: pointer;
  font-size: 20px;
}

.w600 {
    width: 600px;
}


</style>