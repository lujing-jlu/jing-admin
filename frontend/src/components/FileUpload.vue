<template>
  <div class="file-upload-container">
    <!-- 拖拽上传区域 -->
    <div 
      v-if="!hideDropZone"
      :class="['upload-drop-zone', { 'dragover': isDragOver }]"
      @drop="handleDrop"
      @dragover.prevent="isDragOver = true"
      @dragleave="isDragOver = false"
      @click="triggerFileSelect"
    >
      <el-icon :size="48" class="upload-icon">
        <UploadFilled v-if="category === 'avatar'" />
        <Document v-else-if="category === 'document'" />
        <Picture v-else-if="category === 'image'" />
        <Folder v-else />
      </el-icon>
      <div class="upload-text">
        <p class="upload-title">{{ getUploadTitle() }}</p>
        <p class="upload-desc">{{ getUploadDesc() }}</p>
      </div>
      
      <!-- 隐藏的文件输入 -->
      <input
        ref="fileInput"
        type="file"
        :accept="getAcceptTypes()"
        :multiple="allowMultiple"
        @change="handleFileSelect"
        style="display: none"
      />
    </div>

    <!-- 文件选择按钮 -->
    <div v-if="showButton" class="upload-button-container">
      <el-button
        type="primary"
        :icon="Upload"
        :loading="uploading"
        @click="triggerFileSelect"
        :size="buttonSize"
      >
        {{ buttonText }}
      </el-button>
    </div>

    <!-- 文件列表 -->
    <div v-if="fileList.length > 0" class="file-list">
      <div 
        v-for="(file, index) in fileList" 
        :key="index"
        class="file-item"
      >
        <div class="file-info">
          <el-icon class="file-icon">
            <Picture v-if="isImage(file.name)" />
            <Document v-else />
          </el-icon>
          <div class="file-details">
            <div class="file-name">{{ file.name }}</div>
            <div class="file-size">{{ formatFileSize(file.size) }}</div>
          </div>
        </div>
        
        <div class="file-actions">
          <!-- 上传进度 -->
          <el-progress
            v-if="file.status === 'uploading'"
            :percentage="file.progress"
            :stroke-width="4"
            :show-text="false"
            status="success"
          />
          
          <!-- 上传状态 -->
          <el-tag
            v-else-if="file.status === 'success'"
            type="success"
            size="small"
          >
            已上传
          </el-tag>
          
          <el-tag
            v-else-if="file.status === 'error'"
            type="danger"
            size="small"
          >
            上传失败
          </el-tag>
          
          <!-- 删除按钮 -->
          <el-button
            v-if="file.status !== 'uploading'"
            link
            type="danger"
            :icon="Delete"
            @click="removeFile(index)"
            size="small"
          />
        </div>
      </div>
    </div>

    <!-- 预览对话框 -->
    <el-dialog
      v-model="previewVisible"
      title="文件预览"
      width="60%"
      center
    >
      <div class="preview-container">
        <img
          v-if="previewType === 'image'"
          :src="previewUrl"
          alt="预览图片"
          style="max-width: 100%; max-height: 500px;"
        />
        <div v-else class="preview-placeholder">
          <el-icon :size="64"><Document /></el-icon>
          <p>无法预览此类型文件</p>
        </div>
      </div>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { ElMessage } from 'element-plus'
import {
  Upload,
  UploadFilled,
  Document,
  Picture,
  Folder,
  Delete
} from '@element-plus/icons-vue'
import { useAuthStore } from '@/stores/auth'

// 类型定义
interface FileItem {
  name: string
  size: number
  file: File
  status: 'pending' | 'uploading' | 'success' | 'error'
  progress: number
  url?: string
  id?: number
}

interface Props {
  category?: 'avatar' | 'document' | 'image' | 'general'
  allowMultiple?: boolean
  maxSize?: number // MB
  allowedTypes?: string[]
  hideDropZone?: boolean
  showButton?: boolean
  buttonText?: string
  buttonSize?: 'large' | 'default' | 'small'
  autoUpload?: boolean
}

// Props
const props = withDefaults(defineProps<Props>(), {
  category: 'general',
  allowMultiple: false,
  maxSize: 10,
  allowedTypes: () => [],
  hideDropZone: false,
  showButton: true,
  buttonText: '选择文件',
  buttonSize: 'default',
  autoUpload: true
})

// Emits
const emit = defineEmits<{
  success: [file: any]
  error: [error: string]
  progress: [progress: number]
}>()

