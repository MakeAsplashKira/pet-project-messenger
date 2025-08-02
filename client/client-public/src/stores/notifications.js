import { defineStore } from "pinia";

export const useNotifications = defineStore('notifications', {
    state: () => ({
        notifications: []
    }),
    actions: {
        addNotification(notification) {
            const id = Date.now()
            const newNotification = {
                id,
                message: notification.message,
                title: notification.title || notification.type,
                type: notification.type || 'info',
                timeout: notification.timeout || 10000
            }

            this.notifications.push(newNotification)

            if (newNotification.timeout > 0) {
                setTimeout(()=> {
                    this.removeNotification(id)
                }, newNotification.timeout)
            }
        },
        removeNotification(id) {
            this.notifications = this.notifications.filter(n => n.id !== id)
        },
        clearAll() {
            this.notifications = []
        }
    }
})