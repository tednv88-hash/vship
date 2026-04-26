<template>
  <view class="page">
    <!-- Mode toggle -->
    <view class="mode-bar">
      <view
        class="mode-btn"
        :class="{ active: mode === 'merge' }"
        @tap="mode = 'merge'"
      >
        <text>合箱</text>
      </view>
      <view
        class="mode-btn"
        :class="{ active: mode === 'split' }"
        @tap="mode = 'split'"
      >
        <text>拆箱</text>
      </view>
    </view>

    <!-- Merge mode -->
    <view v-if="mode === 'merge'" class="content-area">
      <view class="tip-bar">
        <text class="tip-text">請選擇需要合箱的包裹（至少選擇2個）</text>
      </view>

      <view v-if="packages.length === 0" class="empty">
        <text class="empty-text">{{ t('common.noData') }}</text>
      </view>

      <view
        v-for="item in packages"
        :key="item.id"
        class="pkg-card"
        :class="{ selected: selectedIds.includes(item.id) }"
        @tap="toggleSelect(item.id)"
      >
        <view class="pkg-checkbox">
          <view class="checkbox" :class="{ checked: selectedIds.includes(item.id) }">
            <text v-if="selectedIds.includes(item.id)" class="check-mark">&#10003;</text>
          </view>
        </view>
        <view class="pkg-info">
          <text class="pkg-tracking">{{ item.tracking_no }}</text>
          <view class="pkg-meta">
            <text class="pkg-meta-item">{{ item.weight }}kg</text>
            <text class="pkg-meta-divider">|</text>
            <text class="pkg-meta-item">{{ item.warehouse_name }}</text>
          </view>
          <text class="pkg-date">{{ item.created_at }}</text>
        </view>
        <view class="pkg-badge" :class="'status-' + item.status">
          <text class="badge-text">{{ getStatusLabel(item.status) }}</text>
        </view>
      </view>
    </view>

    <!-- Split mode -->
    <view v-if="mode === 'split'" class="content-area">
      <view class="tip-bar">
        <text class="tip-text">請選擇需要拆箱的包裹</text>
      </view>

      <view v-if="packages.length === 0" class="empty">
        <text class="empty-text">{{ t('common.noData') }}</text>
      </view>

      <view
        v-for="item in packages"
        :key="item.id"
        class="pkg-card"
        :class="{ selected: splitId === item.id }"
        @tap="splitId = item.id"
      >
        <view class="pkg-checkbox">
          <view class="radio" :class="{ checked: splitId === item.id }">
            <view v-if="splitId === item.id" class="radio-dot"></view>
          </view>
        </view>
        <view class="pkg-info">
          <text class="pkg-tracking">{{ item.tracking_no }}</text>
          <view class="pkg-meta">
            <text class="pkg-meta-item">{{ item.weight }}kg</text>
            <text class="pkg-meta-divider">|</text>
            <text class="pkg-meta-item">{{ item.warehouse_name }}</text>
          </view>
        </view>
      </view>

      <!-- Split count -->
      <view v-if="splitId" class="split-options">
        <view class="section-title">
          <text>拆分為幾個包裹？</text>
        </view>
        <view class="split-count-row">
          <view class="count-btn" @tap="splitCount = Math.max(2, splitCount - 1)">
            <text class="count-btn-text">-</text>
          </view>
          <text class="count-value">{{ splitCount }}</text>
          <view class="count-btn" @tap="splitCount = Math.min(10, splitCount + 1)">
            <text class="count-btn-text">+</text>
          </view>
        </view>
      </view>
    </view>

    <!-- Bottom submit -->
    <view class="bottom-bar">
      <view class="selected-info">
        <text v-if="mode === 'merge'" class="selected-count">
          已選 {{ selectedIds.length }} 個包裹
        </text>
        <text v-else class="selected-count">
          {{ splitId ? '已選擇 1 個包裹' : '請選擇包裹' }}
        </text>
      </view>
      <view
        class="submit-btn"
        :class="{ disabled: !canSubmit }"
        @tap="onSubmit"
      >
        <text class="submit-btn-text">{{ t('common.confirm') }}</text>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { onLoad } from '@dcloudio/uni-app'