// 数据
const authStore = useAuthStore()
const fileInput = ref<HTMLInputElement>()
const fileList = ref<FileItem[]>([])
const isDragOver = ref(false)
const uploading = ref(false)
const previewVisible = ref(false)
const previewUrl = ref('')
const previewType = ref<'image' | 'other'>('other')

// 计算属性
const getDefaultAllowedTypes = () => {
  switch (props.category) {
    case 'avatar':
      return ['jpg', 'jpeg', 'png', 'gif']
    case 'document':
      return ['pdf', 'doc', 'docx', 'xls', 'xlsx', 'txt']
    case 'image':
      return ['jpg', 'jpeg', 'png', 'gif', 'bmp', 'webp']
    default:
      return ['jpg', 'jpeg', 'png', 'gif', 'pdf', 'doc', 'docx', 'xls', 'xlsx']
  }
}

const allowedTypes = computed(() => {
  return props.allowedTypes.length > 0 ? props.allowedTypes : getDefaultAllowedTypes()
})

// 方法
const triggerFileSelect = () => {
  fileInput.value?.click()
}

const getUploadTitle = () => {
  switch (props.category) {
    case 'avatar':
      return '上传头像'
    case 'document':
      return '上传文档'
    case 'image':
      return '上传图片'
    default:
      return '上传文件'
  }
}

const getUploadDesc = () => {
  const types = allowedTypes.value.join(', ').toUpperCase()
  return `支持 ${types} 格式，文件大小不超过 ${props.maxSize}MB`
}

const getAcceptTypes = () => {
  const mimeTypes = allowedTypes.value.map(type => {
    switch (type.toLowerCase()) {
      case 'jpg':
      case 'jpeg':
        return 'image/jpeg'
      case 'png':
        return 'image/png'
      case 'gif':
        return 'image/gif'
      case 'pdf':
        return 'application/pdf'
      case 'doc':
        return 'application/msword'
      case 'docx':
        return 'application/vnd.openxmlformats-officedocument.wordprocessingml.document'
      case 'xls':
        return 'application/vnd.ms-excel'
      case 'xlsx':
        return 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet'
      default:
        return ''
    }
  }).filter(Boolean)
  
  return mimeTypes.join(',')
}

const validateFile = (file: File): string | null => {
  // 检查文件大小
  const maxSizeBytes = props.maxSize * 1024 * 1024
  if (file.size > maxSizeBytes) {
    return `文件大小不能超过 ${props.maxSize}MB`
  }

  // 检查文件类型
  const fileExt = file.name.split('.').pop()?.toLowerCase()
  if (!fileExt || !allowedTypes.value.includes(fileExt)) {
    return `不支持的文件类型，仅支持: ${allowedTypes.value.join(', ')}`
  }

  return null
}

const handleFileSelect = (event: Event) => {
  const input = event.target as HTMLInputElement
  if (input.files) {
    handleFiles(Array.from(input.files))
  }
}

const handleDrop = (event: DragEvent) => {
  event.preventDefault()
  isDragOver.value = false
  
  if (event.dataTransfer?.files) {
    handleFiles(Array.from(event.dataTransfer.files))
  }
}

const handleFiles = (files: File[]) => {
  if (!props.allowMultiple && files.length > 1) {
    ElMessage.warning('只能选择一个文件')
    return
  }

  const validFiles: FileItem[] = []

  for (const file of files) {
    const error = validateFile(file)
    if (error) {
      ElMessage.error(`${file.name}: ${error}`)
      continue
    }

    validFiles.push({
      name: file.name,
      size: file.size,
      file,
      status: 'pending',
      progress: 0
    })
  }

  if (validFiles.length === 0) return

  // 如果不允许多文件，替换现有文件
  if (!props.allowMultiple) {
    fileList.value = validFiles
  } else {
    fileList.value.push(...validFiles)
  }

  // 自动上传
  if (props.autoUpload) {
    uploadFiles()
  }
}

const uploadFiles = async () => {
  if (uploading.value) return

  uploading.value = true

  for (const fileItem of fileList.value) {
    if (fileItem.status !== 'pending') continue

    await uploadSingleFile(fileItem)
  }

  uploading.value = false
}

