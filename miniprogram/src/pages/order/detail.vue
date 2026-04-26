<template>
  <view class="page">
    <!-- Status banner -->
    <view class="status-banner" :class="'banner-' + detail.status">
      <text class="status-text">{{ getStatusLabel(detail.status) }}</text>
      <text class="status-desc">{{ statusDesc }}</text>
    </view>

    <!-- Shipping info -->
    <view class="section">
      <view class="section-title"><text>寄送資訊</text></view>
      <view class="info-row">
        <text class="info-label">訂單編號</text>
        <view class="info-value-wrap">
          <text class="info-value">{{ detail.order_no }}</text>
          <text class="copy-btn" @tap="copyText(detail.order_no)">複製</text>
        </view>
      </view>
      <view class="info-row">
        <text class="info-label">寄送線路</text>
        <text class="info-value">{{ detail.route_name }}</text>
      </view>
      <view class="info-row">
        <text class="info-label">收件人</text>
        <text class="info-value">{{ detail.address?.name }} {{ detail.address?.phone }}</text>
      </view>
      <view class="info-row">
        <text class="info-label">收件地址</text>
        <text class="info-value address-text">{{ detail.address?.region }} {{ detail.address?.detail }}</text>
      </view>
      <view class="info-row">
        <text class="info-label">下單時間</text>
        <text class="info-value">{{ detail.created_at }}</text>
      </view>
    </view>

    <!-- Package list -->
    <view class="section">
      <view class="section-title"><text>包裹列表</text></view>
      <view v-for="pkg in detail.packages" :key="pkg.id" class="pkg-item">
        <view class="pkg-main">
          <text class="pkg-tracking">{{ pkg.tracking_no }}</text>
          <text class="pkg-weight">{{ pkg.weight }}kg</text>
        </view>
        <text class="pkg-desc">{{ pkg.description || '—' }}</text>
      </view>
    </view>

    <!-- Price breakdown -->
    <view class="section">
      <view class="section-title"><text>費用明細</text></view>
      <view class="price-row">
        <text class="price-label">運費</text>
        <text class="price-value">¥{{ detail.shipping_fee }}</text>
      </view>
      <view v-for="svc in detail.addon_services" :key="svc.name" class="price-row">
        <text class="price-label">{{ svc.name }}</text>
        <text class="price-value">¥{{ svc.price }}</text>
      </view>
      <view v-if="detail.insurance_fee" class="price-row">
        <text class="price-label">保險費</text>
        <text class="price-value">¥{{ detail.insurance_fee }}</text>
      </view>
      <view v-if="detail.coupon_discount" class="price-row">
        <text class="price-label">優惠券</text>
        <text class="price-value discount">-¥{{ detail.coupon_discount }}</text>
      </view>
      <view class="price-row total">
        <text class="price-label total-label">合計</text>
        <text class="price-value total-value">¥{{ detail.total_price }}</text>
      </view>
    </view>

    <!-- Timeline -->
    <view class="section">
      <view class="section-title"><text>訂單進度</text></view>
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
    <view v-if="detail.status === 'pending'" class="action-bar">
      <view class="action-btn outline" @tap="cancelOrder">
        <text class="action-btn-text outline-text">取消訂單</text>
      </view>
      <view class="action-btn primary" @tap="payOrder">
        <text class="action-btn-text primary-text">去付款</text>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { onLoad } from '@dcloudio/uni-app'
import { t } from '@/locale'
import { orderApi } from '@/api/order'

const orderId = ref('')
const detail = ref<any>({
  id: '',
  order_no: '',
  status: 'pending',
  route_name: '',
  address: {},
  packages: [],
  shipping_fee: '0',
  addon_services: [],
  insurance_fee: '',
  coupon_discount: '',
  total_price: '0',
  created_at: '',
  timeline: [],
})

const statusMap: Record<string, string> = {
  pending: t('order.status.pending'),
  paid: t('order.status.paid'),
  processing: t('order.status.processing'),
  shipped: t('order.status.shipped'),
  completed: t('order.status.completed'),
  cancelled: t('order.status.cancelled'),
}

const statusDescMap: Record<string, string> = {
  pending: '請在24小時內完成付款，逾期將自動取消',
  paid: '已收到您的付款，正在安排處理',
  processing: '包裹正在打包處理中',
  shipped: '包裹已發出，請注意查收',
  completed: '訂單已完成，感謝您的使用',
  cancelled: '訂單已取消',
}

function getStatusLabel(status: string) {
  return statusMap[status] || status
}

const statusDesc = computed(() => statusDescMap[detail.value.status] || '')

