<template>
  <view class="page">
    <!-- Search area -->
    <view class="search-area">
      <view class="search-title">
        <text class="title-text">{{ t('tracking.title') }}</text>
        <text class="title-sub">輸入物流單號查詢包裹追蹤狀態</text>
      </view>
      <view class="search-bar">
        <input
          v-model="trackingNo"
          class="search-input"
          :placeholder="t('tracking.number')"
          confirm-type="search"
          @confirm="onSearch"
        />
        <view class="search-btn" @tap="onSearch">
          <text class="search-btn-text">{{ t('tracking.search') }}</text>
        </view>
      </view>
    </view>

    <!-- Results -->
    <view v-if="searched" class="result-area">
      <!-- No result -->
      <view v-if="!result && !loading" class="empty">
        <text class="empty-text">未找到相關物流資訊</text>
        <text class="empty-sub">請確認單號是否正確</text>
      </view>

      <!-- Loading -->
      <view v-if="loading" class="loading-wrap">
        <text class="loading-text">{{ t('common.loading') }}</text>
      </view>

      <!-- Result card -->
      <view v-if="result && !loading" class="result-card">
        <view class="result-header">
          <view class="courier-row">
            <text class="courier-name">{{ result.courier_name }}</text>
            <view class="status-badge" :class="'status-' + result.status">
              <text class="badge-text">{{ result.status_label }}</text>
            </view>
          </view>
          <view class="tracking-row">
            <text class="tracking-label">單號：</text>
            <text class="tracking-value">{{ result.tracking_no }}</text>
            <text class="copy-btn" @tap="copyText(result.tracking_no)">複製</text>
          </view>
        </view>

        <!-- Timeline -->
        <view class="timeline-section">
          <view class="timeline">
            <view
              v-for="(event, idx) in result.events"
              :key="idx"
              class="timeline-item"
              :class="{ first: idx === 0 }"
            >
              <view class="timeline-left">
                <view class="timeline-dot" :class="{ active: idx === 0 }"></view>
                <view v-if="idx < result.events.length - 1" class="timeline-line"></view>
              </view>
              <view class="timeline-content">
                <text class="timeline-desc">{{ event.description }}</text>
                <text class="timeline-location" v-if="event.location">{{ event.location }}</text>
                <text class="timeline-time">{{ event.time }}</text>
              </view>
            </view>
          </view>
        </view>
      </view>
    </view>

    <!-- Recent searches -->
    <view v-if="!searched && recentSearches.length > 0" class="recent-section">
      <view class="recent-header">
        <text class="recent-title">最近查詢</text>
        <text class="recent-clear" @tap="clearRecent">清除</text>
      </view>
      <view
        v-for="item in recentSearches"
        :key="item"
        class="recent-item"
        @tap="trackingNo = item; onSearch()"
      >
        <text class="recent-text">{{ item }}</text>
        <text class="recent-arrow">&#8250;</text>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { onLoad } from '@dcloudio/uni-app'
import { t } from '@/locale'
import { commonApi } from '@/api/common'

const trackingNo = ref('')
const searched = ref(false)
const loading = ref(false)
const result = ref<any>(null)
const recentSearches = ref<string[]>([])

function loadRecentSearches() {
  try {
    const data = uni.getStorageSync('tracking_recent')
    if (data) {
      recentSearches.value = JSON.parse(data)
    }
  } catch (e) {
    // ignore
  }
}

function saveRecentSearch(no: string) {
  const list = recentSearches.value.filter((s) => s !== no)
  list.unshift(no)
  if (list.length > 10) list.length = 10
  recentSearches.value = list
  uni.setStorageSync('tracking_recent', JSON.stringify(list))
}

function clearRecent() {
  recentSearches.value = []
  uni.removeStorageSync('tracking_recent')
}

async function onSearch() {
  const no = trackingNo.value.trim()
  if (!no) {
    uni.showToast({ title: '請輸入物流單號', icon: 'none' })
    return
  }

  searched.value = true
  loading.value = true
  result.value = null

  try {
    const res = await commonApi.getTracking(no)
    if (res?.data) {
      result.value = res.data
      saveRecentSearch(no)
    }
  } catch (e) {
    // Mock data
    result.value = getMockResult(no)
    saveRecentSearch(no)
  } finally {
    loading.value = false
  }
}

function getMockResult(no: string) {
  return {
    tracking_no: no,
    courier_name: '順豐速運',
    status: 'in_transit',
    status_label: '運送中',
    events: [
      { description: '包裹正在派送中', location: '台北市信義區營業所', time: '2026-03-10 09:30' },
      { description: '包裹已到達目的地城市', location: '台北轉運中心', time: '2026-03-10 06:00' },
      { description: '包裹已通過海關', location: '桃園機場海關', time: '2026-03-09 22:30' },
      { description: '航班已起飛', location: '廣州白雲機場', time: '2026-03-09 14:00' },
      { description: '包裹已交付航空公司', location: '廣州集運倉', time: '2026-03-09 10:00' },
      { description: '包裹已打包完成', location: '廣州集運倉', time: '2026-03-08 16:00' },
      { description: '包裹已入庫', location: '廣州集運倉', time: '2026-03-07 11:00' },
    ],
  }
}

