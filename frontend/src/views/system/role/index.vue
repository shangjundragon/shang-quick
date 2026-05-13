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

    <n-modal v-model:show="showModal" :title="modalTitle" preset="card" style="width: 500px">
      <n-form ref="formRef" :model="form" :rules="rules" label-placement="left" label-width="auto">
        <n-form-item label="角色名称" path="roleName" required>
          <n-input v-model:value="form.roleName" placeholder="请输入角色名称" />
        </n-form-item>
        <n-form-item label="角色编码" path="roleCode" required>
          <n-input v-model:value="form.roleCode" placeholder="请输入角色编码" />
        </n-form-item>
        <n-form-item label="备注" path="remark">
          <n-input v-model:value="form.remark" type="textarea" placeholder="请输入备注" />
        </n-form-item>
        <n-form-item label="状态" path="status">
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

    <n-modal v-model:show="showMenuModal" title="分配菜单" preset="card" style="width: 400px">
      <n-tree
        v-model:checked-keys="checkedMenuIds"
        :data="menuTreeData"
        checkable
        block-line
        cascade
        default-expand-all
      />
      <template #footer>
        <n-space>
          <n-button @click="showMenuModal = false">取消</n-button>
          <n-button type="primary" :loading="menuSubmitLoading" @click="handleAssignMenu">确定</n-button>
        </n-space>
      </template>
    </n-modal>
  </div>
</template>

<script setup>
import { ref, onMounted, h, useTemplateRef } from 'vue'
import { NButton, NSpace, NPopconfirm } from 'naive-ui'
import { getRoleList, addRole, updateRole, deleteRole, getRoleMenuIds, assignRoleMenu } from '@/api/role'
import { getMenuList } from '@/api/menu'
import { formatTimestamp } from '@/utils/format'

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
  { title: '创建时间', key: 'createTime', render(row) { return formatTimestamp(row.createTime) } },
  { title: '操作', key: 'actions', render(row) {
    return h(NSpace, null, {
      default: () => [
        h(NButton, { size: 'small', onClick: () => handleEdit(row) }, { default: () => '编辑' }),
        h(NButton, { size: 'small', onClick: () => handleAssign(row) }, { default: () => '分配菜单' }),
        h(NPopconfirm, { onPositiveClick: () => handleDelete(row) }, {
          trigger: () => h(NButton, { size: 'small', type: 'error' }, { default: () => '删除' }),
          default: () => '确定删除该角色吗？'
        })
      ]
    })
  }}
]

const showModal = ref(false)
const modalTitle = ref('')
const isEdit = ref(false)
const submitLoading = ref(false)
const formRef = useTemplateRef('formRef')
const form = ref({
  id: null,
  roleName: '',
  roleCode: '',
  remark: '',
  status: 1
})

const rules = {
  roleName: { required: true, message: '请输入角色名称', trigger: 'blur' },
  roleCode: { required: true, message: '请输入角色编码', trigger: 'blur' }
}

const showMenuModal = ref(false)
const menuSubmitLoading = ref(false)
const currentRoleId = ref(null)
const checkedMenuIds = ref([])
const menuTreeData = ref([])

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

function handleAdd() {
  isEdit.value = false
  modalTitle.value = '新增角色'
  form.value = {
    id: null,
    roleName: '',
    roleCode: '',
    remark: '',
    status: 1
  }
  showModal.value = true
}

function handleEdit(row) {
  isEdit.value = true
  modalTitle.value = '编辑角色'
  form.value = {
    id: row.id,
    roleName: row.roleName || '',
    roleCode: row.roleCode || '',
    remark: row.remark || '',
    status: row.status
  }
  showModal.value = true
}

async function handleSubmit() {
  try {
    await formRef.value?.validate()
    submitLoading.value = true
    if (isEdit.value) {
      await updateRole(form.value)
    } else {
      await addRole(form.value)
    }
    window.$message.success('操作成功')
    showModal.value = false
    loadData()
  } finally {
    submitLoading.value = false
  }
}

async function handleDelete(row) {
  await deleteRole({ id: row.id })
  window.$message.success('删除成功')
  loadData()
}

async function handleAssign(row) {
  currentRoleId.value = row.id
  const menuRes = await getMenuList()
  menuTreeData.value = buildMenuTree(menuRes)
  const ids = await getRoleMenuIds({ roleId: row.id })
  checkedMenuIds.value = ids || []
  showMenuModal.value = true
}

function buildMenuTree(data) {
  const map = {}
  const result = []
  data.forEach(item => {
    map[item.id] = { label: item.menuName, key: item.id, children: [] }
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

async function handleAssignMenu() {
  menuSubmitLoading.value = true
  try {
    await assignRoleMenu({ roleId: currentRoleId.value, menuIds: checkedMenuIds.value })
    window.$message.success('分配成功')
    showMenuModal.value = false
  } finally {
    menuSubmitLoading.value = false
  }
}
</script>
