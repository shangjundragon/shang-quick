import { useUserStore } from '@/store/user'

export const permission = {
  mounted(el, binding) {
    checkPermission(el, binding)
  },
  updated(el, binding) {
    checkPermission(el, binding)
  }
}

function checkPermission(el, binding) {
  const userStore = useUserStore()
  const { value } = binding

  if (value && value instanceof Array && value.length > 0) {
    const hasPermission = userStore.permissions.some(p => value.includes(p))
    if (!hasPermission) {
      el.parentNode && el.parentNode.removeChild(el)
    }
  }
}
