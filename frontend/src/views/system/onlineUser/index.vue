<template>
  <div>
    <n-space vertical>
      <n-data-table
        :columns="columns"
        :data="tableData"
        :loading="loading"
      />
    </n-space>
  </div>
</template>

<script setup>
import { ref, onMounted, h } from 'vue'
import { NButton, NSpace, NPopconfirm } from 'naive-ui'
import { getOnlineUserList, kickOnlineUser } from '@/api/onlineUser'
import { formatTimestamp } from '@/utils/format'

const loading = ref(false)
const tableData = ref([])

const columns = [
  { title: '用户名', key: 'username' },
  { title: '昵称', key: 'nickname' },
  { title: 'IP地址', key: 'ipAddr' },
  { title: '浏览器', key: 'browser' },
  { title: '操作系统', key: 'os' },
  { title: '登录时间', key: 'loginTime', render(row) { return formatTimestamp(row.loginTime) } },
  { title: '最后活跃', key: 'lastActiveTime', render(row) { return formatTimestamp(row.lastActiveTime) } },
  { title: '操作', key: 'actions', render(row) {
    return h(NPopconfirm, { onPositiveClick: () => handleKick(row) }, {
      trigger: () => h(NButton, { size: 'small', type: 'error', vPermission: ['onlineUser:kick'] }, { default: () => '踢出' }),
      default: () => '确定踢出该用户吗？'
    })
  }}
]

onMounted(() => {
  loadData()
})

async function loadData() {
  loading.value = true
  try {
    const res = await getOnlineUserList()
    tableData.value = res
  } finally {
    loading.value = false
  }
}

async function handleKick(row) {
  try {
    await kickOnlineUser({ tokenId: row.tokenId })
    window.$message.success('踢出成功')
    loadData()
  } catch (e) {
    window.$message.error(e.message || '踢出失败')
  }
}
</script>
