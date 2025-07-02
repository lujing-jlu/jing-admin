<template>
  <el-dialog
    :model-value="visible"
    :title="isEdit ? '编辑用户' : '新建用户'"
    width="600px"
    @close="handleClose"
    destroy-on-close
  >
    <el-form
      ref="formRef"
      :model="form"
      :rules="rules"
      label-width="100px"
    >
      <el-form-item label="用户名" prop="username">
        <el-input
          v-model="form.username"
          :disabled="isEdit"
          placeholder="请输入用户名"
        />
      </el-form-item>
      
      <el-form-item label="邮箱" prop="email">
        <el-input
          v-model="form.email"
          placeholder="请输入邮箱"
          type="email"
        />
      </el-form-item>
      
      <el-form-item label="密码" prop="password" v-if="!isEdit">
        <el-input
          v-model="form.password"
          type="password"
          placeholder="请输入密码"
          show-password
        />
      </el-form-item>
      
      <el-form-item label="角色" prop="role">
        <el-select v-model="form.role" placeholder="请选择角色">
          <el-option label="管理员" value="admin" />
          <el-option label="普通用户" value="user" />
        </el-select>
      </el-form-item>
      
      <el-form-item label="状态" prop="status">
        <el-switch v-model="form.status" />
      </el-form-item>
    </el-form>
    
    <template #footer>
      <span class="dialog-footer">
        <el-button @click="handleClose">取消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="loading">
          {{ isEdit ? '更新' : '创建' }}
        </el-button>
      </span>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, reactive, watch, computed, type PropType } from 'vue'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
import type { User } from '../stores/user'

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
  userData: {
    type: Object as PropType<User | null>,
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
const form = reactive({
  username: '',
  email: '',
  password: '',
  role: 'user',
  status: true
})

// 计算属性
const isEdit = computed(() => props.mode === 'edit')

// 表单验证规则
const rules: FormRules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 3, max: 20, message: '用户名长度在 3 到 20 个字符', trigger: 'blur' }
  ],
  email: [
    { required: true, message: '请输入邮箱地址', trigger: 'blur' },
    { type: 'email', message: '请输入正确的邮箱地址', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, message: '密码长度不能少于 6 个字符', trigger: 'blur' }
  ],
  role: [
    { required: true, message: '请选择角色', trigger: 'change' }
  ]
}

// 监听用户数据变化
watch(() => props.userData, (newVal) => {
  if (newVal && props.mode === 'edit') {
    Object.assign(form, {
      username: newVal.username,
      email: newVal.email,
      password: '',
      role: newVal.role,
      status: newVal.status
    })
  }
}, { immediate: true })

// 监听弹窗显示状态
watch(() => props.visible, (newVal) => {
  if (newVal && props.mode === 'create') {
    // 重置表单数据
    Object.assign(form, {
      username: '',
      email: '',
      password: '',
      role: 'user',
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
    
    // 准备提交数据
    const submitData: any = { ...form }
    if (props.mode === 'edit') {
      // 编辑模式下，如果密码为空则不更新密码
      if (!submitData.password) {
        delete submitData.password
      }
    }
    
    emit('submit', submitData)
  } catch (error) {
    ElMessage.error('请检查表单输入')
  }
}
</script>

<style scoped>
/* 表单输入框样式 */

:deep(.el-form-item__label) {
  color: var(--text-secondary) !important;
}

:deep(.el-input__wrapper) {
  background: var(--bg-input) !important;
  border: 1px solid var(--border-secondary) !important;
  border-radius: var(--radius-md) !important;
}

:deep(.el-input__wrapper:hover) {
  border-color: var(--primary-color) !important;
}

:deep(.el-input__wrapper.is-focus) {
  border-color: var(--primary-color) !important;
  box-shadow: 0 0 0 2px var(--primary-color-light) !important;
}

:deep(.el-select .el-input__wrapper) {
  background: var(--bg-input) !important;
  border: 1px solid var(--border-secondary) !important;
}
</style> 