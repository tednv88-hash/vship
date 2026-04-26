<template>
  <view class="verify-page">
    <!-- Status banner (if already verified or pending) -->
    <view v-if="verifyStatus" class="status-banner" :class="`status-${verifyStatus}`">
      <text class="status-icon">
        {{ verifyStatus === 'verified' ? '✓' : verifyStatus === 'pending' ? '◷' : '✗' }}
      </text>
      <view class="status-info">
        <text class="status-title">{{ statusLabels[verifyStatus] }}</text>
        <text class="status-desc">{{ statusDescs[verifyStatus] }}</text>
      </view>
    </view>

    <!-- Form -->
    <view class="form-section">
      <view class="form-item">
        <text class="label">真實姓名</text>
        <input
          v-model="form.real_name"
          placeholder="請輸入真實姓名"
          class="input"
          :disabled="verifyStatus === 'verified' || verifyStatus === 'pending'"
        />
      </view>

      <view class="form-item">
        <text class="label">身份證號碼</text>
        <input
          v-model="form.id_number"
          placeholder="請輸入身份證號碼"
          class="input"
          maxlength="18"
          :disabled="verifyStatus === 'verified' || verifyStatus === 'pending'"
        />
      </view>

      <!-- ID card front -->
      <view class="form-item">
        <text class="label">身份證正面</text>
        <view
          class="upload-area"
          @tap="chooseImage('front')"
          v-if="verifyStatus !== 'verified' && verifyStatus !== 'pending'"
        >
          <image
            v-if="form.front_image"
            :src="form.front_image"
            class="preview-image"
            mode="aspectFill"
          />
          <view v-else class="upload-placeholder">
            <text class="upload-icon">+</text>
            <text class="upload-hint">上傳正面照片</text>
          </view>
        </view>
        <image
          v-else-if="form.front_image"
          :src="form.front_image"
          class="preview-image static"
          mode="aspectFill"
        />
      </view>

      <!-- ID card back -->
      <view class="form-item">
        <text class="label">身份證反面</text>
        <view
          class="upload-area"
          @tap="chooseImage('back')"
          v-if="verifyStatus !== 'verified' && verifyStatus !== 'pending'"
        >
          <image
            v-if="form.back_image"
            :src="form.back_image"
            class="preview-image"
            mode="aspectFill"
          />
          <view v-else class="upload-placeholder">
            <text class="upload-icon">+</text>
            <text class="upload-hint">上傳反面照片</text>
          </view>
        </view>
        <image
          v-else-if="form.back_image"
          :src="form.back_image"
          class="preview-image static"
          mode="aspectFill"
        />
      </view>
    </view>

    <!-- Submit button -->
    <view
      v-if="verifyStatus !== 'verified' && verifyStatus !== 'pending'"
      class="submit-btn"
      @tap="handleSubmit"
    >
      <text class="submit-btn-text">{{ t('common.submit') }}</text>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { t } from '@/locale'
import { userApi } from '@/api/user'

type VerifyStatus = 'verified' | 'pending' | 'rejected' | null

const verifyStatus = ref<VerifyStatus>(null)
const loading = ref(false)

const form = reactive({
  real_name: '',
  id_number: '',
  front_image: '',
  back_image: '',
})

const statusLabels: Record<string, string> = {
  verified: '已認證',
  pending: '審核中',
  rejected: '認證失敗',
}

const statusDescs: Record<string, string> = {
  verified: '您的身份已通過認證',
  pending: '您的認證資料正在審核中，請耐心等待',
  rejected: '認證失敗，請重新提交資料',
}

async function loadStatus() {
  try {
    const res = (await userApi.getIdentityStatus()) as any
    const data = res?.data || res
    if (data) {
      verifyStatus.value = data.status || null
      if (data.real_name) form.real_name = data.real_name
      if (data.id_number) form.id_number = data.id_number
      if (data.front_image) form.front_image = data.front_image
      if (data.back_image) form.back_image = data.back_image
    }
  } catch {
    // Not verified yet
  }
}

