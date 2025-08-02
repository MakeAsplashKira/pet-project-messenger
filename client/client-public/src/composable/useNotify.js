import { useNotifications } from '@/stores/notifications'

export default function() {
  const notifyStore = useNotifications()
  
  const notify = {
    success: (title, msg) => notifyStore.addNotification({ message: msg, title: title, type: 'success'}),
    error: (title, msg) => notifyStore.addNotification({ message: msg, title: title, type: 'error'}),
    warning: (title, msg) => notifyStore.addNotification({ message: msg, title: title, type: 'warning'}),
    info: (title, msg) => notifyStore.addNotification({ message: msg, title: title, type: 'info'})
  }
  
  return { notify }
}