<template>
  <view class="feedback-page">
    <!-- Type selector -->
    <view class="form-section">
      <text class="section-label">Type</text>
      <view class="type-selector">
        <view
          v-for="option in typeOptions"
          :key="option.value"
          class="type-option"
          :class="{ active: form.type === option.value }"
          @click="form.type = option.value"
        >
          <text class="type-text">{{ option.label }}</text>
        </view>
      </view>
    </view>

    <!-- Content textarea -->
    <view class="form-section">
      <text class="section-label">Content</text>
      <textarea
        class="content-textarea"
        v-model="form.content"
        placeholder="Please describe your feedback..."
        :maxlength="1000"
      />
      <text class="char-count">{{ form.content.length }}/1000</text>
    </view>

    <!-- Image upload -->
    <view class="form-section">
      <text class="section-label">Images (max 3)</text>
      <view class="photo-list">
        <view
          v-for="(photo, idx) in form.images"
          :key="idx"
          class="photo-item"
        >
          <image class="photo-img" :src="photo" mode="aspectFill" />
          <view class="photo-delete" @click="removeImage(idx)">
            <text class="delete-icon">x</text>
          </view>
        </view>
        <view
          v-if="form.images.length < 3"
          class="photo-add"
          @click="chooseImage"
        >
          <text class="add-icon">+</text>
        </view>
      </view>
    </view>

    <!-- Contact info -->
    <view class="form-section">
      <text class="section-label">Contact (optional)</text>
      <input
        class="contact-input"
        v-model="form.contact"
        placeholder="Phone or email"
      />
    </view>

    <!-- Submit -->
    <view class="submit-section">
      <view
        class="submit-btn"
        :class="{ disabled: submitting || !canSubmit }"
        @click="onSubmit"
      >
        <text class="submit-text">
          {{ submitting ? t('common.loading') : t('common.submit') }}
        </text>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { t } from '@/locale'
import { commonApi } from '@/api/common'

interface FeedbackForm {
  type: string
  content: string
  images: string[]
  contact: string
}

const typeOptions = [
  { label: '建議', value: 'suggestion' },
  { label: '投訴', value: 'complaint' },
  { label: '其他', value: 'other' },
]

const form = ref<FeedbackForm>({
  type: 'suggestion',
  content: '',
  images: [],
  contact: '',
})
const submitting = ref(false)

const canSubmit = computed(() => {
  return form.value.content.trim().length > 0
})

function chooseImage() {
  const remaining = 3 - form.value.images.length
  if (remaining <= 0) return
  uni.chooseImage({
    count: remaining,
    sizeType: ['compressed'],
    sourceType: ['album', 'camera'],
    success: (res) => {
      form.value.images.push(...res.tempFilePaths)
    },
  })
}

function removeImage(idx: number) {
  form.value.images.splice(idx, 1)
}

async function onSubmit() {
  if (submitting.value || !canSubmit.value) return

  submitting.value = true
  try {
    await commonApi.submitFeedback({
      type: form.value.type,
      content: form.value.content,
      images: form.value.images,
      contact: form.value.contact,
    })
    uni.showToast({ title: t('common.done'), icon: 'success' })
    setTimeout(() => {
      uni.navigateBack()
    }, 1500)
  } catch (e: any) {
    uni.showToast({ title: e?.message || t('common.retry'), icon: 'none' })
  } finally {
    submitting.value = false
  }
}

onMounted(() => {
  uni.setNavigationBarTitle({ title: t('user.feedback') })
})
</script>

<style scoped>
.feedback-page {
  min-height: 100vh;
  background-color: #f5f5f5;
  padding-bottom: 140rpx;
}

.form-section {
  background-color: #fff;
  padding: 24rpx 32rpx;
  margin-top: 20rpx;
}

.section-label {
  font-size: 28rpx;
  color: #333;
  font-weight: 500;
  margin-bottom: 16rpx;
  display: block;
}

.type-selector {
  display: flex;
  gap: 16rpx;
}

.type-option {
  padding: 12rpx 32rpx;
  border: 2rpx solid #ddd;
  border-radius: 32rpx;
}

.type-option.active {
  border-color: #0f3a57;
  background-color: rgba(15, 58, 87, 0.08);
}

.type-text {
  font-size: 26rpx;
  color: #666;
}

.type-option.active .type-text {
  color: #0f3a57;
  font-weight: 500;
}

.content-textarea {
  width: 100%;
  height: 240rpx;
  font-size: 28rpx;
  color: #333;
  line-height: 1.6;
  box-sizing: border-box;
}

.char-count {
  display: block;
  text-align: right;
  font-size: 24rpx;
  color: #999;
  margin-top: 8rpx;
}

.photo-list {
  display: flex;
  flex-wrap: wrap;
  gap: 16rpx;
}

.photo-item {
  position: relative;
  width: 160rpx;
  height: 160rpx;
}

.photo-img {
  width: 100%;
  height: 100%;
  border-radius: 8rpx;
}

.photo-delete {
  position: absolute;
  top: -12rpx;
  right: -12rpx;
  width: 40rpx;
  height: 40rpx;
  background-color: rgba(0, 0, 0, 0.6);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
}

.delete-icon {
  font-size: 24rpx;
  color: #fff;
}

.photo-add {
  width: 160rpx;
  height: 160rpx;
  border: 2rpx dashed #ccc;
  border-radius: 8rpx;
  display: flex;
  align-items: center;
  justify-content: center;
}

.add-icon {
  font-size: 56rpx;
  color: #ccc;
}

.contact-input {
  height: 72rpx;
  font-size: 28rpx;
  color: #333;
  background-color: #f5f5f5;
  border-radius: 8rpx;
  padding: 0 24rpx;
}

.submit-section {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  padding: 20rpx 32rpx;
  background-color: #fff;
  box-shadow: 0 -2rpx 10rpx rgba(0, 0, 0, 0.05);
}

.submit-btn {
  background-color: #0f3a57;
  border-radius: 44rpx;
  height: 88rpx;
  display: flex;
  align-items: center;
  justify-content: center;
}

.submit-btn.disabled {
  opacity: 0.5;
}

.submit-text {
  font-size: 32rpx;
  color: #fff;
  font-weight: 500;
}
</style>