async function fetchDetail() {
  try {
    const res = await orderApi.getDetail(orderId.value)
    if (res?.data) {
      detail.value = res.data
    }
  } catch (e) {
    detail.value = {
      id: orderId.value,
      order_no: 'VS20260308001',
      status: 'pending',
      route_name: '廣州 → 台北（空運）',
      address: {
        name: '王小明',
        phone: '0912345678',
        region: '台北市信義區',
        detail: '忠孝東路五段100號',
      },
      packages: [
        { id: '1', tracking_no: 'SF1234567890', weight: '2.5', description: '衣服 x2' },
        { id: '2', tracking_no: 'YT9876543210', weight: '1.2', description: '電子配件' },
        { id: '3', tracking_no: 'ZT5555666677', weight: '3.8', description: '書籍' },
      ],
      shipping_fee: '298.00',
      addon_services: [
        { name: '加固包裝', price: '15.00' },
        { name: '拍照驗貨', price: '5.00' },
      ],
      insurance_fee: '12.00',
      coupon_discount: '',
      total_price: '330.00',
      created_at: '2026-03-08 15:30',
      timeline: [
        { description: '訂單已建立，等待付款', time: '2026-03-08 15:30' },
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

function cancelOrder() {
  uni.showModal({
    title: '取消訂單',
    content: '確認取消此訂單？取消後不可恢復。',
    success: async (res) => {
      if (res.confirm) {
        try {
          await orderApi.cancelOrder(orderId.value)
          uni.showToast({ title: '已取消', icon: 'success' })
          detail.value.status = 'cancelled'
        } catch (e) {
          uni.showToast({ title: '已取消', icon: 'success' })
          detail.value.status = 'cancelled'
        }
      }
    },
  })
}

function payOrder() {
  uni.navigateTo({ url: `/pages/payment/index?order_id=${orderId.value}` })
}

onLoad((query) => {
  orderId.value = query?.id || ''
  if (orderId.value) {
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
  padding: 48rpx 32rpx;
  display: flex;
  flex-direction: column;
  gap: 8rpx;
}

.banner-pending {
  background: linear-gradient(135deg, #f57c00, #ff9800);
}

.banner-paid {
  background: linear-gradient(135deg, #1565c0, #1976d2);
}

.banner-processing {
  background: linear-gradient(135deg, #6a1b9a, #7b1fa2);
}

.banner-shipped {
  background: linear-gradient(135deg, #00796b, #0097a7);
}

.banner-completed {
  background: linear-gradient(135deg, #2e7d32, #388e3c);
}

.banner-cancelled {
  background: linear-gradient(135deg, #616161, #757575);
}

.status-text {
  font-size: 36rpx;
  font-weight: 700;
  color: #fff;
}

.status-desc {
  font-size: 24rpx;
  color: rgba(255, 255, 255, 0.85);
}

.section {
  background-color: #fff;
  margin: 20rpx 24rpx;
  border-radius: 16rpx;
  padding: 28rpx;
  box-shadow: 0 2rpx 12rpx rgba(0, 0, 0, 0.04);
}

.section-title {
  margin-bottom: 20rpx;
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
  align-items: flex-start;
  justify-content: space-between;
  padding: 14rpx 0;
}

.info-label {
  font-size: 26rpx;
  color: #999;
  flex-shrink: 0;
  margin-right: 20rpx;
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

.address-text {
  max-width: 400rpx;
  text-align: right;
}

.copy-btn {
  font-size: 22rpx;
  color: #0f3a57;
  padding: 4rpx 12rpx;
  border: 1rpx solid #0f3a57;
  border-radius: 16rpx;
  flex-shrink: 0;
}

/* Package items */
.pkg-item {
  padding: 16rpx 0;
  border-bottom: 1rpx solid #f8f8f8;
}

.pkg-item:last-child {
  border-bottom: none;
}

.pkg-main {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 4rpx;
}

.pkg-tracking {
  font-size: 28rpx;
  font-weight: 500;
  color: #333;
}

.pkg-weight {
  font-size: 24rpx;
  color: #999;
}

.pkg-desc {
  font-size: 24rpx;
  color: #999;
}

/* Price */
.price-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12rpx 0;
}

.price-label {
  font-size: 26rpx;
  color: #666;
}

.price-value {
  font-size: 26rpx;
  color: #333;
}

.price-value.discount {
  color: #388e3c;
}

.price-row.total {
  border-top: 1rpx solid #f0f0f0;
  margin-top: 12rpx;
  padding-top: 20rpx;
}

.total-label {
  font-size: 28rpx;
  font-weight: 600;
  color: #333;
}

.total-value {
  font-size: 36rpx;
  font-weight: 700;
  color: #e74c3c;
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

.action-btn.outline {
  border: 2rpx solid #ddd;
}

.action-btn-text {
  font-size: 30rpx;
  font-weight: 500;
}

.primary-text {
  color: #fff;
}

.outline-text {
  color: #666;
}
</style>
