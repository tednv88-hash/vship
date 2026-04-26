<template>
  <view class="help-detail-page">
    <!-- Loading -->
    <view v-if="loading" class="loading-wrap">
      <text class="loading-text">{{ t('common.loading') }}</text>
    </view>

    <!-- Error -->
    <view v-else-if="error" class="error-wrap">
      <text class="error-text">{{ error }}</text>
      <view class="retry-btn" @click="loadDetail">
        <text class="retry-text">{{ t('common.retry') }}</text>
      </view>
    </view>

    <!-- Content -->
    <view v-else class="detail-content">
      <view class="detail-header">
        <text class="detail-title">{{ detail.title }}</text>
      </view>
      <view class="detail-body">
        <rich-text :nodes="detail.content" />
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, onMounted, getCurrentInstance } from 'vue'
import { t } from '@/locale'
import { commonApi } from '@/api/common'

const loading = ref(true)
const error = ref('')
const detail = ref<{ title: string; content: string }>({
  title: '',
  content: '',
})

let articleId = ''

async function loadDetail() {
  if (!articleId) return
  loading.value = true
  error.value = ''
  try {
    const res: any = await commonApi.getHelpDetail(articleId)
    const data = res?.data || res
    detail.value = {
      title: data.title || '',
      content: data.content || '',
    }
  } catch (e: any) {
    error.value = e?.message || t('common.retry')
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  uni.setNavigationBarTitle({ title: t('help.detail') })
  const instance = getCurrentInstance()
  const options = (instance?.proxy as any)?.$page?.options || {}
  articleId = options.id || ''
  loadDetail()
})
</script>

<style scoped>
.help-detail-page {
  min-height: 100vh;
  background-color: #fff;
}

.loading-wrap,
.error-wrap {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 120rpx 0;
}

.loading-text {
  font-size: 28rpx;
  color: #999;
}

.error-text {
  font-size: 28rpx;
  color: #e64340;
  margin-bottom: 24rpx;
}

.retry-btn {
  padding: 16rpx 48rpx;
  background-color: #0f3a57;
  border-radius: 8rpx;
}

.retry-text {
  font-size: 28rpx;
  color: #fff;
}

.detail-header {
  padding: 32rpx;
  border-bottom: 1rpx solid #eee;
}

.detail-title {
  font-size: 36rpx;
  font-weight: 600;
  color: #333;
  line-height: 1.5;
}

.detail-body {
  padding: 32rpx;
  font-size: 28rpx;
  color: #333;
  line-height: 1.8;
}
</style>
