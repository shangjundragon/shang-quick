<template>
  <div>
    <n-space vertical>
      <n-form inline>
        <n-form-item label="文件名">
          <n-input v-model:value="queryParams.originalName" placeholder="请输入文件名" clearable />
        </n-form-item>
        <n-form-item>
          <n-button type="primary" @click="handleQuery">搜索</n-button>
          <n-button @click="resetQuery">重置</n-button>
        </n-form-item>
      </n-form>

      <n-space>
        <n-upload
          v-permission="['file:upload']"
          :custom-request="handleUpload"
          :show-file-list="false"
          accept="*"
        >
          <n-button type="primary">上传文件</n-button>
        </n-upload>
      </n-space>

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
import { ref, onMounted, h } from 'vue'
import { NButton, NSpace, NImage, NPopconfirm, NUpload } from 'naive-ui'
import { getFileList, deleteFile, getFileConfig } from '@/api/file'
import { formatTimestamp } from '@/utils/format'

const loading = ref(false)
const tableData = ref([])
const fileUrlPrefix = ref('')
const queryParams = ref({
  pageNum: 1,
  pageSize: 10,
  originalName: ''
})
const pagination = ref({
  page: 1,
  pageSize: 10,
  itemCount: 0,
  showSizePicker: true,
  pageSizes: [10, 20, 50]
})

const columns = [
  {
    title: '预览',
    key: 'preview',
    width: 100,
    render(row) {
      if (row.isImage === 1) {
        return h(NImage, {
          width: 60,
          height: 60,
          src: fileUrlPrefix.value + row.filePath,
          objectFit: 'cover',
          style: 'border-radius: 4px;'
        })
      }
      return h('div', {
        style: 'width: 60px; height: 60px; display: flex; align-items: center; justify-content: center; background: #f5f5f5; border-radius: 4px; font-size: 12px; color: #999;'
      }, '文件')
    }
  },
  { title: '原始文件名', key: 'originalName' },
  { title: '文件大小', key: 'fileSizeStr', width: 120 },
  { title: '文件类型', key: 'fileType', width: 180 },
  { title: '上传时间', key: 'createTime', width: 180, render(row) { return formatTimestamp(row.createTime) } },
  {
    title: '操作',
    key: 'actions',
    width: 200,
    render(row) {
      return h(NSpace, null, {
        default: () => [
          h(NButton, {
            size: 'small',
            onClick: () => copyUrl(row.filePath)
          }, { default: () => '复制链接' }),
          h(NPopconfirm, {
            onPositiveClick: () => handleDelete(row)
          }, {
            trigger: () => h(NButton, { size: 'small', type: 'error' }, { default: () => '删除' }),
            default: () => '确定删除该文件吗？'
          })
        ]
      })
    }
  }
]

onMounted(() => {
  loadConfig()
  loadData()
})

async function loadConfig() {
  try {
    const res = await getFileConfig()
    fileUrlPrefix.value = res.fileUrlPrefix || ''
  } catch (error) {
    console.error('获取文件配置失败', error)
  }
}

async function loadData() {
  loading.value = true
  try {
    const res = await getFileList(queryParams.value)
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
    originalName: ''
  }
  loadData()
}

function handlePageChange(page) {
  queryParams.value.pageNum = page
  loadData()
}

async function handleUpload({ file, onFinish, onError }) {
  const formData = new FormData()
  formData.append('file', file.file)

  try {
    const { uploadFile } = await import('@/api/file')
    await uploadFile(formData)
    window.$message.success('上传成功')
    loadData()
    onFinish()
  } catch (error) {
    window.$message.error(error?.message || '上传失败')
    onError()
  }
}

async function handleDelete(row) {
  await deleteFile({ id: row.id })
  window.$message.success('删除成功')
  loadData()
}

function copyUrl(filePath) {
  const fullUrl = fileUrlPrefix.value + filePath
  navigator.clipboard.writeText(fullUrl).then(() => {
    window.$message.success('链接已复制')
  }).catch(() => {
    window.$message.error('复制失败')
  })
}
</script>