function chooseImage(side: 'front' | 'back') {
  uni.chooseImage({
    count: 1,
    sizeType: ['compressed'],
    sourceType: ['album', 'camera'],
    success: (res) => {
      const tempPath = res.tempFilePaths[0]
      if (side === 'front') {
        form.front_image = tempPath
      } else {
        form.back_image = tempPath
      }
    },
  })
}

async function handleSubmit() {
  if (!form.real_name) {
    uni.showToast({ title: '請輸入真實姓名', icon: 'none' })
    return
  }
  if (!form.id_number) {
    uni.showToast({ title: '請輸入身份證號碼', icon: 'none' })
    return
  }
  if (!form.front_image) {
    uni.showToast({ title: '請上傳身份證正面', icon: 'none' })
    return
  }
  if (!form.back_image) {
    uni.showToast({ title: '請上傳身份證反面', icon: 'none' })
    return
  }

  if (loading.value) return
  loading.value = true

  try {
    await userApi.verifyIdentity({
      real_name: form.real_name,
      id_number: form.id_number,
      front_image: form.front_image,
      back_image: form.back_image,
    })
    verifyStatus.value = 'pending'
    uni.showToast({ title: '提交成功', icon: 'success' })
  } catch {
    uni.showToast({ title: '提交失敗', icon: 'none' })
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadStatus()
})
</script>

<style scoped>
.verify-page {
  min-height: 100vh;
  background: #f7f8fa;
  padding: 24rpx;
}

.status-banner {
  display: flex;
  align-items: center;
  padding: 30rpx;
  border-radius: 16rpx;
  margin-bottom: 24rpx;
}

.status-verified {
  background: rgba(7, 193, 96, 0.1);
}

.status-pending {
  background: rgba(255, 165, 0, 0.1);
}

.status-rejected {
  background: rgba(231, 76, 60, 0.1);
}

.status-icon {
  font-size: 48rpx;
  margin-right: 24rpx;
  width: 60rpx;
  text-align: center;
}

.status-verified .status-icon {
  color: #07c160;
}

.status-pending .status-icon {
  color: #ffa500;
}

.status-rejected .status-icon {
  color: #e74c3c;
}

.status-info {
  flex: 1;
}

.status-title {
  display: block;
  font-size: 30rpx;
  font-weight: 600;
  color: #333;
  margin-bottom: 6rpx;
}

.status-desc {
  font-size: 24rpx;
  color: #666;
}

.form-section {
  background: #fff;
  border-radius: 16rpx;
  padding: 10rpx 30rpx;
}

.form-item {
  padding: 24rpx 0;
  border-bottom: 1rpx solid #f0f0f0;
}

.form-item:last-child {
  border-bottom: none;
}

.label {
  display: block;
  font-size: 28rpx;
  color: #333;
  font-weight: 500;
  margin-bottom: 16rpx;
}

.input {
  width: 100%;
  height: 80rpx;
  background: #f7f8fa;
  border-radius: 12rpx;
  padding: 0 24rpx;
  font-size: 28rpx;
  color: #333;
  box-sizing: border-box;
}

.upload-area {
  width: 100%;
  height: 300rpx;
  background: #f7f8fa;
  border-radius: 12rpx;
  overflow: hidden;
  display: flex;
  align-items: center;
  justify-content: center;
  border: 2rpx dashed #ddd;
}

.upload-placeholder {
  display: flex;
  flex-direction: column;
  align-items: center;
}

.upload-icon {
  font-size: 72rpx;
  color: #ccc;
  line-height: 1;
}

.upload-hint {
  font-size: 24rpx;
  color: #999;
  margin-top: 12rpx;
}

.preview-image {
  width: 100%;
  height: 300rpx;
  border-radius: 12rpx;
}

.preview-image.static {
  display: block;
}

.submit-btn {
  width: 100%;
  height: 96rpx;
  background: #0f3a57;
  border-radius: 16rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-top: 40rpx;
}

.submit-btn-text {
  font-size: 32rpx;
  color: #fff;
  font-weight: 600;
}
</style>
