<template>
  <div class="notifications-container">
    <div 
      v-for="notification in notifications" 
      :key="notification.id"
      :class="['notification', `notification-${notification.type}`]"
    >
      <div class="notification-content">
        {{ notification.message }}
      </div>
      <button 
        class="notification-close"
        @click="remove(notification.id)"
      >
        &times;
      </button>
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

.notification {
  position: relative;
  padding: 15px;
  margin-bottom: 10px;
  border-radius: 4px;
  color: white;
  min-width: 250px;
  box-shadow: 0 2px 10px rgba(0,0,0,0.1);
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.notification-info {
  background-color: #2196F3;
}

.notification-success {
  background-color: #4CAF50;
}

.notification-warning {
  background-color: #FF9800;
}

.notification-error {
  background-color: #F44336;
}

.notification-close {
  background: transparent;
  border: none;
  color: white;
  font-size: 20px;
  cursor: pointer;
  margin-left: 10px;
}
</style>