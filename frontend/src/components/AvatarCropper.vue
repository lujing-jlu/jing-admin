<template>
  <el-dialog
    :visible="visible"
    @update:visible="(val: boolean) => { console.log('el-dialog update:visible', val); emit('update:visible', val) }"
    title="更换头像"
    width="400px"
    :close-on-click-modal="false"
    :destroy-on-close="true"
    center
    :append-to-body="true"
  >
    <div class="avatar-cropper-container">
      <div>AvatarCropper 渲染，visible: {{ visible }} imageUrl: {{ imageUrl }}</div>
      <div v-if="!imageUrl" class="upload-area">
        <el-button type="primary" icon="Upload" @click="triggerFileInput">选择图片</el-button>
        <input ref="fileInput" type="file" accept="image/*" style="display:none" @change="onFileInputChange" />
      </div>
      <div v-if="imageUrl" class="cropper-area">
        <canvas ref="canvasRef" :width="256" :height="256" style="border-radius: 50%; background: #f8fafc;" />
      </div>
      <div v-if="imageUrl" class="cropper-actions">
        <el-button type="success" :loading="uploading" @click="uploadAvatar">上传头像</el-button>
        <el-button @click="cancelCrop">取消</el-button>
      </div>
      <el-progress v-if="uploading" :percentage="uploadProgress" style="width: 100%; margin-top: 10px;" />
    </div>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, watch, nextTick, onMounted, onUpdated } from 'vue'
import { ElMessage } from 'element-plus'

const props = defineProps({
  visible: Boolean,
  currentAvatar: String,
  imageSrc: String
})
const emit = defineEmits(['update:visible', 'success'])

const imageUrl = ref<string | null>(null)
const imageFile = ref<File | null>(null)
const uploading = ref(false)
const uploadProgress = ref(0)
const fileInput = ref<HTMLInputElement | null>(null)
const canvasRef = ref<HTMLCanvasElement | null>(null)

watch(() => props.visible, (val) => {
  if (val) {
    // 弹窗打开时显示传入图片
    if (props.imageSrc) {
      imageUrl.value = props.imageSrc
      nextTick(drawToCanvas)
    }
  } else {
    imageUrl.value = null
    imageFile.value = null
    uploadProgress.value = 0
    uploading.value = false
  }
})

// 监听imageSrc变化，动态更新imageUrl和canvas
watch(() => props.imageSrc, (val) => {
  if (props.visible && val) {
    imageUrl.value = val
    nextTick(drawToCanvas)
  }
})

function triggerFileInput() {
  fileInput.value?.click()
}

function onFileInputChange(e: Event) {
  const files = (e.target as HTMLInputElement).files
  if (!files || files.length === 0) return
  const file = files[0]
  if (!file.type.startsWith('image/')) {
    ElMessage.error('只能上传图片文件')
    return
  }
  if (file.size / 1024 / 1024 > 2) {
    ElMessage.error('图片大小不能超过2MB')
    return
  }
  imageFile.value = file
  const reader = new FileReader()
  reader.onload = (ev: ProgressEvent<FileReader>) => {
    imageUrl.value = ev.target?.result as string
    nextTick(drawToCanvas)
  }
  reader.readAsDataURL(file)
}

function drawToCanvas() {
  if (!imageUrl.value || !canvasRef.value) return
  const ctx = canvasRef.value.getContext('2d')
  if (!ctx) return
  const img = new window.Image()
  img.onload = function () {
    // 居中裁剪为正方形
    const size = Math.min(img.width, img.height)
    const sx = (img.width - size) / 2
    const sy = (img.height - size) / 2
    ctx.clearRect(0, 0, 256, 256)
    ctx.save()
    ctx.beginPath()
    ctx.arc(128, 128, 128, 0, Math.PI * 2)
    ctx.closePath()
    ctx.clip()
    ctx.drawImage(img, sx, sy, size, size, 0, 0, 256, 256)
    ctx.restore()
  }
  img.src = imageUrl.value
}

onUpdated(() => {
  if (imageUrl.value) drawToCanvas()
})

async function uploadAvatar() {
  if (!canvasRef.value) {
    ElMessage.error('请先选择图片')
    return
  }
  uploading.value = true
  uploadProgress.value = 0
  try {
    // 导出canvas为blob
    canvasRef.value.toBlob(async (blob) => {
      if (!blob) {
        ElMessage.error('裁剪失败')
        uploading.value = false
        return
      }
      console.log('canvas toBlob success', blob)
      const formData = new FormData()
      formData.append('file', blob, 'avatar.png')
      formData.append('category', 'avatar')
      const token = localStorage.getItem('token')
      console.log('准备上传', formData, token)
      const response = await fetch('http://localhost:8081/api/files/upload', {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${token}`
        },
        body: formData
      })
      const result = await response.json()
      console.log('上传响应', result)
      if (result.code === 200 && result.data && result.data.url) {
        ElMessage.success('头像上传成功')
        emit('success', result.data.url)
        emit('update:visible', false)
      } else {
        ElMessage.error(result.message || '头像上传失败')
      }
      uploading.value = false
    }, 'image/png')
  } catch (e) {
    ElMessage.error('网络错误，上传失败')
    uploading.value = false
  }
}

function cancelCrop() {
  imageUrl.value = null
  imageFile.value = null
  nextTick(() => fileInput.value?.click())
}
</script>

<style scoped>
.avatar-cropper-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 260px;
}
.upload-area {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  width: 100%;
  min-height: 180px;
}
.upload-tip {
  color: #888;
  font-size: 13px;
  margin-top: 8px;
}
.cropper-area {
  width: 256px;
  height: 256px;
  margin: 0 auto;
  background: #f8fafc;
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 2px 8px rgba(0,0,0,0.06);
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
}
.cropper-actions {
  display: flex;
  justify-content: space-between;
  width: 100%;
  margin-top: 12px;
}
</style> 