function copyText(text: string) {
  uni.setClipboardData({
    data: text,
    success: () => {
      uni.showToast({ title: t('common.copied'), icon: 'success' })
    },
  })
}

onLoad((query) => {
  if (query?.no) {
    trackingNo.value = query.no
    onSearch()
  }
  loadRecentSearches()
})
</script>

<style scoped>
.page {
  min-height: 100vh;
  background-color: #f5f6f8;
}

.search-area {
  background: linear-gradient(135deg, #0f3a57, #1a5276);
  padding: 48rpx 32rpx 40rpx;
}

.search-title {
  display: flex;
  flex-direction: column;
  margin-bottom: 28rpx;
}

.title-text {
  font-size: 40rpx;
  font-weight: 700;
  color: #fff;
  margin-bottom: 8rpx;
}

.title-sub {
  font-size: 24rpx;
  color: rgba(255, 255, 255, 0.7);
}

.search-bar {
  display: flex;
  align-items: center;
  gap: 16rpx;
}

.search-input {
  flex: 1;
  height: 80rpx;
  background-color: #fff;
  border-radius: 12rpx;
  padding: 0 24rpx;
  font-size: 28rpx;
  color: #333;
}

.search-btn {
  height: 80rpx;
  padding: 0 36rpx;
  background-color: rgba(255, 255, 255, 0.2);
  border: 2rpx solid rgba(255, 255, 255, 0.4);
  border-radius: 12rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.search-btn-text {
  color: #fff;
  font-size: 28rpx;
  font-weight: 500;
}

.result-area {
  padding: 20rpx 24rpx;
}

.result-card {
  background-color: #fff;
  border-radius: 16rpx;
  box-shadow: 0 2rpx 12rpx rgba(0, 0, 0, 0.04);
  overflow: hidden;
}

.result-header {
  padding: 28rpx;
  border-bottom: 1rpx solid #f0f0f0;
}

.courier-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 16rpx;
}

.courier-name {
  font-size: 32rpx;
  font-weight: 600;
  color: #333;
}

.status-badge {
  padding: 4rpx 16rpx;
  border-radius: 20rpx;
}

.badge-text {
  font-size: 22rpx;
}

.status-in_transit {
  background-color: #e3f2fd;
  color: #1976d2;
}

.status-delivered {
  background-color: #e8f5e9;
  color: #388e3c;
}

.status-pending {
  background-color: #fff3e0;
  color: #f57c00;
}

.tracking-row {
  display: flex;
  align-items: center;
  gap: 8rpx;
}

.tracking-label {
  font-size: 26rpx;
  color: #999;
}

.tracking-value {
  font-size: 26rpx;
  color: #333;
  font-weight: 500;
}

.copy-btn {
  font-size: 22rpx;
  color: #0f3a57;
  padding: 4rpx 12rpx;
  border: 1rpx solid #0f3a57;
  border-radius: 16rpx;
  margin-left: 8rpx;
}

/* Timeline */
.timeline-section {
  padding: 28rpx;
}

.timeline {
  padding-left: 0;
}

.timeline-item {
  display: flex;
  position: relative;
  padding-bottom: 36rpx;
}

.timeline-item:last-child {
  padding-bottom: 0;
}

.timeline-left {
  width: 40rpx;
  display: flex;
  flex-direction: column;
  align-items: center;
  flex-shrink: 0;
}

.timeline-dot {
  width: 16rpx;
  height: 16rpx;
  border-radius: 50%;
  background-color: #ddd;
  margin-top: 6rpx;
}

.timeline-dot.active {
  width: 20rpx;
  height: 20rpx;
  background-color: #0f3a57;
  box-shadow: 0 0 0 6rpx rgba(15, 58, 87, 0.15);
  margin-top: 4rpx;
}

.timeline-line {
  flex: 1;
  width: 2rpx;
  background-color: #eee;
  margin-top: 8rpx;
}

.timeline-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 4rpx;
  padding-left: 16rpx;
}

.timeline-desc {
  font-size: 26rpx;
  color: #333;
}

.timeline-item.first .timeline-desc {
  color: #0f3a57;
  font-weight: 600;
}

.timeline-location {
  font-size: 22rpx;
  color: #999;
}

.timeline-time {
  font-size: 22rpx;
  color: #bbb;
}

.empty {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 100rpx 0;
  gap: 12rpx;
}

.empty-text {
  font-size: 28rpx;
  color: #666;
}

.empty-sub {
  font-size: 24rpx;
  color: #999;
}

.loading-wrap {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 100rpx 0;
}

.loading-text {
  font-size: 28rpx;
  color: #999;
}

/* Recent searches */
.recent-section {
  padding: 20rpx 24rpx;
}

.recent-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 16rpx;
}

.recent-title {
  font-size: 28rpx;
  font-weight: 600;
  color: #333;
}

.recent-clear {
  font-size: 24rpx;
  color: #999;
}

.recent-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 24rpx 28rpx;
  background-color: #fff;
  border-radius: 12rpx;
  margin-bottom: 12rpx;
}

.recent-text {
  font-size: 28rpx;
  color: #333;
}

.recent-arrow {
  font-size: 28rpx;
  color: #ccc;
}
</style>
