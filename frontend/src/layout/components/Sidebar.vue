<template>
  <n-layout-sider
    collapsible
    :collapsed-width="64"
    :width="200"
    :collapsed="collapsed"
    @update:collapsed="$emit('update:collapsed', $event)"
    bordered
  >
    <div class="logo">
      <span v-if="!collapsed">Admin</span>
    </div>
    <n-menu
      :value="activeKey"
      :collapsed="collapsed"
      :collapsed-width="64"
      :collapsed-icon-size="22"
      :options="menuOptions"
      @update:value="handleMenuSelect"
    />
  </n-layout-sider>
</template>

<script setup>
import { computed, h } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useUserStore } from '@/store/user'
import * as Icons from '@vicons/ionicons5'

const props = defineProps({
  collapsed: Boolean
})

defineEmits(['update:collapsed'])

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

const activeKey = computed(() => route.path)

const menuOptions = computed(() => {
  return generateMenu(userStore.menus)
})

function generateMenu(menus) {
  const menuMap = {}
  const result = []

  // 只保留目录(0)和菜单(1)，过滤按钮(2)
  const displayMenus = menus.filter(menu => menu.menuType !== 2)

  // 计算完整路径
  function getFullPath(menu) {
    if (!menu.path) return ''
    if (menu.path.startsWith('/')) return menu.path
    
    // 查找父菜单
    const parent = displayMenus.find(m => m.id === menu.parentId)
    if (parent) {
      const parentPath = getFullPath(parent)
      if (parentPath) {
        return parentPath.replace(/\/$/, '') + '/' + menu.path
      }
    }
    return '/' + menu.path
  }

  displayMenus.forEach(menu => {
    const iconName = menu.icon || ''
    const iconComponent = Icons[iconName]
    const fullPath = getFullPath(menu)
    
    menuMap[menu.id] = {
      label: menu.menuName,
      key: fullPath,
      icon: iconComponent ? () => h(iconComponent) : undefined
    }
  })

  displayMenus.forEach(menu => {
    if (menu.parentId === 0) {
      result.push(menuMap[menu.id])
    } else if (menuMap[menu.parentId]) {
      if (!menuMap[menu.parentId].children) {
        menuMap[menu.parentId].children = []
      }
      menuMap[menu.parentId].children.push(menuMap[menu.id])
    }
  })

  return result
}

function handleMenuSelect(key) {
  if (key) {
    router.push(key)
  }
}
</script>

<style scoped>
.logo {
  height: 64px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 20px;
  font-weight: bold;
  color: #18a058;
  border-bottom: 1px solid #eee;
}
</style>