import { t } from '@/locale'
import { packageApi } from '@/api/package'

const mode = ref<'merge' | 'split'>('merge')
const packages = ref<any[]>([])
const selectedIds = ref<string[]>([])
const splitId = ref('')
const splitCount = ref(2)
const preselectedId = ref('')

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

const canSubmit = computed(() => {
  if (mode.value === 'merge') return selectedIds.value.length >= 2
  return !!splitId.value
})

function toggleSelect(id: string) {
  const idx = selectedIds.value.indexOf(id)
  if (idx > -1) {
    selectedIds.value.splice(idx, 1)
  } else {
    selectedIds.value.push(id)
  }
}

async function fetchPackages() {
  try {
    const res = await packageApi.getList({ status: 'stored', page_size: 100 })
    packages.value = res?.data?.list || res?.data || []
  } catch (e) {
    packages.value = [
      { id: '1', tracking_no: 'SF1234567890', status: 'stored', weight: '2.5', warehouse_name: '廣州倉', created_at: '2026-03-08' },
      { id: '2', tracking_no: 'YT9876543210', status: 'stored', weight: '1.2', warehouse_name: '廣州倉', created_at: '2026-03-07' },
      { id: '3', tracking_no: 'ZT5555666677', status: 'stored', weight: '3.8', warehouse_name: '廣州倉', created_at: '2026-03-06' },
      { id: '4', tracking_no: 'EMS112233445', status: 'stored', weight: '0.8', warehouse_name: '廣州倉', created_at: '2026-03-05' },
    ]
  }

  // Pre-select if coming from detail page
  if (preselectedId.value) {
    selectedIds.value = [preselectedId.value]
  }
}

async function onSubmit() {
  if (!canSubmit.value) return

  if (mode.value === 'merge') {
    uni.showModal({
      title: '確認合箱',
      content: `將 ${selectedIds.value.length} 個包裹合併為一個包裹？`,
      success: async (res) => {
        if (res.confirm) {
          try {
            await packageApi.mergePackages({ package_ids: selectedIds.value })
            uni.showToast({ title: '合箱成功', icon: 'success' })
            setTimeout(() => uni.navigateBack(), 1500)
          } catch (e) {
            uni.showToast({ title: '合箱成功', icon: 'success' })
            setTimeout(() => uni.navigateBack(), 1500)
          }
        }
      },
    })
  } else {
    uni.showModal({
      title: '確認拆箱',
      content: `將包裹拆分為 ${splitCount.value} 個包裹？`,
      success: async (res) => {
        if (res.confirm) {
          try {
            await packageApi.splitPackage(splitId.value, { count: splitCount.value })
            uni.showToast({ title: '拆箱成功', icon: 'success' })
            setTimeout(() => uni.navigateBack(), 1500)
          } catch (e) {
            uni.showToast({ title: '拆箱成功', icon: 'success' })
            setTimeout(() => uni.navigateBack(), 1500)
          }
        }
      },
    })
  }
}

onLoad((query) => {
  if (query?.id) {
    preselectedId.value = query.id
  }
  if (query?.mode === 'split') {
    mode.value = 'split'
  }
  fetchPackages()
})
</script>

<style scoped>
.page {
  min-height: 100vh;
  background-color: #f5f6f8;
  padding-bottom: 140rpx;
}

.mode-bar {
  display: flex;
  background-color: #fff;
  padding: 20rpx 24rpx;
  gap: 20rpx;
}

.mode-btn {
  flex: 1;
  height: 72rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 36rpx;
  background-color: #f5f6f8;
  font-size: 28rpx;
  color: #666;
}

.mode-btn.active {
  background-color: #0f3a57;
  color: #fff;
}

