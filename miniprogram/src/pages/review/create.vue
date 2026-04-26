<template>
  <view class="review-create-page">
    <!-- Star rating -->
    <view class="rating-section">
      <text class="section-label">Rating</text>
      <view class="star-selector">
        <view
          v-for="i in 5"
          :key="i"
          class="star-btn"
          @click="form.rating = i"
        >
          <text class="star-icon" :class="{ filled: i <= form.rating }">
            {{ i <= form.rating ? '\u2605' : '\u2606' }}
          </text>
        </view>
      </view>
    </view>

    <!-- Text input -->
    <view class="content-section">
      <textarea
        class="review-textarea"
        v-model="form.content"
        placeholder="Share your experience..."
        :maxlength="500"
      />
      <text class="char-count">{{ form.content.length }}/500</text>
    </view>

    <!-- Photo upload -->
    <view class="photo-section">
      <text class="section-label">Photos (max 5)</text>
      <view class="photo-list">
        <view
          v-for="(photo, idx) in form.photos"
          :key="idx"
          class="photo-item"
        >
          <image class="photo-img" :src="photo" mode="aspectFill" />
          <view class="photo-delete" @click="removePhoto(idx)">
            <text class="delete-icon">x</text>
          </view>
        </view>
        <view
          v-if="form.photos.length < 5"
          class="photo-add"
          @click="choosePhoto"
        >
          <text class="add-icon">+</text>
        </view>
      </view>
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
import { ref, computed, onMounted, getCurrentInstance } from 'vue'
import { t } from '@/locale'
import { commonApi } from '@/api/common'

interface ReviewForm {
  rating: number
  content: string
  photos: string[]
}

const form = ref<ReviewForm>({
  rating: 5,
  content: '',
  photos: [],
})
const submitting = ref(false)
let orderId = ''

const canSubmit = computed(() => form.value.rating > 0)

function choosePhoto() {
  const remaining = 5 - form.value.photos.length
  if (remaining <= 0) return
  uni.chooseImage({
    count: remaining,
    sizeType: ['compressed'],
    sourceType: ['album', 'camera'],
    success: (res) => {
      form.value.photos.push(...res.tempFilePaths)
    },
  })
}

function removePhoto(idx: number) {
  form.value.photos.splice(idx, 1)
}

async function onSubmit() {
  if (submitting.value || !canSubmit.value) return
  submitting.value = true
  try {
    await commonApi.createReview({
      order_id: orderId,
      rating: form.value.rating,
      content: form.value.content,
      photos: form.value.photos,
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
  uni.setNavigationBarTitle({ title: t('review.create') })
  const instance = getCurrentInstance()
  const options = (instance?.proxy as any)?.$page?.options || {}
  orderId = options.order_id || ''
})
</script>

<style scoped>
.review-create-page {
  min-height: 100vh;
  background-color: #f5f5f5;
  padding-bottom: 120rpx;
}

.section-label {
  font-size: 28rpx;
  color: #333;
  font-weight: 500;
  margin-bottom: 16rpx;
  display: block;
}

.rating-section {
  background-color: #fff;
  padding: 32rpx;
}

.star-selector {
  display: flex;
  gap: 16rpx;
}

.star-btn {
  padding: 8rpx;
}

.star-icon {
  font-size: 56rpx;
  color: #ddd;
}

.star-icon.filled {
  color: #f5a623;
}

.content-section {
  background-color: #fff;
  margin-top: 20rpx;
  padding: 32rpx;
}

.review-textarea {
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

.photo-section {
  background-color: #fff;
  margin-top: 20rpx;
  padding: 32rpx;
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
