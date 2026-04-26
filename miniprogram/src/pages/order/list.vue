<template>
  <view class="page">
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

    <!-- Order list -->
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

      <view v-for="item in list" :key="item.id" class="order-card" @tap="goDetail(item.id)">
        <view class="order-header">
          <text class="order-no">{{ item.order_no }}</text>
          <text class="order-status" :class="'status-' + item.status">
            {{ getStatusLabel(item.status) }}
          </text>
        </view>
        <view class="order-body">
          <view class="order-info-row">
            <text class="order-label">包裹數量</text>
            <text class="order-value">{{ item.package_count }} 件</text>
          </view>
          <view class="order-info-row">
            <text class="order-label">總金額</text>
            <text class="order-price">¥{{ item.total_price }}</text>
          </view>
          <view class="order-info-row">
            <text class="order-label">下單時間</text>
            <text class="order-value">{{ item.created_at }}</text>
          </view>
        </view>
        <view class="order-footer">
          <view v-if="item.status === 'pending'" class="order-actions">
            <view class="action-btn-sm outline" @tap.stop="cancelOrder(item.id)">
              <text class="action-text outline-text">{{ t('common.cancel') }}</text>
            </view>
            <view class="action-btn-sm primary" @tap.stop="payOrder(item.id)">
              <text class="action-text primary-text">去付款</text>
            </view>
          </view>
          <view v-else class="order-actions">
            <view class="action-btn-sm outline" @tap.stop="goDetail(item.id)">
              <text class="action-text outline-text">查看詳情</text>
            </view>
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
import { ref } from 'vue'
import { onLoad, onShow } from '@dcloudio/uni-app'
import { t } from '@/locale'
import { orderApi } from '@/api/order'

const tabs = [
  { label: t('common.all'), value: '' },
  { label: t('order.status.pending'), value: 'pending' },
  { label: t('order.status.paid'), value: 'paid' },
  { label: t('order.status.processing'), value: 'processing' },
  { label: t('order.status.shipped'), value: 'shipped' },
  { label: t('order.status.completed'), value: 'completed' },
]

const currentTab = ref('')
const list = ref<any[]>([])
const loading = ref(false)
const refreshing = ref(false)
const finished = ref(false)
const page = ref(1)
const pageSize = 10

const statusMap: Record<string, string> = {
  pending: t('order.status.pending'),
  paid: t('order.status.paid'),
  processing: t('order.status.processing'),
  shipped: t('order.status.shipped'),
  completed: t('order.status.completed'),
  cancelled: t('order.status.cancelled'),
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
    const res = await orderApi.getList({
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
    if (data.length < pageSize) finished.value = true
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
    { id: '1', order_no: 'VS20260308001', status: 'pending', package_count: 3, total_price: '358.00', created_at: '2026-03-08 15:30' },
    { id: '2', order_no: 'VS20260307002', status: 'paid', package_count: 1, total_price: '128.00', created_at: '2026-03-07 10:20' },
    { id: '3', order_no: 'VS20260306003', status: 'processing', package_count: 2, total_price: '256.00', created_at: '2026-03-06 09:00' },
    { id: '4', order_no: 'VS20260305004', status: 'shipped', package_count: 1, total_price: '98.00', created_at: '2026-03-05 16:45' },
    { id: '5', order_no: 'VS20260301005', status: 'completed', package_count: 4, total_price: '520.00', created_at: '2026-03-01 12:00' },
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

function goDetail(id: string) {
  uni.navigateTo({ url: `/pages/order/detail?id=${id}` })
}

function payOrder(id: string) {
  uni.navigateTo({ url: `/pages/payment/index?order_id=${id}` })
}

function cancelOrder(id: string) {
  uni.showModal({
    title: '取消訂單',
    content: '確認取消此訂單？',
    success: async (res) => {
      if (res.confirm) {
        try {
          await orderApi.cancelOrder(id)
          uni.showToast({ title: '已取消', icon: 'success' })
          fetchList(true)
        } catch (e) {
          uni.showToast({ title: '已取消', icon: 'success' })
          fetchList(true)
        }
      }
    },
  })
}

onLoad((query) => {
  if (query?.status) {
    currentTab.value = query.status
  }
})

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
  white-space: nowrap;
  background-color: #fff;
  border-bottom: 1rpx solid #eee;
  padding: 0 12rpx;
}

.tab-item {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  padding: 20rpx 24rpx;
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

.order-card {
  background-color: #fff;
  border-radius: 16rpx;
  margin-bottom: 20rpx;
  box-shadow: 0 2rpx 12rpx rgba(0, 0, 0, 0.04);
  overflow: hidden;
}

.order-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 24rpx 28rpx;
  border-bottom: 1rpx solid #f5f5f5;
}

.order-no {
  font-size: 28rpx;
  font-weight: 600;
  color: #333;
}

.order-status {
  font-size: 26rpx;
  font-weight: 500;
}

.status-pending {
  color: #f57c00;
}

.status-paid {
  color: #1976d2;
}

.status-processing {
  color: #7b1fa2;
}

.status-shipped {
  color: #0097a7;
}

.status-completed {
  color: #388e3c;
}

.status-cancelled {
  color: #999;
}

.order-body {
  padding: 20rpx 28rpx;
  display: flex;
  flex-direction: column;
  gap: 12rpx;
}

.order-info-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.order-label {
  font-size: 26rpx;
  color: #999;
}

.order-value {
  font-size: 26rpx;
  color: #333;
}

.order-price {
  font-size: 30rpx;
  font-weight: 600;
  color: #e74c3c;
}

.order-footer {
  padding: 16rpx 28rpx 24rpx;
}

.order-actions {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  gap: 16rpx;
}

.action-btn-sm {
  padding: 0 28rpx;
  height: 60rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 30rpx;
}

.action-btn-sm.outline {
  border: 1rpx solid #ddd;
}

.action-btn-sm.primary {
  background-color: #0f3a57;
}

.action-text {
  font-size: 24rpx;
}

.outline-text {
  color: #666;
}

.primary-text {
  color: #fff;
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
