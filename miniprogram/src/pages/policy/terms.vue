<template>
  <view class="policy-page">
    <!-- Loading -->
    <view v-if="loading" class="loading-wrap">
      <text class="loading-text">{{ t('common.loading') }}</text>
    </view>

    <!-- Error -->
    <view v-else-if="error" class="error-wrap">
      <text class="error-text">{{ error }}</text>
      <view class="retry-btn" @click="loadContent">
        <text class="retry-text">{{ t('common.retry') }}</text>
      </view>
    </view>

    <!-- Content -->
    <scroll-view v-else class="policy-scroll" scroll-y>
      <view class="policy-content">
        <rich-text :nodes="content" />
      </view>
    </scroll-view>
  </view>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { t } from '@/locale'
import { commonApi } from '@/api/common'

const loading = ref(true)
const error = ref('')
const content = ref('')

async function loadContent() {
  loading.value = true
  error.value = ''
  try {
    const res: any = await commonApi.getTerms()
    const data = res?.data || res
    content.value = data.content || data || ''
  } catch (e: any) {
    error.value = e?.message || t('common.retry')
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  uni.setNavigationBarTitle({ title: t('policy.terms') })
  loadContent()
})
</script>

<style scoped>
.policy-page {
  min-height: 100vh;
  background-color: #fff;
}

.policy-scroll {
  height: 100vh;
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

.policy-content {
  padding: 32rpx;
  font-size: 28rpx;
  color: #333;
  line-height: 1.8;
}
</style>