const uploadSingleFile = async (fileItem: FileItem) => {
  fileItem.status = 'uploading'
  fileItem.progress = 0

  const formData = new FormData()
  formData.append('file', fileItem.file)
  formData.append('category', props.category)

  try {
    const xhr = new XMLHttpRequest()

    // 上传进度监听
    xhr.upload.addEventListener('progress', (e) => {
      if (e.lengthComputable) {
        fileItem.progress = Math.round((e.loaded / e.total) * 100)
        emit('progress', fileItem.progress)
      }
    })

    // 完成监听
    xhr.addEventListener('load', () => {
      if (xhr.status === 200) {
        const response = JSON.parse(xhr.responseText)
        fileItem.status = 'success'
        fileItem.url = response.data.url
        fileItem.id = response.data.file.id
        emit('success', response.data.file)
      } else {
        fileItem.status = 'error'
        const errorData = JSON.parse(xhr.responseText)
        emit('error', errorData.message || '上传失败')
        ElMessage.error(`${fileItem.name} 上传失败`)
      }
    })

    // 错误监听
    xhr.addEventListener('error', () => {
      fileItem.status = 'error'
      emit('error', '网络错误')
      ElMessage.error(`${fileItem.name} 上传失败`)
    })

    xhr.open('POST', '/api/files/upload')
    xhr.setRequestHeader('Authorization', `Bearer ${authStore.token}`)
    xhr.send(formData)

  } catch (error) {
    fileItem.status = 'error'
    console.error('Upload error:', error)
    emit('error', '上传失败')
    ElMessage.error(`${fileItem.name} 上传失败`)
  }
}

const removeFile = (index: number) => {
  fileList.value.splice(index, 1)
}

const isImage = (fileName: string): boolean => {
  const imageExts = ['jpg', 'jpeg', 'png', 'gif', 'bmp', 'webp']
  const ext = fileName.split('.').pop()?.toLowerCase()
  return ext ? imageExts.includes(ext) : false
}

const formatFileSize = (bytes: number): string => {
  if (bytes === 0) return '0 B'
  
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

// 暴露方法给父组件
defineExpose({
  uploadFiles,
  clearFiles: () => {
    fileList.value = []
  },
  getFiles: () => fileList.value
})
</script>

<style scoped>
.file-upload-container {
  width: 100%;
}

.upload-drop-zone {
  border: 2px dashed var(--border-secondary);
  border-radius: var(--radius-lg);
  padding: var(--spacing-xl);
  text-align: center;
  cursor: pointer;
  transition: all var(--duration-normal) ease;
  background: var(--bg-secondary);
  min-height: 200px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: var(--spacing-lg);
}

.upload-drop-zone:hover,
.upload-drop-zone.dragover {
  border-color: var(--primary-color);
  background: var(--primary-light);
}

.upload-icon {
  color: var(--text-tertiary);
  transition: color var(--duration-normal) ease;
}

.upload-drop-zone:hover .upload-icon,
.upload-drop-zone.dragover .upload-icon {
  color: var(--primary-color);
}

.upload-text {
  color: var(--text-secondary);
}

.upload-title {
  font-size: var(--font-lg);
  font-weight: var(--font-medium);
  margin-bottom: var(--spacing-xs);
  color: var(--text-primary);
}

.upload-desc {
  font-size: var(--font-sm);
  margin: 0;
}

.upload-button-container {
  margin-top: var(--spacing-md);
  text-align: center;
}

.file-list {
  margin-top: var(--spacing-lg);
  border: 1px solid var(--border-secondary);
  border-radius: var(--radius-md);
  overflow: hidden;
}

.file-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: var(--spacing-md) var(--spacing-lg);
  border-bottom: 1px solid var(--border-secondary);
  background: var(--bg-card);
  transition: background var(--duration-normal) ease;
}

.file-item:last-child {
  border-bottom: none;
}

.file-item:hover {
  background: var(--bg-hover);
}

.file-info {
  display: flex;
  align-items: center;
  gap: var(--spacing-md);
  flex: 1;
}

.file-icon {
  color: var(--text-tertiary);
  font-size: 20px;
}

.file-details {
  flex: 1;
}

.file-name {
  font-weight: var(--font-medium);
  color: var(--text-primary);
  margin-bottom: var(--spacing-xs);
}

.file-size {
  font-size: var(--font-sm);
  color: var(--text-secondary);
}

.file-actions {
  display: flex;
  align-items: center;
  gap: var(--spacing-md);
}

.preview-container {
  text-align: center;
}

.preview-placeholder {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: var(--spacing-md);
  padding: var(--spacing-xl);
  color: var(--text-tertiary);
}

/* 响应式调整 */
@media (max-width: 768px) {
  .upload-drop-zone {
    min-height: 150px;
    padding: var(--spacing-lg);
  }

  .file-item {
    flex-direction: column;
    align-items: flex-start;
    gap: var(--spacing-md);
  }

  .file-actions {
    width: 100%;
    justify-content: space-between;
  }
}
</style> 