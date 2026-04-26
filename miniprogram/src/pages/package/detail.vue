<template>
  <view class="page">
    <!-- Status banner -->
    <view class="status-banner">
      <view class="status-icon-wrap">
        <text class="status-icon-text">{{ statusEmoji }}</text>
      </view>
      <text class="status-label">{{ getStatusLabel(detail.status) }}</text>
    </view>

    <!-- Package info -->
    <view class="section">
      <view class="section-title">
        <text>包裹資訊</text>
      </view>
      <view class="info-row">
        <text class="info-label">物流單號</text>
        <view class="info-value-wrap">
          <text class="info-value">{{ detail.tracking_no }}</text>
          <text class="copy-btn" @tap="copyText(detail.tracking_no)">{{ t('common.copy') }}</text>
        </view>
      </view>
      <view class="info-row">
        <text class="info-label">重量</text>
        <text class="info-value">{{ detail.weight }}kg</text>
      </view>
      <view class="info-row">
        <text class="info-label">尺寸</text>
        <text class="info-value">{{ detail.length }}×{{ detail.width }}×{{ detail.height }} cm</text>
      </view>
      <view class="info-row">
        <text class="info-label">入庫時間</text>
        <text class="info-value">{{ detail.stored_at || '-' }}</text>
      </view>
    </view>

    <!-- Photos -->
    <view v-if="detail.photos && detail.photos.length > 0" class="section">
      <view class="section-title">
        <text>包裹照片</text>
      </view>
      <scroll-view class="photos-scroll" scroll-x enable-flex>
        <image
          v-for="(photo, idx) in detail.photos"
          :key="idx"
          class="photo-item"
          :src="photo"
          mode="aspectFill"
          @tap="previewImage(idx)"
        />
      </scroll-view>
    </view>

    <!-- Warehouse info -->
    <view class="section">
      <view class="section-title">
        <text>倉庫資訊</text>
      </view>
      <view class="info-row">
        <text class="info-label">倉庫名稱</text>
        <text class="info-value">{{ detail.warehouse_name }}</text>
      </view>
      <view class="info-row">
        <text class="info-label">倉庫地址</text>
        <text class="info-value">{{ detail.warehouse_address }}</text>
      </view>
      <view class="info-row">
        <text class="info-label">貨架號</text>
        <text class="info-value">{{ detail.shelf_no || '-' }}</text>
      </view>
    </view>

    <!-- Timeline -->
    <view class="section">
      <view class="section-title">
        <text>狀態歷史</text>
      </view>
      <view class="timeline">
        <view
          v-for="(event, idx) in detail.timeline"
          :key="idx"
          class="timeline-item"
          :class="{ first: idx === 0 }"
        >
          <view class="timeline-dot" :class="{ active: idx === 0 }"></view>
          <view v-if="idx < detail.timeline.length - 1" class="timeline-line"></view>
          <view class="timeline-content">
            <text class="timeline-desc">{{ event.description }}</text>
            <text class="timeline-time">{{ event.time }}</text>
          </view>
        </view>
      </view>
    </view>

    <!-- Action buttons -->
    <view class="action-bar">
      <view
        v-if="detail.status === 'stored'"
        class="action-btn primary"
        @tap="applyPack"
      >
        <text class="action-btn-text">申請打包</text>
      </view>
      <view
        v-if="detail.status === 'stored'"
        class="action-btn outline"
        @tap="goMerge"
      >
        <text class="action-btn-text outline-text">合箱</text>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { onLoad } from '@dcloudio/uni-app'
import { t } from '@/locale'
import { packageApi } from '@/api/package'

const packageId = ref('')
const detail = ref<any>({
  id: '',
  tracking_no: '',
  status: 'pending',
  weight: '0',
  length: '0',
  width: '0',
  height: '0',
  stored_at: '',
  warehouse_name: '',
  warehouse_address: '',
  shelf_no: '',
  photos: [],
  timeline: [],
})

const statusMap: Record<string, string> = {
  pending: t('package.status.pending'),
  stored: t('package.status.stored'),
  packed: t('package.status.packed'),
  shipped: t('package.status.shipped'),
  delivered: t('package.status.delivered'),
}

const statusEmojiMap: Record<string, string> = {
  pending: '📦',
  stored: '🏭',
  packed: '📋',
  shipped: '🚚',
  delivered: '✅',
}

const statusEmoji = computed(() => statusEmojiMap[detail.value.status] || '📦')

function getStatusLabel(status: string) {
  return statusMap[status] || status
}

