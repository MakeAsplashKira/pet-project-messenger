<template>
  <div class="notifications-container">
    <div 
      v-for="notification in notifications" 
      :key="notification.id"
      class="notifications"
    >
      <div :class="['notification-type', `notification-${notification.type}`]"></div>
      <div class="notification-icon">
        <svg  v-if="notification.type == 'success'" fill="none" stroke="#4dfa53" stroke-width="2" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7"></path></svg>
        <svg v-if="notification.type == 'error'" fill="none" stroke="#ef4444" stroke-width="2" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12"></path></svg>
        <svg v-if="notification.type == 'warning'" fill="none" stroke="#f59e0b" stroke-width="2" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" d="M12 9v2m0 4h.01M12 5c-3.866 0-7 3.134-7 7s3.134 7 7 7 7-3.134 7-7-3.134-7-7-7z"></path></svg>
        <svg v-if="notification.type == 'info'" fill="none" stroke="#f59e42" stroke-width="2" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" d="M13 16h-1v-4h-1m1-4h.01M12 20c4.418 0 8-3.582 8-8s-3.582-8-8-8-8 3.582-8 8 3.582 8 8 8z"></path></svg>
      </div>
      <div class="notification-content">
        <div class="notification-title">{{notification.title}}</div>
        <div class="notification-text" lang="ru">{{ notification.message }}</div>
      </div>
      <div class="notification-close" @click="remove(notification.id)">
        <svg  aria-hidden="true" display="block" viewBox="0 0 24 24" fill="currentColor"><path d="M7.536 6.264a.9.9 0 0 0-1.272 1.272L10.727 12l-4.463 4.464a.9.9 0 0 0 1.272 1.272L12 13.273l4.464 4.463a.9.9 0 1 0 1.272-1.272L13.273 12l4.463-4.464a.9.9 0 1 0-1.272-1.272L12 10.727z"></path></svg>
      </div>
      <div v-if="notification.message" class="notification-copy" @click="copieToClipboard(notification.message)">
        <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24"><path fill="#ffffff" d="M3.25 9A5.75 5.75 0 0 1 9 3.25h7.013a.75.75 0 0 1 0 1.5H9A4.25 4.25 0 0 0 4.75 9v7.107a.75.75 0 0 1-1.5 0z"/><path fill="#ffffff" d="M18.403 6.793a44.372 44.372 0 0 0-9.806 0a2.011 2.011 0 0 0-1.774 1.76a42.581 42.581 0 0 0 0 9.894a2.01 2.01 0 0 0 1.774 1.76c3.241.362 6.565.362 9.806 0a2.01 2.01 0 0 0 1.774-1.76a42.579 42.579 0 0 0 0-9.894a2.011 2.011 0 0 0-1.774-1.76"/></svg>
      </div>
    </div>
  </div>
</template>

<script>
import { useNotifications } from '@/stores/notifications'
import { mapState, mapActions } from 'pinia'

export default {
  computed: {
    ...mapState(useNotifications, ['notifications'])
  },
  methods: {
    ...mapActions(useNotifications, ['removeNotification']),
    remove(id) {
      this.removeNotification(id)
    },
    copieToClipboard (message) {
      if(!message) return
      navigator.clipboard.writeText(message)
    }
  }
}

</script>

<style scoped>
.notifications-container {
  position: fixed;
  bottom: 20px;
  left: 20px;
  z-index: 1000;
}

.notifications {
  position: relative;
  margin-bottom: 10px;
  color: white;
  max-width: 300px;
  box-shadow: 0 3px 6px black;
  animation: notifictaion-appear 1s ease-in-out;
}

@keyframes notifictaion-appear {
  0% {
    transform: translateX(-150%);
  }
  100% {
    transform: translateX(0%);
  }
}

.notification-info {
  background-color: #f59e42;
}

.notification-success {
  background-color: #4dfa53;
}

.notification-warning {
  background-color: #FF9800;
}

.notification-error {
  background-color: #F44336;
}

.notification-icon {
  position: absolute;
  top: 50%;
  left: 20px;
  transform: translateY(-50%);
}

.notification-icon > svg {
  width: 30px;
  height: 30px;
}

.notification-type {
  position: absolute;
  bottom: 0;
  top: 0;
  left: 0;
  width: 12px;
  border-top-left-radius: 8px;
  border-bottom-left-radius: 8px;
}

.notification-window {
  position: relative;
  width: 100%;
  height: 30px;
  background-color: #0e1621;
  box-shadow: 0px 2px 6px black;
}

.notification-content {
  display: flex;
  flex-direction: column;
  justify-content: center;
  background-color: #1b1f25;
  padding: 10px 10px;
  padding-left: 55px;
  padding-right: 40px;
  font-size: 15px;
  margin-left: 10px;
}

.notification-title {
  font-size: 18px;
  opacity: .9;
  font-weight: bold;
  user-select: none;
}

.notification-text {
  margin-top: 5px;
  font-size: 14px;
  opacity: .8;
  overflow-wrap: break-word;
  word-break: break-all;
  line-height: 18px;
  display: -webkit-box;
  -webkit-line-clamp: 5;       /* Ограничиваем 3 строками */
  -webkit-box-orient: vertical;
  overflow: hidden;
  text-overflow: ellipsis;
  line-height: 1.3;            /* Важно для корректного подсчёта строк */
  max-height: calc(5em * 1.3); /* 2 строки с учётом line-height */
}

.notification-close {
  position: absolute;
  top: 7px;
  right: 7px;
  border-radius: 8px;
  color: rgb(0, 0, 0);
  font-size: 20px;
  cursor: pointer;
  background-color: #ffffff00;
  transition: all .2s ease-in-out;
  padding: 2px;
}

.notification-close > svg {
  width: 20px;
  height: 20px;
  transition: all .2s ease-in-out;
  opacity: .5;
}

.notification-close:hover {
  background-color: #ffffff1c;
  box-shadow: 1px 1px 3px rgba(0, 0, 0, 0.575);
}

.notification-close:hover > svg {
  transform: rotate(-90deg);
  opacity: .9;
}

.notification-copy {
  position: absolute;
  top: 35px;
  right: 7px;
  border-radius: 8px;
  color: rgb(0, 0, 0);
  font-size: 20px;
  cursor: pointer;
  background-color: #ffffff00;
  transition: all .2s ease-in-out;
  padding: 2px;
}
.notification-copy > svg {
  width: 20px;
  height: 20px;
  transition: all .2s ease-in-out;
  opacity: .5;
  transform: translateY(8%);
}

.notification-copy:hover {
   background-color: #ffffff1c;
  box-shadow: 1px 1px 3px rgba(0, 0, 0, 0.575);
}

.notification-copy:hover > svg {
  opacity: .9;
}


</style>