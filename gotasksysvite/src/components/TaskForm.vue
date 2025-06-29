<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useTaskTypeStore } from '@/stores/taskType'
import type { CreateTaskInput } from '@/types/task'
import { ElForm, ElFormItem, ElInput, ElSelect, ElOption, ElInputNumber, ElButton } from 'element-plus'
import type { FormInstance } from 'element-plus'

// 定义组件要触发的事件
const emit = defineEmits(['submit'])

const formRef = ref<FormInstance>()
const formData = ref<CreateTaskInput>({
  title: '',
  description: '',
  priority: 'P2', // 默认优先级
  effort: 8,      // 默认工时
  task_type_id: '',
})

// 从store获取任务类型
const taskTypeStore = useTaskTypeStore()
onMounted(() => {
  taskTypeStore.fetchTaskTypes()
})

const submitForm = async (formEl: FormInstance | undefined) => {
  if (!formEl) return
  await formEl.validate((valid) => {
    if (valid) {
      emit('submit', formData.value)
    } else {
      console.log('error submit!')
    }
  })
}
</script>

<template>
  <el-form ref="formRef" :model="formData" label-width="80px">
    <el-form-item label="任务标题" prop="title" :rules="[{ required: true, message: '请输入任务标题' }]">
      <el-input v-model="formData.title" />
    </el-form-item>
    <el-form-item label="任务类型" prop="task_type_id" :rules="[{ required: true, message: '请选择任务类型' }]">
      <el-select v-model="formData.task_type_id" placeholder="请选择任务类型" style="width: 100%;">
        <el-option
          v-for="item in taskTypeStore.taskTypes"
          :key="item.id"
          :label="item.name"
          :value="item.id"
        />
      </el-select>
    </el-form-item>
    <el-form-item label="优先级" prop="priority">
      <el-select v-model="formData.priority" style="width: 100%;">
        <el-option label="P0 - 紧急" value="P0" />
        <el-option label="P1 - 高" value="P1" />
        <el-option label="P2 - 中" value="P2" />
        <el-option label="P3 - 低" value="P3" />
      </el-select>
    </el-form-item>
    <el-form-item label="预估工时" prop="effort">
      <el-input-number v-model="formData.effort" :min="1" :step="1" /> 小时
    </el-form-item>
    <el-form-item label="任务描述" prop="description">
      <el-input v-model="formData.description" type="textarea" :rows="4" />
    </el-form-item>
    <el-form-item>
      <el-button type="primary" @click="submitForm(formRef)">创建任务</el-button>
    </el-form-item>
  </el-form>
</template>