async function fetchDetail() {
  try {
    const res = await packageApi.getDetail(packageId.value)
    if (res?.data) {
      detail.value = res.data
    }
  } catch (e) {
    // Mock data
    detail.value = {
      id: packageId.value,
      tracking_no: 'SF1234567890',
      status: 'stored',
      weight: '2.5',
      length: '30',
      width: '20',
      height: '15',
      stored_at: '2026-03-08 14:30',
      warehouse_name: '廣州倉',
      warehouse_address: '廣東省廣州市白雲區某某路88號',
      shelf_no: 'A-12-03',
      photos: [
        '/static/mock/pkg1.jpg',
        '/static/mock/pkg2.jpg',
      ],
      timeline: [
        { description: '包裹已入庫，貨架號 A-12-03', time: '2026-03-08 14:30' },
        { description: '包裹已到達廣州倉', time: '2026-03-08 10:20' },
        { description: '快遞已簽收', time: '2026-03-07 18:00' },
        { description: '快遞運輸中', time: '2026-03-06 09:15' },
        { description: '預報已建立', time: '2026-03-05 11:00' },
      ],
    }
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

function previewImage(index: number) {
  uni.previewImage({
    current: index,
    urls: detail.value.photos,
  })
}

function applyPack() {
  uni.showModal({
    title: '申請打包',
    content: '確認將此包裹申請打包？',
    success: (res) => {
      if (res.confirm) {
        uni.showToast({ title: '已提交申請', icon: 'success' })
      }
    },
  })
}

function goMerge() {
  uni.navigateTo({ url: `/pages/package/merge?id=${packageId.value}` })
}

onLoad((query) => {
  packageId.value = query?.id || ''
  if (packageId.value) {
    fetchDetail()
  }
})
</script>

<style scoped>
.page {
  min-height: 100vh;
  background-color: #f5f6f8;
  padding-bottom: 140rpx;
}

.status-banner {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #0f3a57, #1a5276);
  padding: 60rpx 0 48rpx;
}

.status-icon-wrap {
  width: 96rpx;
  height: 96rpx;
  border-radius: 48rpx;
  background-color: rgba(255, 255, 255, 0.15);
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 16rpx;
}

.status-icon-text {
  font-size: 48rpx;
}

.status-label {
  font-size: 32rpx;
  color: #fff;
  font-weight: 600;
}

.section {
  background-color: #fff;
  margin: 20rpx 24rpx;
  border-radius: 16rpx;
  padding: 28rpx;
  box-shadow: 0 2rpx 12rpx rgba(0, 0, 0, 0.04);
}

.section-title {
  margin-bottom: 24rpx;
  padding-bottom: 16rpx;
  border-bottom: 1rpx solid #f0f0f0;
}

.section-title text {
  font-size: 30rpx;
  font-weight: 600;
  color: #333;
}

.info-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 14rpx 0;
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
  text-align: right;
}

.copy-btn {
  font-size: 22rpx;
  color: #0f3a57;
  padding: 4rpx 12rpx;
  border: 1rpx solid #0f3a57;
  border-radius: 16rpx;
}

.photos-scroll {
  display: flex;
  white-space: nowrap;
  gap: 16rpx;
}

.photo-item {
  width: 180rpx;
  height: 180rpx;
  border-radius: 12rpx;
  flex-shrink: 0;
  margin-right: 16rpx;
}

/* Timeline */
.timeline {
  padding-left: 8rpx;
}

.timeline-item {
  display: flex;
  position: relative;
  padding-bottom: 32rpx;
  padding-left: 36rpx;
}

.timeline-item:last-child {
  padding-bottom: 0;
}

.timeline-dot {
  position: absolute;
  left: 0;
  top: 8rpx;
  width: 16rpx;
  height: 16rpx;
  border-radius: 50%;
  background-color: #ddd;
}

.timeline-dot.active {
  background-color: #0f3a57;
  box-shadow: 0 0 0 6rpx rgba(15, 58, 87, 0.15);
}

.timeline-line {
  position: absolute;
  left: 7rpx;
  top: 28rpx;
  width: 2rpx;
  bottom: 0;
  background-color: #eee;
}

.timeline-content {
  display: flex;
  flex-direction: column;
  gap: 6rpx;
}

.timeline-desc {
  font-size: 26rpx;
  color: #333;
}

.timeline-item.first .timeline-desc {
  color: #0f3a57;
  font-weight: 500;
}

.timeline-time {
  font-size: 22rpx;
  color: #999;
}

/* Action bar */
.action-bar {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  display: flex;
  align-items: center;
  gap: 20rpx;
  padding: 20rpx 32rpx;
  padding-bottom: calc(20rpx + env(safe-area-inset-bottom));
  background-color: #fff;
  box-shadow: 0 -2rpx 12rpx rgba(0, 0, 0, 0.06);
}

.action-btn {
  flex: 1;
  height: 84rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 42rpx;
}

.action-btn.primary {
  background-color: #0f3a57;
}

.action-btn.primary .action-btn-text {
  color: #fff;
  font-size: 30rpx;
  font-weight: 500;
}

.action-btn.outline {
  border: 2rpx solid #0f3a57;
  background-color: #fff;
}

.action-btn.outline .outline-text {
  color: #0f3a57;
  font-size: 30rpx;
  font-weight: 500;
}
</style>
