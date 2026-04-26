<template>
  <view class="address-list-page">
    <view v-if="list.length === 0 && !loading" class="empty">
      <text class="empty-text">{{ t('common.noData') }}</text>
    </view>

    <view v-for="item in list" :key="item.id" class="address-card">
      <view
        class="card-content"
        @touchstart="onTouchStart($event, item.id)"
        @touchmove="onTouchMove($event, item.id)"
        @touchend="onTouchEnd(item.id)"
        :style="{ transform: `translateX(${getOffset(item.id)}rpx)` }"
      >
        <view class="address-info" @tap="selectAddress(item)">
          <view class="name-row">
            <text class="name">{{ item.name }}</text>
            <text class="phone">{{ item.phone }}</text>
            <view v-if="item.is_default" class="default-badge">
              <text class="default-text">默認</text>
            </view>
          </view>
          <text class="full-address">
            {{ item.province }}{{ item.city }}{{ item.district }}{{ item.address }}
          </text>
        </view>
      </view>

      <!-- Swipe actions -->
      <view
        class="swipe-actions"
        :style="{ opacity: getOffset(item.id) < -20 ? 1 : 0 }"
      >
        <view class="action-btn edit-btn" @tap="handleEdit(item.id)">
          <text class="action-text">{{ t('common.edit') }}</text>
        </view>
        <view class="action-btn delete-btn" @tap="handleDelete(item.id)">
          <text class="action-text">{{ t('common.delete') }}</text>
        </view>
      </view>
    </view>

    <!-- Add button -->
    <view class="bottom-bar">
      <view class="add-btn" @tap="handleAdd">
        <text class="add-btn-text">{{ t('address.add') }}</text>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { t } from '@/locale'
import { commonApi } from '@/api/common'

interface AddressItem {
  id: string
  name: string
  phone: string
  province: string
  city: string
  district: string
  address: string
  is_default: boolean
}

const list = ref<AddressItem[]>([])
const loading = ref(false)

// Swipe state
const swipeState = reactive<Record<string, { startX: number; offset: number }>>({})

function getOffset(id: string): number {
  return swipeState[id]?.offset || 0
}

function onTouchStart(e: any, id: string) {
  const touch = e.touches[0]
  if (!swipeState[id]) {
    swipeState[id] = { startX: 0, offset: 0 }
  }
  swipeState[id].startX = touch.clientX
}

function onTouchMove(e: any, id: string) {
  const touch = e.touches[0]
  const diff = touch.clientX - swipeState[id].startX
  // Reset other items
  Object.keys(swipeState).forEach((key) => {
    if (key !== id) swipeState[key].offset = 0
  })
  if (diff < 0) {
    swipeState[id].offset = Math.max(diff * 2, -260)
  } else {
    swipeState[id].offset = 0
  }
}

function onTouchEnd(id: string) {
  if (swipeState[id]?.offset < -100) {
    swipeState[id].offset = -260
  } else if (swipeState[id]) {
    swipeState[id].offset = 0
  }
}

function selectAddress(item: AddressItem) {
  // If opened from order page, return selected address
  const pages = getCurrentPages()
  if (pages.length > 1) {
    uni.$emit('selectAddress', item)
    uni.navigateBack()
  }
}

async function loadList() {
  loading.value = true
  try {
    const res = (await commonApi.getAddresses()) as any
    list.value = res?.data || res || []
  } catch {
    uni.showToast({ title: '載入失敗', icon: 'none' })
  } finally {
    loading.value = false
  }
}

function handleAdd() {
  uni.navigateTo({ url: '/pages/address/edit' })
}

function handleEdit(id: string) {
  uni.navigateTo({ url: `/pages/address/edit?id=${id}` })
}

async function handleDelete(id: string) {
  const [, res] = (await uni.showModal({
    title: t('common.confirm'),
    content: '確定刪除此地址？',
  })) as any
  if (res?.confirm) {
    try {
      await commonApi.deleteAddress(id)
      list.value = list.value.filter((item) => item.id !== id)
      uni.showToast({ title: '刪除成功', icon: 'success' })
    } catch {
      uni.showToast({ title: '刪除失敗', icon: 'none' })
    }
  }
}

onMounted(() => {
  loadList()
})

// Refresh when returning from edit page
uni.$on('addressUpdated', () => {
  loadList()
})
</script>

<style scoped>
.address-list-page {
  min-height: 100vh;
  background: #f7f8fa;
  padding: 20rpx 24rpx 180rpx;
}

.empty {
  display: flex;
  justify-content: center;
  align-items: center;
  padding-top: 300rpx;
}

.empty-text {
  font-size: 28rpx;
  color: #999;
}

.address-card {
  position: relative;
  margin-bottom: 20rpx;
  border-radius: 16rpx;
  overflow: hidden;
}

.card-content {
  position: relative;
  z-index: 2;
  background: #fff;
  border-radius: 16rpx;
  padding: 30rpx;
  transition: transform 0.2s ease;
}

.address-info {
  width: 100%;
}

.name-row {
  display: flex;
  align-items: center;
  margin-bottom: 12rpx;
}

.name {
  font-size: 32rpx;
  font-weight: 600;
  color: #333;
  margin-right: 20rpx;
}

.phone {
  font-size: 28rpx;
  color: #666;
  margin-right: 16rpx;
}

.default-badge {
  background: #0f3a57;
  border-radius: 6rpx;
  padding: 4rpx 12rpx;
}

.default-text {
  font-size: 20rpx;
  color: #fff;
}

.full-address {
  font-size: 26rpx;
  color: #666;
  line-height: 1.5;
}

.swipe-actions {
  position: absolute;
  right: 0;
  top: 0;
  bottom: 0;
  display: flex;
  align-items: stretch;
  z-index: 1;
  transition: opacity 0.2s;
}

.action-btn {
  width: 130rpx;
  display: flex;
  align-items: center;
  justify-content: center;
}

.edit-btn {
  background: #0f3a57;
}

.delete-btn {
  background: #e74c3c;
  border-radius: 0 16rpx 16rpx 0;
}

.action-text {
  font-size: 26rpx;
  color: #fff;
}

.bottom-bar {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  padding: 20rpx 40rpx;
  padding-bottom: calc(20rpx + env(safe-area-inset-bottom));
  background: #fff;
  box-shadow: 0 -4rpx 12rpx rgba(0, 0, 0, 0.06);
}

.add-btn {
  width: 100%;
  height: 96rpx;
  background: #0f3a57;
  border-radius: 16rpx;
  display: flex;
  align-items: center;
  justify-content: center;
}

.add-btn-text {
  font-size: 32rpx;
  color: #fff;
  font-weight: 600;
}
</style>
