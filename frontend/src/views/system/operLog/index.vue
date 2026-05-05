<template>
  <div>
    <n-space vertical>
      <n-form inline>
        <n-form-item label="操作模块">
          <n-input v-model:value="queryParams.title" placeholder="请输入操作模块" clearable />
        </n-form-item>
        <n-form-item label="操作人员">
          <n-input v-model:value="queryParams.operName" placeholder="请输入操作人员" clearable />
        </n-form-item>
        <n-form-item>
          <n-button type="primary" @click="handleQuery">搜索</n-button>
          <n-button @click="resetQuery">重置</n-button>
        </n-form-item>
      </n-form>

      <n-data-table
        :columns="columns"
        :data="tableData"
        :loading="loading"
        :pagination="pagination"
        @update:page="handlePageChange"
      />
    </n-space>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { getOperLogList } from '@/api/operLog'
import { formatTimestamp } from '@/utils/format'

const loading = ref(false)
const tableData = ref([])
const queryParams = ref({
  pageNum: 1,
  pageSize: 10,
  title: '',
  operName: ''
})
const pagination = ref({
  page: 1,
  pageSize: 10,
  itemCount: 0,
  showSizePicker: true,
  pageSizes: [10, 20, 50]
})

const columns = [
  { title: '操作模块', key: 'title' },
  { title: '请求方式', key: 'method' },
  { title: '请求URL', key: 'requestUrl' },
  { title: '操作人员', key: 'operName' },
  { title: '操作IP', key: 'operIp' },
  { title: '操作时间', key: 'operTime', render(row) { return formatTimestamp(row.operTime) } },
  { title: '状态', key: 'status', render(row) {
    return row.status === 1 ? '成功' : '失败'
  }}
]

onMounted(() => {
  loadData()
})

async function loadData() {
  loading.value = true
  try {
    const res = await getOperLogList(queryParams.value)
    tableData.value = res.list
    pagination.value.itemCount = res.total
  } finally {
    loading.value = false
  }
}

function handleQuery() {
  queryParams.value.pageNum = 1
  loadData()
}

function resetQuery() {
  queryParams.value = {
    pageNum: 1,
    pageSize: 10,
    title: '',
    operName: ''
  }
  loadData()
}

function handlePageChange(page) {
  queryParams.value.pageNum = page
  loadData()
}
</script>
