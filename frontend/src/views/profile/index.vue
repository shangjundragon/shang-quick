<template>
  <div>
    <n-card title="个人中心">
      <n-form label-placement="left" label-width="auto">
        <n-form-item label="用户名">
          <n-input v-model:value="userInfo.username" disabled />
        </n-form-item>
        <n-form-item label="昵称">
          <n-input v-model:value="userInfo.nickname" />
        </n-form-item>
        <n-form-item label="手机号">
          <n-input v-model:value="userInfo.phone" />
        </n-form-item>
        <n-form-item label="邮箱">
          <n-input v-model:value="userInfo.email" />
        </n-form-item>
        <n-form-item>
          <n-button type="primary" @click="handleUpdate">保存</n-button>
        </n-form-item>
      </n-form>
    </n-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useUserStore } from '@/store/user'

const userStore = useUserStore()
const userInfo = ref({})

onMounted(() => {
  userInfo.value = { ...userStore.userInfo }
})

function handleUpdate() {
  window.$message.success('保存成功')
}
</script>
