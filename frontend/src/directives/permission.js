import { useUserStore } from '@/store/user'

// v-permission 指令：传权限标识数组，用户有任一权限即显示元素
// 用法：<n-button v-permission="['user:add']">新增</n-button>
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
    // 权限匹配规则：用户权限列表中包含任一所列标识即可
    const hasPermission = userStore.permissions.some(p => value.includes(p))
    if (!hasPermission) {
      el.parentNode && el.parentNode.removeChild(el)
    }
  }
}
