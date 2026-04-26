<template>
  <view class="dealer-orders">
    <!-- Status filter -->
    <view class="filter-bar">
      <view
        v-for="item in statusFilters"
        :key="item.key"
        class="filter-item"
        :class="{ active: currentFilter === item.key }"
        @click="onFilterChange(item.key)"
      >
        <text>{{ item.label }}</text>
      </view>
    </view>

    <!-- Orders list -->
    <view class="orders-list">
      <view v-for="order in orderList" :key="order.id" class="order-card">
        <view class="order-header">
          <text class="order-no">訂單號: {{ order.order_no }}</text>
          <text class="order-status" :class="order.status">{{ getStatusLabel(order.status) }}</text>
        </view>

        <view class="order-body">
          <view class="buyer-info">
            <image class="buyer-avatar" :src="order.buyer_avatar || '/static/default-avatar.png'" />
            <view class="buyer-detail">
              <text class="buyer-name">{{ order.buyer_name }}</text>
              <text class="order-date">{{ order.created_at }}</text>
            </view>
          </view>
          <view class="commission-info">
            <text class="commission-label">佣金</text>
            <text class="commission-value">¥{{ order.commission }}</text>
          </view>
        </view>

        <view class="order-footer">
          <text class="order-amount">訂單金額: ¥{{ order.amount }}</text>
        </view>
      </view>
    </view>

    <!-- Empty / Loading -->
    <view v-if="!loading && orderList.length === 0" class="empty">
      <text>{{ t('common.noData') }}</text>
    </view>
    <view class="load-more">
      <text v-if="loading">{{ t('common.loading') }}</text>
      <text v-else-if="finished && orderList.length > 0">— 已經到底了 —</text>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { t } from '@/locale'
import { commonApi } from '@/api/common'

const statusFilters = [
  { key: 'all', label: t('common.all') },
  { key: 'pending', label: '待結算' },
  { key: 'settled', label: '已結算' },
  { key: 'cancelled', label: '已取消' },
]

const currentFilter = ref('all')
const orderList = ref<any[]>([])
const loading = ref(false)
const finished = ref(false)
const page = ref(1)
const limit = 10

onMounted(() => {
  loadOrders()
})

function onFilterChange(key: string) {
  currentFilter.value = key
  resetList()
}

function resetList() {
  page.value = 1
  finished.value = false
  orderList.value = []
  loadOrders()
}

async function loadOrders() {
  if (loading.value || finished.value) return
  loading.value = true
  try {
    const params: any = { page: page.value, limit }
    if (currentFilter.value !== 'all') {
      params.status = currentFilter.value
    }
    const res = await commonApi.getDealerOrders(params)
    const list = res?.data?.list || res?.data || []
    if (list.length < limit) {
      finished.value = true
    }
    orderList.value.push(...list)
    page.value++
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

function getStatusLabel(status: string): string {
  const map: Record<string, string> = {
    pending: '待結算',
    settled: '已結算',
    cancelled: '已取消',
  }
  return map[status] || status
}

onReachBottom(() => {
  loadOrders()
})
</script>

<style scoped>
.dealer-orders {
  min-height: 100vh;
  background: #f5f5f5;
}

.filter-bar {
  display: flex;
  background: #fff;
  padding: 20rpx 0;
  margin-bottom: 16rpx;
  position: sticky;
  top: 0;
  z-index: 10;
}

.filter-item {
  flex: 1;
  text-align: center;
  font-size: 26rpx;
  color: #666;
  padding: 12rpx 0;
  position: relative;
}

.filter-item.active {
  color: #0f3a57;
  font-weight: 600;
}

.filter-item.active::after {
  content: '';
  position: absolute;
  bottom: 0;
  left: 50%;
  transform: translateX(-50%);
  width: 40rpx;
  height: 4rpx;
  background: #0f3a57;
  border-radius: 2rpx;
}

.orders-list {
  padding: 0 24rpx;
}

.order-card {
  background: #fff;
  border-radius: 16rpx;
  padding: 24rpx;
  margin-bottom: 16rpx;
}

.order-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-bottom: 20rpx;
  border-bottom: 1rpx solid #f5f5f5;
}

.order-no {
  font-size: 24rpx;
  color: #999;
}

.order-status {
  font-size: 24rpx;
  font-weight: 500;
}

.order-status.pending {
  color: #ff9800;
}

.order-status.settled {
  color: #4caf50;
}

.order-status.cancelled {
  color: #999;
}

.order-body {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20rpx 0;
}

.buyer-info {
  display: flex;
  align-items: center;
  gap: 16rpx;
}

.buyer-avatar {
  width: 64rpx;
  height: 64rpx;
  border-radius: 50%;
}

.buyer-detail {
  display: flex;
  flex-direction: column;
  gap: 6rpx;
}

.buyer-name {
  font-size: 28rpx;
  color: #333;
}

.order-date {
  font-size: 22rpx;
  color: #999;
}

.commission-info {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  gap: 4rpx;
}

.commission-label {
  font-size: 22rpx;
  color: #999;
}

.commission-value {
  font-size: 34rpx;
  color: #e64340;
  font-weight: 600;
}

.order-footer {
  border-top: 1rpx solid #f5f5f5;
  padding-top: 16rpx;
}

.order-amount {
  font-size: 24rpx;
  color: #666;
}

.empty {
  text-align: center;
  padding: 100rpx 0;
  color: #999;
  font-size: 28rpx;
}

.load-more {
  text-align: center;
  padding: 30rpx 0;
  font-size: 24rpx;
  color: #999;
}
</style>
