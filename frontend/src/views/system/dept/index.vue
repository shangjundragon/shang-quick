<template>
  <div>
    <n-space vertical>
      <n-space>
        <n-button v-permission="['dept:add']" type="primary" @click="handleAdd">新增</n-button>
      </n-space>
      <n-tree
        :data="treeData"
        :render-suffix="renderSuffix"
        block-line
        default-expand-all
      />
    </n-space>

    <n-modal v-model:show="showModal" :title="modalTitle" preset="card" style="width: 500px">
      <n-form :model="form" label-placement="left" label-width="auto">
        <n-form-item label="上级部门">
          <n-tree-select
            v-model:value="form.parentId"
            :options="treeSelectOptions"
            placeholder="请选择上级部门"
            clearable
            :default-value="0"
          />
        </n-form-item>
        <n-form-item label="部门名称" required>
          <n-input v-model:value="form.deptName" placeholder="请输入部门名称" />
        </n-form-item>
        <n-form-item label="显示排序">
          <n-input-number v-model:value="form.orderNum" :min="0" style="width: 100%" />
        </n-form-item>
        <n-form-item label="负责人">
          <n-input v-model:value="form.leader" placeholder="请输入负责人" />
        </n-form-item>
        <n-form-item label="联系电话">
          <n-input v-model:value="form.phone" placeholder="请输入联系电话" />
        </n-form-item>
        <n-form-item label="邮箱">
          <n-input v-model:value="form.email" placeholder="请输入邮箱" />
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
import { NButton, NSpace, NPopconfirm } from 'naive-ui'
import { getDeptList, addDept, updateDept, deleteDept } from '@/api/dept'

const treeData = ref([])
const treeSelectOptions = ref([])
const showModal = ref(false)
const modalTitle = ref('')
const isEdit = ref(false)
const submitLoading = ref(false)
const form = ref({
  id: null,
  parentId: 0,
  deptName: '',
  orderNum: 0,
  leader: '',
  phone: '',
  email: '',
  status: 1
})

onMounted(() => {
  loadData()
})

async function loadData() {
  const res = await getDeptList()
  const tree = buildTree(res)
  treeData.value = tree
  treeSelectOptions.value = buildTreeSelectOptions(res)
}

function buildTree(data) {
  const map = {}
  const result = []
  data.forEach(item => {
    map[item.id] = { ...item, label: item.deptName, key: item.id, children: [] }
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

function buildTreeSelectOptions(data, excludeId) {
  const map = {}
  const result = []
  data.forEach(item => {
    if (item.id === excludeId) return
    map[item.id] = { label: item.deptName, key: item.id, value: item.id, children: [] }
  })
  data.forEach(item => {
    if (item.id === excludeId) return
    if (item.parentId == 0) {
      result.push(map[item.id])
    } else if (map[item.parentId]) {
      map[item.parentId].children.push(map[item.id])
    }
  })
  return result
}

function renderSuffix({ option }) {
  return h(NSpace, null, {
    default: () => [
      h(NButton, { size: 'small', onClick: () => handleEdit(option) }, { default: () => '编辑' }),
      h(NPopconfirm, { onPositiveClick: () => handleDelete(option) }, {
        trigger: () => h(NButton, { size: 'small', type: 'error' }, { default: () => '删除' }),
        default: () => '确定删除该部门吗？'
      })
    ]
  })
}

function handleAdd() {
  isEdit.value = false
  modalTitle.value = '新增部门'
  form.value = {
    id: null,
    parentId: 0,
    deptName: '',
    orderNum: 0,
    leader: '',
    phone: '',
    email: '',
    status: 1
  }
  showModal.value = true
}

function handleEdit(row) {
  isEdit.value = true
  modalTitle.value = '编辑部门'
  form.value = {
    id: row.id,
    parentId: row.parentId || 0,
    deptName: row.deptName || '',
    orderNum: row.orderNum || 0,
    leader: row.leader || '',
    phone: row.phone || '',
    email: row.email || '',
    status: row.status
  }
  treeSelectOptions.value = buildTreeSelectOptions(treeData.value, row.id)
  showModal.value = true
}

async function handleSubmit() {
  submitLoading.value = true
  try {
    if (isEdit.value) {
      await updateDept(form.value)
    } else {
      await addDept(form.value)
    }
    window.$message.success('操作成功')
    showModal.value = false
    loadData()
  } finally {
    submitLoading.value = false
  }
}

async function handleDelete(row) {
  await deleteDept({ id: row.id })
  window.$message.success('删除成功')
  loadData()
}
</script>
