<template>
  <div class="navbar">
    <n-space align="center">
      <n-button text @click="$emit('toggle-collapse')">
        <n-icon size="20">
          <menu-outline />
        </n-icon>
      </n-button>
      <span class="breadcrumb">{{ route.meta.title || '仪表盘' }}</span>
    </n-space>
    <n-space align="center">
      <n-dropdown :options="userOptions" @select="handleSelect">
        <n-space align="center" style="cursor: pointer">
          <n-avatar round size="small" :src="userStore.userInfo?.avatar" />
          <span>{{ userStore.userInfo?.nickname || userStore.userInfo?.username }}</span>
        </n-space>
      </n-dropdown>
    </n-space>
  </div>
</template>

<script setup>
import { useRoute, useRouter } from 'vue-router'
import { useUserStore } from '@/store/user'
import { MenuOutline } from '@vicons/ionicons5'
import { h } from 'vue'
import { NIcon } from 'naive-ui'
import { LogOutOutline, PersonOutline } from '@vicons/ionicons5'

defineEmits(['toggle-collapse'])

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

const userOptions = [
  {
    label: '个人中心',
    key: 'profile',
    icon: () => h(NIcon, null, { default: () => h(PersonOutline) })
  },
  {
    label: '退出登录',
    key: 'logout',
    icon: () => h(NIcon, null, { default: () => h(LogOutOutline) })
  }
]

function handleSelect(key) {
  if (key === 'profile') {
    router.push('/profile')
  } else if (key === 'logout') {
    userStore.logout()
    router.push('/login')
  }
}
</script>

<style scoped>
.navbar {
  height: 64px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 20px;
  border-bottom: 1px solid #eee;
  background: #fff;
}

.breadcrumb {
  font-size: 16px;
  font-weight: 500;
}
</style>