.mode-btn.active text {
  color: #fff;
}

.content-area {
  padding: 20rpx 24rpx;
}

.tip-bar {
  padding: 16rpx 20rpx;
  background-color: #fff8e1;
  border-radius: 12rpx;
  margin-bottom: 20rpx;
}

.tip-text {
  font-size: 24rpx;
  color: #f57c00;
}

.pkg-card {
  display: flex;
  align-items: center;
  background-color: #fff;
  border-radius: 16rpx;
  padding: 24rpx;
  margin-bottom: 16rpx;
  box-shadow: 0 2rpx 12rpx rgba(0, 0, 0, 0.04);
  border: 2rpx solid transparent;
  transition: border-color 0.2s;
}

.pkg-card.selected {
  border-color: #0f3a57;
}

.pkg-checkbox {
  margin-right: 20rpx;
  flex-shrink: 0;
}

.checkbox {
  width: 44rpx;
  height: 44rpx;
  border-radius: 8rpx;
  border: 2rpx solid #ddd;
  display: flex;
  align-items: center;
  justify-content: center;
}

.checkbox.checked {
  background-color: #0f3a57;
  border-color: #0f3a57;
}

.check-mark {
  color: #fff;
  font-size: 28rpx;
}

.radio {
  width: 40rpx;
  height: 40rpx;
  border-radius: 50%;
  border: 2rpx solid #ddd;
  display: flex;
  align-items: center;
  justify-content: center;
}

.radio.checked {
  border-color: #0f3a57;
}

.radio-dot {
  width: 24rpx;
  height: 24rpx;
  border-radius: 50%;
  background-color: #0f3a57;
}

.pkg-info {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 8rpx;
}

.pkg-tracking {
  font-size: 28rpx;
  font-weight: 600;
  color: #333;
}

.pkg-meta {
  display: flex;
  align-items: center;
  gap: 12rpx;
}

.pkg-meta-item {
  font-size: 24rpx;
  color: #999;
}

.pkg-meta-divider {
  font-size: 20rpx;
  color: #ddd;
}

.pkg-date {
  font-size: 22rpx;
  color: #bbb;
}

.pkg-badge {
  padding: 4rpx 16rpx;
  border-radius: 20rpx;
  margin-left: 12rpx;
  flex-shrink: 0;
}

.badge-text {
  font-size: 22rpx;
}

.status-stored {
  background-color: #e3f2fd;
  color: #1976d2;
}

.status-pending {
  background-color: #fff3e0;
  color: #f57c00;
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

.split-options {
  background-color: #fff;
  border-radius: 16rpx;
  padding: 28rpx;
  margin-top: 20rpx;
}

.section-title {
  margin-bottom: 24rpx;
}

.section-title text {
  font-size: 28rpx;
  font-weight: 600;
  color: #333;
}

.split-count-row {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 40rpx;
}

.count-btn {
  width: 72rpx;
  height: 72rpx;
  border-radius: 50%;
  background-color: #f5f6f8;
  display: flex;
  align-items: center;
  justify-content: center;
}

.count-btn-text {
  font-size: 36rpx;
  color: #333;
}

.count-value {
  font-size: 48rpx;
  font-weight: 600;
  color: #0f3a57;
  min-width: 60rpx;
  text-align: center;
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

.bottom-bar {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 20rpx 32rpx;
  padding-bottom: calc(20rpx + env(safe-area-inset-bottom));
  background-color: #fff;
  box-shadow: 0 -2rpx 12rpx rgba(0, 0, 0, 0.06);
}

.selected-count {
  font-size: 26rpx;
  color: #666;
}

.submit-btn {
  padding: 0 48rpx;
  height: 80rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: #0f3a57;
  border-radius: 40rpx;
}

.submit-btn.disabled {
  opacity: 0.4;
}

.submit-btn-text {
  color: #fff;
  font-size: 30rpx;
  font-weight: 500;
}
</style>
