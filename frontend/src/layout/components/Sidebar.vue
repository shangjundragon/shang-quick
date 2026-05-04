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

  menus.forEach(menu => {
    menuMap[menu.id] = {
      label: menu.menuName,
      key: menu.path,
      icon: () => h('i', { class: menu.icon }),
      children: []
    }
  })

  menus.forEach(menu => {
    if (menu.parentId === 0) {
      result.push(menuMap[menu.id])
    } else if (menuMap[menu.parentId]) {
      menuMap[menu.parentId].children.push(menuMap[menu.id])
    }
  })

  return result
}

function handleMenuSelect(key) {
  router.push(key)
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
