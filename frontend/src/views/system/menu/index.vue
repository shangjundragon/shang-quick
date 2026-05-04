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
      />
    </n-space>
  </div>
</template>

<script setup>
import { ref, onMounted, h } from 'vue'
import { NButton, NSpace } from 'naive-ui'
import { getMenuList, addMenu, updateMenu, deleteMenu } from '@/api/menu'

const treeData = ref([])

onMounted(() => {
  loadData()
})

async function loadData() {
  const res = await getMenuList()
  treeData.value = buildTree(res)
}

function buildTree(data) {
  const map = {}
  const result = []
  data.forEach(item => {
    map[item.id] = { ...item, label: item.menuName, key: item.id, children: [] }
  })
  data.forEach(item => {
    if (item.parentId === 0) {
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
      h(NButton, { size: 'small', type: 'error', onClick: () => handleDelete(option) }, { default: () => '删除' })
    ]
  })
}

function handleAdd() {}
function handleEdit(row) {}
async function handleDelete(row) {
  await deleteMenu({ id: row.id })
  window.$message.success('删除成功')
  loadData()
}
</script>
