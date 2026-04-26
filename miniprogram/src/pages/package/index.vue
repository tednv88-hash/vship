<template>
  <view class="page">
    <!-- Search bar -->
    <view class="search-bar">
      <view class="search-input-wrap">
        <text class="search-icon">&#xe612;</text>
        <input
          v-model="keyword"
          class="search-input"
          :placeholder="t('common.search')"
          confirm-type="search"
          @confirm="onSearch"
        />
        <text v-if="keyword" class="search-clear" @tap="keyword = ''; onSearch()">&#xe613;</text>
      </view>
    </view>

    <!-- Tabs -->
    <scroll-view class="tabs" scroll-x enable-flex>
      <view
        v-for="tab in tabs"
        :key="tab.value"
        class="tab-item"
        :class="{ active: currentTab === tab.value }"
        @tap="switchTab(tab.value)"
      >
        <text>{{ tab.label }}</text>
      </view>
    </scroll-view>

    <!-- Package list -->
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

      <view v-for="item in list" :key="item.id" class="pkg-card" @tap="goDetail(item.id)">
        <view class="pkg-header">
          <text class="pkg-tracking">{{ item.tracking_no }}</text>
          <view class="pkg-badge" :class="'status-' + item.status">
            <text class="badge-text">{{ getStatusLabel(item.status) }}</text>
          </view>
        </view>
        <view class="pkg-body">
          <view class="pkg-info-row">
            <text class="pkg-label">重量</text>
            <text class="pkg-value">{{ item.weight }}kg</text>
          </view>
          <view class="pkg-info-row">
            <text class="pkg-label">倉庫</text>
            <text class="pkg-value">{{ item.warehouse_name }}</text>
          </view>
          <view class="pkg-info-row">
            <text class="pkg-label">日期</text>
            <text class="pkg-value">{{ item.created_at }}</text>
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
  </view>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { onLoad, onShow } from '@dcloudio/uni-app'
import { t } from '@/locale'
import { packageApi } from '@/api/package'

const tabs = [
  { label: t('common.all'), value: '' },
  { label: t('package.status.pending'), value: 'pending' },
  { label: t('package.status.stored'), value: 'stored' },
  { label: t('package.status.packed'), value: 'packed' },
  { label: t('package.status.shipped'), value: 'shipped' },
  { label: t('package.status.delivered'), value: 'delivered' },
]

const currentTab = ref('')
const keyword = ref('')
const list = ref<any[]>([])
const loading = ref(false)
const refreshing = ref(false)
const finished = ref(false)
const page = ref(1)
const pageSize = 10

const statusMap: Record<string, string> = {
  pending: t('package.status.pending'),
  stored: t('package.status.stored'),
  packed: t('package.status.packed'),
  shipped: t('package.status.shipped'),
  delivered: t('package.status.delivered'),
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
    const res = await packageApi.getList({
      page: page.value,
      page_size: pageSize,
      status: currentTab.value || undefined,
      keyword: keyword.value || undefined,
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
    // Use mock data in development
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
    { id: '1', tracking_no: 'SF1234567890', status: 'stored', weight: '2.5', warehouse_name: '廣州倉', created_at: '2026-03-08' },
    { id: '2', tracking_no: 'YT9876543210', status: 'pending', weight: '1.2', warehouse_name: '深圳倉', created_at: '2026-03-07' },
    { id: '3', tracking_no: 'ZT5555666677', status: 'packed', weight: '3.8', warehouse_name: '廣州倉', created_at: '2026-03-06' },
    { id: '4', tracking_no: 'EMS112233445', status: 'shipped', weight: '0.8', warehouse_name: '上海倉', created_at: '2026-03-05' },
    { id: '5', tracking_no: 'JD8877665544', status: 'delivered', weight: '5.0', warehouse_name: '廣州倉', created_at: '2026-03-01' },
  ]
}

function switchTab(val: string) {
  currentTab.value = val
  fetchList(true)
}

function onSearch() {
  fetchList(true)
}

function onRefresh() {
  refreshing.value = true
  fetchList(true)
}

function onLoadMore() {
  fetchList(false)
}

function goDetail(id: string) {
  uni.navigateTo({ url: `/pages/package/detail?id=${id}` })
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

.search-bar {
  padding: 20rpx 24rpx;
  background-color: #0f3a57;
}

.search-input-wrap {
  display: flex;
  align-items: center;
  background-color: #fff;
  border-radius: 36rpx;
  padding: 0 24rpx;
  height: 72rpx;
}

.search-icon {
  font-size: 32rpx;
  color: #999;
  margin-right: 12rpx;
}

.search-input {
  flex: 1;
  font-size: 28rpx;
  height: 72rpx;
}

.search-clear {
  font-size: 28rpx;
  color: #ccc;
  padding: 10rpx;
}

.tabs {
  display: flex;
  white-space: nowrap;
  background-color: #fff;
  border-bottom: 1rpx solid #eee;
  padding: 0 12rpx;
}

.tab-item {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  padding: 20rpx 28rpx;
  font-size: 28rpx;
  color: #666;
  position: relative;
  flex-shrink: 0;
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

.pkg-card {
  background-color: #fff;
  border-radius: 16rpx;
  padding: 28rpx;
  margin-bottom: 20rpx;
  box-shadow: 0 2rpx 12rpx rgba(0, 0, 0, 0.04);
}

.pkg-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 20rpx;
}

.pkg-tracking {
  font-size: 30rpx;
  font-weight: 600;
  color: #333;
}

.pkg-badge {
  padding: 4rpx 16rpx;
  border-radius: 20rpx;
  font-size: 22rpx;
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

.status-packed {
  background-color: #e8f5e9;
  color: #388e3c;
}

.status-shipped {
  background-color: #ede7f6;
  color: #7b1fa2;
}

.status-delivered {
  background-color: #e0f2f1;
  color: #00796b;
}

.pkg-body {
  display: flex;
  flex-direction: column;
  gap: 12rpx;
}

.pkg-info-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.pkg-label {
  font-size: 26rpx;
  color: #999;
}

.pkg-value {
  font-size: 26rpx;
  color: #333;
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
</style>
