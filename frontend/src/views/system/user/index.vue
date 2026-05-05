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
          <n-select v-model:value="queryParams.status" placeholder="请选择状态" clearable style="width: 120px" :options="statusOptions" />
        </n-form-item>
        <n-form-item label="部门">
          <n-tree-select
            v-model:value="queryParams.deptId"
            :options="deptTreeOptions"
            placeholder="请选择部门"
            clearable
            style="width: 200px"
          />
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
        <n-form-item label="部门">
          <n-tree-select
            v-model:value="form.deptId"
            :options="deptTreeOptions"
            placeholder="请选择部门"
            clearable
          />
        </n-form-item>
        <n-form-item label="角色">
          <n-select
            v-model:value="form.roleIds"
            :options="roleOptions"
            placeholder="请选择角色"
            multiple
            clearable
          />
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
import { NButton, NSpace, NSwitch, NPopconfirm } from 'naive-ui'
import { getUserList, addUser, updateUser, changeUserStatus, resetUserPwd, deleteUser, getUserRoleIds } from '@/api/user'
import { getDeptList } from '@/api/dept'
import { getRoleList } from '@/api/role'

const statusOptions = [
  { label: '启用', value: 1 },
  { label: '禁用', value: 0 }
]

const loading = ref(false)
const tableData = ref([])
const queryParams = ref({
  pageNum: 1,
  pageSize: 10,
  username: '',
  phone: '',
  status: null,
  deptId: null
})
const pagination = ref({
  page: 1,
  pageSize: 10,
  itemCount: 0,
  showSizePicker: true,
  pageSizes: [10, 20, 50]
})

const deptTreeOptions = ref([])
const roleOptions = ref([])

const columns = [
  { title: '用户名', key: 'username' },
  { title: '昵称', key: 'nickname' },
  { title: '手机号', key: 'phone' },
  { title: '邮箱', key: 'email' },
  { title: '部门', key: 'deptName' },
  { title: '角色', key: 'roleNames', render(row) {
    return row.roleNames ? row.roleNames.join(', ') : ''
  }},
  { title: '状态', key: 'status', render(row) {
    return h(NSwitch, {
      value: row.status === 1,
      onUpdateValue: (val) => handleStatusChange(row, val)
    })
  }},
  { title: '创建时间', key: 'createTime' },
  { title: '操作', key: 'actions', render(row) {
    return h(NSpace, null, {
      default: () => [
        h(NButton, { size: 'small', onClick: () => handleEdit(row) }, { default: () => '编辑' }),
        h(NPopconfirm, { onPositiveClick: () => handleResetPwd(row) }, {
          trigger: () => h(NButton, { size: 'small', type: 'warning' }, { default: () => '重置密码' }),
          default: () => '确定重置密码为 123456 吗？'
        }),
        h(NPopconfirm, { onPositiveClick: () => handleDelete(row) }, {
          trigger: () => h(NButton, { size: 'small', type: 'error' }, { default: () => '删除' }),
          default: () => '确定删除该用户吗？'
        })
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
  deptId: null,
  roleIds: [],
  status: 1
})

onMounted(() => {
  loadData()
  loadDeptList()
  loadRoleList()
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

async function loadDeptList() {
  const res = await getDeptList()
  deptTreeOptions.value = buildTreeOptions(res)
}

function buildTreeOptions(data) {
  const map = {}
  const result = []
  data.forEach(item => {
    map[item.id] = { label: item.deptName, key: item.id, value: item.id, children: [] }
  })
  data.forEach(item => {
    if (item.parentId == 0) {
      result.push(map[item.id])
    } else if (map[item.parentId]) {
      map[item.parentId].children.push(map[item.id])
    }
  })
  return result
}

async function loadRoleList() {
  const res = await getRoleList({ pageNum: 1, pageSize: 999 })
  roleOptions.value = res.list.map(item => ({ label: item.roleName, value: item.id }))
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
    status: null,
    deptId: null
  }
  loadData()
}

function handlePageChange(page) {
  queryParams.value.pageNum = page
  loadData()
}

async function handleStatusChange(row, val) {
  const status = val ? 1 : 0
  await changeUserStatus({ id: row.id, status })
  row.status = status
  window.$message.success('操作成功')
}

async function handleResetPwd(row) {
  await resetUserPwd({ id: row.id })
  window.$message.success('密码已重置为 123456')
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
    deptId: null,
    roleIds: [],
    status: 1
  }
  showModal.value = true
}

async function handleEdit(row) {
  isEdit.value = true
  modalTitle.value = '编辑用户'
  form.value = { 
    id: row.id,
    username: row.username,
    nickname: row.nickname || '',
    phone: row.phone || '',
    email: row.email || '',
    deptId: row.deptId || null,
    roleIds: [],
    status: row.status
  }
  const roleIdsRes = await getUserRoleIds({ userId: row.id })
  form.value.roleIds = (roleIdsRes || []).map(String)
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
