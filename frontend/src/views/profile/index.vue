<template>
  <div>
    <n-card title="个人中心">
      <n-form ref="formRef" :model="userInfo" :rules="rules" label-placement="left" label-width="auto">
        <n-form-item label="用户名" path="username">
          <n-input v-model:value="userInfo.username" disabled />
        </n-form-item>
        <n-form-item label="昵称" path="nickname">
          <n-input v-model:value="userInfo.nickname" />
        </n-form-item>
        <n-form-item label="手机号" path="phone">
          <n-input v-model:value="userInfo.phone" />
        </n-form-item>
        <n-form-item label="邮箱" path="email">
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
import { ref, onMounted, useTemplateRef } from 'vue'
import { useUserStore } from '@/store/user'
import { updateProfile } from '@/api/profile'

const userStore = useUserStore()
const formRef = useTemplateRef('formRef')
const userInfo = ref({})

const rules = {
  nickname: { required: true, message: '请输入昵称', trigger: 'blur' }
}

onMounted(() => {
  userInfo.value = { ...userStore.userInfo }
})

async function handleUpdate() {
  try {
    await formRef.value?.validate()
    const res = await updateProfile(userInfo.value)
    userStore.setUserInfo(res)
    window.$message.success('保存成功')
  } catch (error) {
    window.$message.error('保存失败')
  }
}
</script>
