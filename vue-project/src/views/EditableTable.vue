<script setup>
import { ref, nextTick } from 'vue'

// ========== 数据：类比 C++ 的 std::vector<RowData> ==========
const tableData = ref([
  { id: 1, name: '张三', age: 28, address: '北京' },
  { id: 2, name: '李四', age: 32, address: '上海' },
  { id: 3, name: '王五', age: 25, address: '深圳' },
])

let nextId = 4
const editingRow = ref(null)
const editingCol = ref(null)

function startEdit(rowIndex, colName) {
  editingRow.value = rowIndex
  editingCol.value = colName
  nextTick(() => {
    const input = document.querySelector('.editing-cell input')
    if (input) input.focus()
  })
}

function saveEdit() {
  editingRow.value = null
  editingCol.value = null
}

function addRow() {
  tableData.value.push({ id: nextId++, name: '', age: 0, address: '' })
  startEdit(tableData.value.length - 1, 'name')
}

function deleteRow(index) {
  tableData.value.splice(index, 1)
}

function isEditing(rowIndex, colName) {
  return editingRow.value === rowIndex && editingCol.value === colName
}
</script>

<template>
  <div class="demo-container">
    <h1>Element Plus 可编辑表格 Demo</h1>
    <el-button type="primary" @click="addRow" style="margin-bottom: 16px">➕ 新增一行</el-button>

    <el-table :data="tableData" border stripe style="width: 100%">
      <el-table-column prop="id" label="ID" width="80" />

      <el-table-column prop="name" label="姓名" width="160">
        <template #default="{ row, $index }">
          <div v-if="isEditing($index, 'name')" class="editing-cell">
            <el-input v-model="row.name" size="small" @blur="saveEdit" @keyup.enter="saveEdit" />
          </div>
          <span v-else @click="startEdit($index, 'name')" class="cell-clickable">
            {{ row.name || '点击编辑' }}
          </span>
        </template>
      </el-table-column>

      <el-table-column prop="age" label="年龄" width="120">
        <template #default="{ row, $index }">
          <div v-if="isEditing($index, 'age')" class="editing-cell">
            <el-input-number v-model="row.age" size="small" :min="0" :max="150" @blur="saveEdit" @keyup.enter="saveEdit" />
          </div>
          <span v-else @click="startEdit($index, 'age')" class="cell-clickable">{{ row.age }}</span>
        </template>
      </el-table-column>

      <el-table-column prop="address" label="地址" min-width="200">
        <template #default="{ row, $index }">
          <div v-if="isEditing($index, 'address')" class="editing-cell">
            <el-input v-model="row.address" size="small" @blur="saveEdit" @keyup.enter="saveEdit" />
          </div>
          <span v-else @click="startEdit($index, 'address')" class="cell-clickable">
            {{ row.address || '点击编辑' }}
          </span>
        </template>
      </el-table-column>

      <el-table-column label="操作" width="100" fixed="right">
        <template #default="{ $index }">
          <el-popconfirm title="确定删除这一行吗？" @confirm="deleteRow($index)">
            <template #reference>
              <el-button type="danger" size="small">🗑 删除</el-button>
            </template>
          </el-popconfirm>
        </template>
      </el-table-column>
    </el-table>
  </div>
</template>

<style>
body {
  margin: 0;
  padding: 24px;
  background: #f5f7fa;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
}
</style>

<style scoped>
.demo-container { max-width: 900px; margin: 0 auto; }
h1 { font-size: 22px; margin-bottom: 16px; color: #303133; }
.cell-clickable {
  display: inline-block; width: 100%; min-height: 24px; cursor: pointer;
  padding: 4px 8px; border-radius: 4px; transition: background 0.2s; color: #606266;
}
.cell-clickable:hover { background: #ecf5ff; color: #409eff; }
.editing-cell { margin: -4px -8px; }
</style>
