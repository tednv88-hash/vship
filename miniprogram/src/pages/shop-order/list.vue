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
        <!-- Shop header -->
        <view class="order-header">
          <text class="shop-name">国韵好运商城</text>
          <text class="order-status" :class="'status-' + item.status">
            {{ getStatusLabel(item.status) }}
          </text>
        </view>

        <!-- Products -->
        <view class="products">
          <view v-for="prod in item.products" :key="prod.id" class="product-item">
            <image class="product-img" :src="prod.image" mode="aspectFill" />
            <view class="product-info">
              <text class="product-name">{{ prod.name }}</text>
              <text class="product-spec" v-if="prod.spec">{{ prod.spec }}</text>
            </view>
            <view class="product-right">
              <text class="product-price">¥{{ prod.price }}</text>
              <text class="product-qty">x{{ prod.quantity }}</text>
            </view>
          </view>
        </view>

        <!-- Footer -->
        <view class="order-footer">
          <text class="order-total">
            共 {{ item.item_count }} 件，合計：
            <text class="total-price">¥{{ item.total_price }}</text>
          </text>
          <view class="order-actions">
            <view v-if="item.status === 'pending'" class="action-btn-sm primary" @tap.stop="payOrder(item.id)">
              <text class="action-text primary-text">去付款</text>
            </view>
            <view v-if="item.status === 'shipped'" class="action-btn-sm primary" @tap.stop="confirmReceive(item.id)">
              <text class="action-text primary-text">確認收貨</text>
            </view>
            <view v-if="item.status === 'pending'" class="action-btn-sm outline" @tap.stop="cancelOrder(item.id)">
              <text class="action-text outline-text">取消</text>
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
  { label: '待發貨', value: 'paid' },
  { label: '待收貨', value: 'shipped' },
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
  paid: '待發貨',
  shipped: '待收貨',
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
    const res = await orderApi.getShopOrders({
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
    {
      id: '1',
      order_no: 'SM20260308001',
      status: 'pending',
      item_count: 2,
      total_price: '199.00',
      products: [
        { id: 'p1', name: '日本限定零食禮盒', spec: '大份裝', price: '129.00', quantity: 1, image: '/static/mock/goods1.jpg' },
        { id: 'p2', name: '韓國面膜套裝', spec: '10片裝', price: '70.00', quantity: 1, image: '/static/mock/goods2.jpg' },
      ],
    },
    {
      id: '2',
      order_no: 'SM20260307002',
      status: 'paid',
      item_count: 1,
      total_price: '89.00',
      products: [
        { id: 'p3', name: '台灣茶葉禮盒', spec: '烏龍茶', price: '89.00', quantity: 1, image: '/static/mock/goods3.jpg' },
      ],
    },
    {
      id: '3',
      order_no: 'SM20260305003',
      status: 'shipped',
      item_count: 3,
      total_price: '356.00',
      products: [
        { id: 'p4', name: '護膚品套裝', spec: '保濕系列', price: '256.00', quantity: 1, image: '/static/mock/goods4.jpg' },
        { id: 'p5', name: '洗面乳', spec: '150ml', price: '50.00', quantity: 2, image: '/static/mock/goods5.jpg' },
      ],
    },
    {
      id: '4',
      order_no: 'SM20260301004',
      status: 'completed',
      item_count: 1,
      total_price: '450.00',
      products: [
        { id: 'p6', name: '藍牙耳機', spec: '黑色', price: '450.00', quantity: 1, image: '/static/mock/goods6.jpg' },
      ],
    },
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
  uni.navigateTo({ url: `/pages/shop-order/detail?id=${id}` })
}

function payOrder(id: string) {
  uni.navigateTo({ url: `/pages/payment/index?order_id=${id}&type=shop` })
}

function cancelOrder(id: string) {
  uni.showModal({
    title: '取消訂單',
    content: '確認取消此訂單？',
    success: (res) => {
      if (res.confirm) {
        uni.showToast({ title: '已取消', icon: 'success' })
        fetchList(true)
      }
    },
  })
}

function confirmReceive(id: string) {
  uni.showModal({
    title: '確認收貨',
    content: '確認已收到商品？',
    success: (res) => {
      if (res.confirm) {
        uni.showToast({ title: '已確認', icon: 'success' })
        fetchList(true)
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
  padding: 20rpx 28rpx;
  border-bottom: 1rpx solid #f5f5f5;
}

.shop-name {
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

.status-shipped {
  color: #0097a7;
}

.status-completed {
  color: #388e3c;
}

.status-cancelled {
  color: #999;
}

.products {
  padding: 16rpx 28rpx;
}

.product-item {
  display: flex;
  align-items: center;
  padding: 12rpx 0;
  gap: 20rpx;
}

.product-img {
  width: 140rpx;
  height: 140rpx;
  border-radius: 12rpx;
  background-color: #f5f5f5;
  flex-shrink: 0;
}

.product-info {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 8rpx;
}

.product-name {
  font-size: 28rpx;
  color: #333;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.product-spec {
  font-size: 24rpx;
  color: #999;
  background-color: #f5f6f8;
  padding: 2rpx 12rpx;
  border-radius: 4rpx;
  display: inline;
  align-self: flex-start;
}

.product-right {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  gap: 6rpx;
  flex-shrink: 0;
}

.product-price {
  font-size: 28rpx;
  color: #333;
  font-weight: 500;
}

.product-qty {
  font-size: 24rpx;
  color: #999;
}

.order-footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16rpx 28rpx 20rpx;
  border-top: 1rpx solid #f5f5f5;
}

.order-total {
  font-size: 24rpx;
  color: #666;
}

.total-price {
  font-size: 30rpx;
  font-weight: 600;
  color: #e74c3c;
}

.order-actions {
  display: flex;
  align-items: center;
  gap: 12rpx;
}

.action-btn-sm {
  padding: 0 24rpx;
  height: 56rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 28rpx;
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
