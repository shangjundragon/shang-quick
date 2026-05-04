<template>
  <div>
    <n-space vertical>
      <n-form inline>
        <n-form-item label="角色名称">
          <n-input v-model:value="queryParams.roleName" placeholder="请输入角色名称" clearable />
        </n-form-item>
        <n-form-item>
          <n-button type="primary" @click="handleQuery">搜索</n-button>
          <n-button @click="resetQuery">重置</n-button>
        </n-form-item>
      </n-form>

      <n-space>
        <n-button v-permission="['role:add']" type="primary" @click="handleAdd">新增</n-button>
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
import { NButton, NSpace } from 'naive-ui'
import { getRoleList, addRole, updateRole, deleteRole } from '@/api/role'

const loading = ref(false)
const tableData = ref([])
const queryParams = ref({
  pageNum: 1,
  pageSize: 10,
  roleName: ''
})
const pagination = ref({
  page: 1,
  pageSize: 10,
  itemCount: 0,
  showSizePicker: true,
  pageSizes: [10, 20, 50]
})

const columns = [
  { title: '角色名称', key: 'roleName' },
  { title: '角色编码', key: 'roleCode' },
  { title: '备注', key: 'remark' },
  { title: '状态', key: 'status', render(row) {
    return row.status === 1 ? '启用' : '禁用'
  }},
  { title: '操作', key: 'actions', render(row) {
    return h(NSpace, null, {
      default: () => [
        h(NButton, { size: 'small', onClick: () => handleEdit(row) }, { default: () => '编辑' }),
        h(NButton, { size: 'small', type: 'error', onClick: () => handleDelete(row) }, { default: () => '删除' })
      ]
    })
  }}
]

onMounted(() => {
  loadData()
})

async function loadData() {
  loading.value = true
  try {
    const res = await getRoleList(queryParams.value)
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
    roleName: ''
  }
  loadData()
}

function handlePageChange(page) {
  queryParams.value.pageNum = page
  loadData()
}

function handleAdd() {}
function handleEdit(row) {}
async function handleDelete(row) {
  await deleteRole({ id: row.id })
  window.$message.success('删除成功')
  loadData()
}
</script>
