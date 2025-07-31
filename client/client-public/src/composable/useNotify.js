import { useNotifications } from '@/stores/notifications'

export default function() {
  const notifyStore = useNotifications()
  
  const notify = {
    success: (msg) => notifyStore.addNotification({ message: msg, type: 'success' }),
    error: (msg) => notifyStore.addNotification({ message: msg, type: 'error' }),
    warning: (msg) => notifyStore.addNotification({ message: msg, type: 'warning' })
  }
  
  return { notify }
}