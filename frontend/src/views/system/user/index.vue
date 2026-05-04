<template>
  <div>
    <n-space vertical>
      <n-form inline>
        <n-form-item label="用户名">
          <n-input v-model:value="queryParams.username" placeholder="请输入用户名" clearable />
        </n-form-item>
        <n-form-item label="手机号">
          <n-input v-model:value="queryParams.phone" placeholder="请输入手机号" clearable />
        </n-form-item>
        <n-form-item label="状态">
          <n-select v-model:value="queryParams.status" placeholder="请选择状态" clearable style="width: 120px">
            <n-option :value="1" label="启用" />
            <n-option :value="0" label="禁用" />
          </n-select>
        </n-form-item>
        <n-form-item>
          <n-button type="primary" @click="handleQuery">搜索</n-button>
          <n-button @click="resetQuery">重置</n-button>
        </n-form-item>
      </n-form>

      <n-space>
        <n-button v-permission="['user:add']" type="primary" @click="handleAdd">新增</n-button>
      </n-space>

      <n-data-table
        :columns="columns"
        :data="tableData"
        :loading="loading"
        :pagination="pagination"
        @update:page="handlePageChange"
      />
    </n-space>

    <n-modal v-model:show="showModal" :title="modalTitle" preset="card" style="width: 600px">
      <n-form :model="form" label-placement="left" label-width="auto">
        <n-form-item label="用户名" required>
          <n-input v-model:value="form.username" :disabled="isEdit" />
        </n-form-item>
        <n-form-item v-if="!isEdit" label="密码" required>
          <n-input v-model:value="form.password" type="password" />
        </n-form-item>
        <n-form-item label="昵称">
          <n-input v-model:value="form.nickname" />
        </n-form-item>
        <n-form-item label="手机号">
          <n-input v-model:value="form.phone" />
        </n-form-item>
        <n-form-item label="邮箱">
          <n-input v-model:value="form.email" />
        </n-form-item>
        <n-form-item label="状态">
          <n-radio-group v-model:value="form.status">
            <n-radio :value="1">启用</n-radio>
            <n-radio :value="0">禁用</n-radio>
          </n-radio-group>
        </n-form-item>
      </n-form>
      <template #footer>
        <n-space>
          <n-button @click="showModal = false">取消</n-button>
          <n-button type="primary" :loading="submitLoading" @click="handleSubmit">确定</n-button>
        </n-space>
      </template>
    </n-modal>
  </div>
</template>

<script setup>
import { ref, onMounted, h } from 'vue'
import { NButton, NSpace } from 'naive-ui'
import { getUserList, addUser, updateUser, changeUserStatus, resetUserPwd, deleteUser } from '@/api/user'

const loading = ref(false)
const tableData = ref([])
const queryParams = ref({
  pageNum: 1,
  pageSize: 10,
  username: '',
  phone: '',
  status: null
})
const pagination = ref({
  page: 1,
  pageSize: 10,
  itemCount: 0,
  showSizePicker: true,
  pageSizes: [10, 20, 50]
})

const columns = [
  { title: '用户名', key: 'username' },
  { title: '昵称', key: 'nickname' },
  { title: '手机号', key: 'phone' },
  { title: '邮箱', key: 'email' },
  { title: '状态', key: 'status', render(row) {
    return row.status === 1 ? '启用' : '禁用'
  }},
  { title: '创建时间', key: 'createTime' },
  { title: '操作', key: 'actions', render(row) {
    return h(NSpace, null, {
      default: () => [
        h(NButton, { size: 'small', onClick: () => handleEdit(row) }, { default: () => '编辑' }),
        h(NButton, { size: 'small', type: 'error', onClick: () => handleDelete(row) }, { default: () => '删除' })
      ]
    })
  }}
]

const showModal = ref(false)
const modalTitle = ref('')
const isEdit = ref(false)
const submitLoading = ref(false)
const form = ref({
  id: null,
  username: '',
  password: '',
  nickname: '',
  phone: '',
  email: '',
  status: 1
})

onMounted(() => {
  loadData()
})

async function loadData() {
  loading.value = true
  try {
    const res = await getUserList(queryParams.value)
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
    username: '',
    phone: '',
    status: null
  }
  loadData()
}

function handlePageChange(page) {
  queryParams.value.pageNum = page
  loadData()
}

function handleAdd() {
  isEdit.value = false
  modalTitle.value = '新增用户'
  form.value = {
    id: null,
    username: '',
    password: '',
    nickname: '',
    phone: '',
    email: '',
    status: 1
  }
  showModal.value = true
}

function handleEdit(row) {
  isEdit.value = true
  modalTitle.value = '编辑用户'
  form.value = { ...row }
  showModal.value = true
}

async function handleSubmit() {
  submitLoading.value = true
  try {
    if (isEdit.value) {
      await updateUser(form.value)
    } else {
      await addUser(form.value)
    }
    window.$message.success('操作成功')
    showModal.value = false
    loadData()
  } finally {
    submitLoading.value = false
  }
}

async function handleDelete(row) {
  await deleteUser({ id: row.id })
  window.$message.success('删除成功')
  loadData()
}
</script>
