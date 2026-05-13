<template>
  <div>
    <n-space vertical>
      <n-space>
        <n-button v-permission="['menu:add']" type="primary" @click="handleAdd">新增</n-button>
      </n-space>
      <n-tree
        :data="treeData"
        :render-suffix="renderSuffix"
        block-line
        default-expand-all
        :pattern="filterPattern"
      />
    </n-space>

    <n-modal v-model:show="showModal" :title="modalTitle" preset="card" style="width: 600px">
      <n-form ref="formRef" :model="form" :rules="rules" label-placement="left" label-width="auto">
        <n-form-item label="上级菜单" path="parentId">
          <n-tree-select
            v-model:value="form.parentId"
            :options="treeSelectOptions"
            placeholder="请选择上级菜单"
            clearable
            :default-value="0"
          />
        </n-form-item>
        <n-form-item label="菜单类型" path="menuType">
          <n-radio-group v-model:value="form.menuType">
            <n-radio :value="0">目录</n-radio>
            <n-radio :value="1">菜单</n-radio>
            <n-radio :value="2">按钮</n-radio>
          </n-radio-group>
        </n-form-item>
        <n-form-item label="菜单名称" path="menuName" required>
          <n-input v-model:value="form.menuName" placeholder="请输入菜单名称" />
        </n-form-item>
        <n-form-item v-if="form.menuType !== 2" label="图标" path="icon">
          <n-input v-model:value="form.icon" placeholder="请输入图标名称" />
        </n-form-item>
        <n-form-item v-if="form.menuType === 1" label="路由地址" path="path">
          <n-input v-model:value="form.path" placeholder="请输入路由地址" />
        </n-form-item>
        <n-form-item v-if="form.menuType === 1" label="组件路径" path="component">
          <n-input v-model:value="form.component" placeholder="请输入组件路径" />
        </n-form-item>
        <n-form-item v-if="form.menuType === 2" label="权限标识" path="perm">
          <n-input v-model:value="form.perm" placeholder="例如：user:list" />
        </n-form-item>
        <n-form-item label="显示排序" path="orderNum">
          <n-input-number v-model:value="form.orderNum" :min="0" style="width: 100%" />
        </n-form-item>
        <n-form-item v-if="form.menuType === 1" label="是否外链" path="isFrame">
          <n-radio-group v-model:value="form.isFrame">
            <n-radio :value="0">否</n-radio>
            <n-radio :value="1">是</n-radio>
          </n-radio-group>
        </n-form-item>
        <n-form-item v-if="form.menuType === 1" label="是否缓存" path="isCache">
          <n-radio-group v-model:value="form.isCache">
            <n-radio :value="0">不缓存</n-radio>
            <n-radio :value="1">缓存</n-radio>
          </n-radio-group>
        </n-form-item>
        <n-form-item v-if="form.menuType !== 2" label="是否可见" path="isVisible">
          <n-radio-group v-model:value="form.isVisible">
            <n-radio :value="1">显示</n-radio>
            <n-radio :value="0">隐藏</n-radio>
          </n-radio-group>
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
  </div>
</template>

<script setup>
import { ref, onMounted, h, useTemplateRef } from 'vue'
import { NButton, NSpace, NPopconfirm } from 'naive-ui'
import { getMenuList, addMenu, updateMenu, deleteMenu } from '@/api/menu'

const treeData = ref([])
const treeSelectOptions = ref([])
const filterPattern = ref('')
const showModal = ref(false)
const modalTitle = ref('')
const isEdit = ref(false)
const submitLoading = ref(false)
const formRef = useTemplateRef('formRef')
const form = ref({
  id: null,
  parentId: 0,
  menuName: '',
  menuType: 0,
  icon: '',
  path: '',
  component: '',
  perm: '',
  orderNum: 0,
  isFrame: 0,
  isCache: 0,
  isVisible: 1,
  status: 1
})

const rules = {
  menuName: { required: true, message: '请输入菜单名称', trigger: 'blur' }
}

onMounted(() => {
  loadData()
})

async function loadData() {
  const res = await getMenuList()
  const tree = buildTree(res)
  treeData.value = tree
  treeSelectOptions.value = buildTreeSelectOptions(res)
}

function buildTree(data) {
  const map = {}
  const result = []
  data.forEach(item => {
    map[item.id] = { ...item, label: item.menuName, key: item.id, children: [] }
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
    map[item.id] = { label: item.menuName, key: item.id, value: item.id, children: [] }
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
        default: () => '确定删除该菜单及其子菜单吗？'
      })
    ]
  })
}

function handleAdd() {
  isEdit.value = false
  modalTitle.value = '新增菜单'
  form.value = {
    id: null,
  parentId: 0,
    menuName: '',
    menuType: 0,
    icon: '',
    path: '',
    component: '',
    perm: '',
    orderNum: 0,
    isFrame: 0,
    isCache: 0,
    isVisible: 1,
    status: 1
  }
  showModal.value = true
}

function handleEdit(row) {
  isEdit.value = true
  modalTitle.value = '编辑菜单'
  form.value = {
    id: row.id,
    parentId: row.parentId || 0,
    menuName: row.menuName || '',
    menuType: row.menuType ?? 0,
    icon: row.icon || '',
    path: row.path || '',
    component: row.component || '',
    perm: row.perm || '',
    orderNum: row.orderNum || 0,
    isFrame: row.isFrame ?? 0,
    isCache: row.isCache ?? 0,
    isVisible: row.isVisible ?? 1,
    status: row.status
  }
  treeSelectOptions.value = buildTreeSelectOptions(treeData.value, row.id)
  showModal.value = true
}

async function handleSubmit() {
  try {
    await formRef.value?.validate()
    submitLoading.value = true
    if (isEdit.value) {
      await updateMenu(form.value)
    } else {
      await addMenu(form.value)
    }
    window.$message.success('操作成功')
    showModal.value = false
    loadData()
  } finally {
    submitLoading.value = false
  }
}

async function handleDelete(row) {
  await deleteMenu({ id: row.id })
  window.$message.success('删除成功')
  loadData()
}
</script>
