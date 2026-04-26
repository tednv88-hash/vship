<template>
  <view class="prohibited-page">
    <!-- Loading -->
    <view v-if="loading" class="loading-wrap">
      <text class="loading-text">{{ t('common.loading') }}</text>
    </view>

    <!-- Error -->
    <view v-else-if="error" class="error-wrap">
      <text class="error-text">{{ error }}</text>
      <view class="retry-btn" @click="loadData">
        <text class="retry-text">{{ t('common.retry') }}</text>
      </view>
    </view>

    <!-- Content -->
    <view v-else class="prohibited-content">
      <!-- Warning banner -->
      <view class="warning-banner">
        <text class="warning-icon">\u26A0</text>
        <text class="warning-text">
          Please check the following prohibited and restricted items before shipping.
        </text>
      </view>

      <!-- Strictly prohibited -->
      <view class="category-section">
        <view class="category-header strictly">
          <text class="category-title">\u{1F6AB} Strictly Prohibited Items</text>
        </view>
        <view class="category-items">
          <view
            v-for="item in strictlyProhibited"
            :key="item.id"
            class="prohibited-item"
          >
            <view class="item-icon-wrap">
              <text class="item-icon">{{ item.icon || '\u2716' }}</text>
            </view>
            <view class="item-content">
              <text class="item-name">{{ item.name }}</text>
              <text class="item-desc">{{ item.description }}</text>
            </view>
          </view>
        </view>
      </view>

      <!-- Restricted -->
      <view class="category-section">
        <view class="category-header restricted">
          <text class="category-title">\u26A0 Restricted Items</text>
        </view>
        <view class="category-items">
          <view
            v-for="item in restrictedItems"
            :key="item.id"
            class="prohibited-item"
          >
            <view class="item-icon-wrap warning">
              <text class="item-icon">{{ item.icon || '\u26A0' }}</text>
            </view>
            <view class="item-content">
              <text class="item-name">{{ item.name }}</text>
              <text class="item-desc">{{ item.description }}</text>
            </view>
          </view>
        </view>
      </view>

      <!-- Empty -->
      <view
        v-if="strictlyProhibited.length === 0 && restrictedItems.length === 0"
        class="empty-wrap"
      >
        <text class="empty-text">{{ t('common.noData') }}</text>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { t } from '@/locale'
import { commonApi } from '@/api/common'

interface ProhibitedItem {
  id: string
  name: string
  description: string
  icon: string
  type: string // 'strictly' | 'restricted'
}

const loading = ref(true)
const error = ref('')
const strictlyProhibited = ref<ProhibitedItem[]>([])
const restrictedItems = ref<ProhibitedItem[]>([])

async function loadData() {
  loading.value = true
  error.value = ''
  try {
    const res: any = await commonApi.getProhibitedItems()
    const data = res?.data || res
    const items: ProhibitedItem[] = Array.isArray(data) ? data : data?.list || []
    strictlyProhibited.value = items.filter((i) => i.type === 'strictly')
    restrictedItems.value = items.filter((i) => i.type === 'restricted')
  } catch (e: any) {
    error.value = e?.message || t('common.retry')
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  uni.setNavigationBarTitle({ title: t('prohibited.title') })
  loadData()
})
</script>

<style scoped>
.prohibited-page {
  min-height: 100vh;
  background-color: #f5f5f5;
}

.loading-wrap,
.error-wrap,
.empty-wrap {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 120rpx 0;
}

.loading-text,
.empty-text {
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

.warning-banner {
  display: flex;
  align-items: center;
  background-color: #fff8e1;
  padding: 24rpx 32rpx;
  margin: 20rpx 24rpx;
  border-radius: 12rpx;
  border-left: 6rpx solid #f5a623;
}

.warning-icon {
  font-size: 36rpx;
  margin-right: 16rpx;
  flex-shrink: 0;
}

.warning-text {
  font-size: 26rpx;
  color: #666;
  line-height: 1.5;
}

.category-section {
  margin: 20rpx 24rpx;
}

.category-header {
  padding: 24rpx;
  border-radius: 12rpx 12rpx 0 0;
}

.category-header.strictly {
  background-color: #fde8e8;
}

.category-header.restricted {
  background-color: #fff3e0;
}

.category-title {
  font-size: 30rpx;
  font-weight: 600;
  color: #333;
}

.category-items {
  background-color: #fff;
  border-radius: 0 0 12rpx 12rpx;
}

.prohibited-item {
  display: flex;
  align-items: flex-start;
  padding: 24rpx;
  border-bottom: 1rpx solid #f0f0f0;
}

.prohibited-item:last-child {
  border-bottom: none;
}

.item-icon-wrap {
  width: 56rpx;
  height: 56rpx;
  border-radius: 50%;
  background-color: #fde8e8;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: 20rpx;
  flex-shrink: 0;
}

.item-icon-wrap.warning {
  background-color: #fff3e0;
}

.item-icon {
  font-size: 28rpx;
}

.item-content {
  flex: 1;
}

.item-name {
  font-size: 28rpx;
  color: #333;
  font-weight: 500;
  display: block;
  margin-bottom: 8rpx;
}

.item-desc {
  font-size: 24rpx;
  color: #999;
  line-height: 1.5;
}
</style>
