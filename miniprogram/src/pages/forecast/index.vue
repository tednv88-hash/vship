<template>
  <view class="page">
    <!-- Tabs -->
    <view class="tabs">
      <view
        v-for="tab in tabs"
        :key="tab.value"
        class="tab-item"
        :class="{ active: currentTab === tab.value }"
        @tap="switchTab(tab.value)"
      >
        <text>{{ tab.label }}</text>
      </view>
    </view>

    <!-- Forecast list -->
    <scroll-view
      class="list-wrap"
      scroll-y
      refresher-enabled
      :refresher-triggered="refreshing"
      @refresherrefresh="onRefresh"
      @scrolltolower="onLoadMore"
    >
      <view v-if="list.length === 0 && !loading" class="empty">
        <text class="empty-text">{{ t('common.noData') }}</text>
      </view>

      <view v-for="item in list" :key="item.id" class="forecast-card">
        <view class="forecast-header">
          <view class="courier-info">
            <text class="courier-name">{{ item.courier_name }}</text>
          </view>
          <view class="forecast-badge" :class="'status-' + item.status">
            <text class="badge-text">{{ getStatusLabel(item.status) }}</text>
          </view>
        </view>
        <view class="forecast-body">
          <view class="info-row">
            <text class="info-label">物流單號</text>
            <view class="info-value-wrap">
              <text class="info-value">{{ item.tracking_no }}</text>
              <text class="copy-btn" @tap.stop="copyText(item.tracking_no)">複製</text>
            </view>
          </view>
          <view v-if="item.description" class="info-row">
            <text class="info-label">商品描述</text>
            <text class="info-value">{{ item.description }}</text>
          </view>
          <view class="info-row">
            <text class="info-label">建立時間</text>
            <text class="info-value">{{ item.created_at }}</text>
          </view>
        </view>
      </view>

      <view v-if="loading" class="loading-more">
        <text>{{ t('common.loading') }}</text>
      </view>
      <view v-else-if="finished && list.length > 0" class="loading-more">
        <text class="no-more-text">— 沒有更多了 —</text>
      </view>
    </scroll-view>

    <!-- Floating add button -->
    <view class="fab" @tap="goCreate">
      <text class="fab-text">+</text>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { onShow } from '@dcloudio/uni-app'
import { t } from '@/locale'
import { packageApi } from '@/api/package'

const tabs = [
  { label: t('common.all'), value: '' },
  { label: t('package.status.pending'), value: 'pending' },
  { label: t('package.status.stored'), value: 'stored' },
]

const currentTab = ref('')
const list = ref<any[]>([])
const loading = ref(false)
const refreshing = ref(false)
const finished = ref(false)
const page = ref(1)
const pageSize = 10

const statusMap: Record<string, string> = {
  pending: t('package.status.pending'),
  stored: t('package.status.stored'),
}

function getStatusLabel(status: string) {
  return statusMap[status] || status
}

async function fetchList(reset = false) {
  if (loading.value) return
  if (!reset && finished.value) return

  if (reset) {
    page.value = 1
    finished.value = false
  }

  loading.value = true
  try {
    const res = await packageApi.getForecastList({
      page: page.value,
      page_size: pageSize,
      status: currentTab.value || undefined,
    })
    const data = res?.data?.list || res?.data || []
    if (reset) {
      list.value = data
    } else {
      list.value.push(...data)
    }
    if (data.length < pageSize) {
      finished.value = true
    }
    page.value++
  } catch (e) {
    if (reset) {
      list.value = getMockData()
    }
  } finally {
    loading.value = false
    refreshing.value = false
  }
}

function getMockData() {
  return [
    { id: '1', tracking_no: 'SF1234567890', courier_name: '順豐速運', status: 'stored', description: '衣服 x2', created_at: '2026-03-08' },
    { id: '2', tracking_no: 'YT9876543210', courier_name: '圓通快遞', status: 'pending', description: '電子配件', created_at: '2026-03-07' },
    { id: '3', tracking_no: 'ZT5555666677', courier_name: '中通快遞', status: 'pending', description: '書籍', created_at: '2026-03-06' },
    { id: '4', tracking_no: 'JD8877665544', courier_name: '京東物流', status: 'stored', description: '', created_at: '2026-03-04' },
  ]
}

function switchTab(val: string) {
  currentTab.value = val
  fetchList(true)
}

function onRefresh() {
  refreshing.value = true
  fetchList(true)
}

function onLoadMore() {
  fetchList(false)
}

function copyText(text: string) {
  uni.setClipboardData({
    data: text,
    success: () => {
      uni.showToast({ title: t('common.copied'), icon: 'success' })
    },
  })
}

function goCreate() {
  uni.navigateTo({ url: '/pages/forecast/create' })
}

onShow(() => {
  fetchList(true)
})
</script>

<style scoped>
.page {
  display: flex;
  flex-direction: column;
  height: 100vh;
  background-color: #f5f6f8;
}

.tabs {
  display: flex;
  background-color: #fff;
  border-bottom: 1rpx solid #eee;
}

.tab-item {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 24rpx 0;
  font-size: 28rpx;
  color: #666;
  position: relative;
}

.tab-item.active {
  color: #0f3a57;
  font-weight: 600;
}

.tab-item.active::after {
  content: '';
  position: absolute;
  bottom: 0;
  left: 50%;
  transform: translateX(-50%);
  width: 48rpx;
  height: 4rpx;
  background-color: #0f3a57;
  border-radius: 2rpx;
}

.list-wrap {
  flex: 1;
  padding: 20rpx 24rpx;
}

.forecast-card {
  background-color: #fff;
  border-radius: 16rpx;
  padding: 28rpx;
  margin-bottom: 20rpx;
  box-shadow: 0 2rpx 12rpx rgba(0, 0, 0, 0.04);
}

.forecast-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 20rpx;
  padding-bottom: 16rpx;
  border-bottom: 1rpx solid #f5f5f5;
}

.courier-name {
  font-size: 30rpx;
  font-weight: 600;
  color: #333;
}

.forecast-badge {
  padding: 4rpx 16rpx;
  border-radius: 20rpx;
}

.badge-text {
  font-size: 22rpx;
}

.status-pending {
  background-color: #fff3e0;
  color: #f57c00;
}

.status-stored {
  background-color: #e3f2fd;
  color: #1976d2;
}

.forecast-body {
  display: flex;
  flex-direction: column;
  gap: 14rpx;
}

.info-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.info-label {
  font-size: 26rpx;
  color: #999;
  flex-shrink: 0;
}

.info-value-wrap {
  display: flex;
  align-items: center;
  gap: 12rpx;
}

.info-value {
  font-size: 26rpx;
  color: #333;
}

.copy-btn {
  font-size: 22rpx;
  color: #0f3a57;
  padding: 4rpx 12rpx;
  border: 1rpx solid #0f3a57;
  border-radius: 16rpx;
}

.empty {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 120rpx 0;
}

.empty-text {
  font-size: 28rpx;
  color: #999;
}

.loading-more {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 30rpx 0;
}

.loading-more text {
  font-size: 24rpx;
  color: #999;
}

.no-more-text {
  color: #ccc;
}

.fab {
  position: fixed;
  right: 40rpx;
  bottom: 120rpx;
  width: 100rpx;
  height: 100rpx;
  border-radius: 50%;
  background-color: #0f3a57;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 8rpx 24rpx rgba(15, 58, 87, 0.35);
}

.fab-text {
  color: #fff;
  font-size: 48rpx;
  font-weight: 300;
  line-height: 1;
}
</style>
