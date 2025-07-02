<template>
  <el-dialog
    :model-value="visible"
    :title="mode === 'create' ? '新增角色' : '编辑角色'"
    width="500px"
    @update:model-value="handleClose"
    @close="handleClose"
  >
    <el-form
      ref="formRef"
      :model="formData"
      :rules="formRules"
      label-width="80px"
    >
      <el-form-item label="角色标识" prop="name">
        <el-input
          v-model="formData.name"
          placeholder="请输入角色标识（如：editor）"
          :disabled="mode === 'edit'"
          clearable
        />
      </el-form-item>
      
      <el-form-item label="角色名称" prop="display_name">
        <el-input
          v-model="formData.display_name"
          placeholder="请输入角色名称（如：编辑员）"
          clearable
        />
      </el-form-item>
      
      <el-form-item label="描述" prop="description">
        <el-input
          v-model="formData.description"
          type="textarea"
          :rows="3"
          placeholder="请输入角色描述"
          clearable
        />
      </el-form-item>
      
      <el-form-item label="状态" prop="status">
        <el-switch
          v-model="formData.status"
          active-text="启用"
          inactive-text="禁用"
        />
      </el-form-item>
    </el-form>
    
    <template #footer>
      <span class="dialog-footer">
        <el-button @click="handleClose">取消</el-button>
        <el-button type="primary" :loading="loading" @click="handleSubmit">
          {{ mode === 'create' ? '创建' : '更新' }}
        </el-button>
      </span>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, reactive, watch, type PropType } from 'vue'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
import type { Role } from '../stores/role'

// Props
const props = defineProps({
  visible: {
    type: Boolean,
    default: false
  },
  mode: {
    type: String as PropType<'create' | 'edit'>,
    default: 'create'
  },
  roleData: {
    type: Object as PropType<Role | null>,
    default: null
  },
  loading: {
    type: Boolean,
    default: false
  }
})

// Emits
const emit = defineEmits<{
  'update:visible': [value: boolean]
  'submit': [data: any]
}>()

// Form 引用
const formRef = ref<FormInstance>()

// 表单数据
const formData = reactive({
  name: '',
  display_name: '',
  description: '',
  status: true
})

// 表单验证规则
const formRules: FormRules = {
  name: [
    { required: true, message: '请输入角色标识', trigger: 'blur' },
    { min: 2, max: 20, message: '角色标识长度在 2 到 20 个字符', trigger: 'blur' },
    { pattern: /^[a-z_]+$/, message: '角色标识只能包含小写字母和下划线', trigger: 'blur' }
  ],
  display_name: [
    { required: true, message: '请输入角色名称', trigger: 'blur' },
    { min: 2, max: 20, message: '角色名称长度在 2 到 20 个字符', trigger: 'blur' }
  ]
}

// 监听角色数据变化
watch(() => props.roleData, (newVal) => {
  if (newVal && props.mode === 'edit') {
    Object.assign(formData, {
      name: newVal.name,
      display_name: newVal.display_name,
      description: newVal.description,
      status: newVal.status
    })
  }
}, { immediate: true })

// 监听弹窗显示状态
watch(() => props.visible, (newVal) => {
  if (newVal && props.mode === 'create') {
    // 重置表单数据
    Object.assign(formData, {
      name: '',
      display_name: '',
      description: '',
      status: true
    })
    // 清除验证
    formRef.value?.clearValidate()
  }
})

// 处理关闭
function handleClose() {
  emit('update:visible', false)
}

// 处理提交
async function handleSubmit() {
  if (!formRef.value) return
  
  try {
    await formRef.value.validate()
    emit('submit', { ...formData })
  } catch (error) {
    ElMessage.error('请检查表单输入')
  }
}
</script>

<style scoped>
.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}
</style> 