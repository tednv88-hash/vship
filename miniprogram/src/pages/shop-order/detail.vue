<template>
  <view class="page">
    <!-- Status banner -->
    <view class="status-banner" :class="'banner-' + detail.status">
      <text class="status-text">{{ getStatusLabel(detail.status) }}</text>
      <text class="status-desc">{{ statusDesc }}</text>
    </view>

    <!-- Shipping info (for shipped orders) -->
    <view v-if="detail.shipping" class="section shipping-section" @tap="goTracking">
      <view class="shipping-header">
        <text class="shipping-courier">{{ detail.shipping.courier_name }}</text>
        <text class="shipping-arrow">&#8250;</text>
      </view>
      <view class="shipping-info">
        <text class="shipping-no">{{ detail.shipping.tracking_no }}</text>
        <text class="shipping-status">{{ detail.shipping.latest_event }}</text>
      </view>
    </view>

    <!-- Product list -->
    <view class="section">
      <view class="section-title"><text>商品資訊</text></view>
      <view v-for="prod in detail.products" :key="prod.id" class="product-item">
        <image class="product-img" :src="prod.image" mode="aspectFill" />
        <view class="product-info">
          <text class="product-name">{{ prod.name }}</text>
          <text class="product-spec" v-if="prod.spec">{{ prod.spec }}</text>
          <view class="product-bottom">
            <text class="product-price">¥{{ prod.price }}</text>
            <text class="product-qty">x{{ prod.quantity }}</text>
          </view>
        </view>
      </view>
    </view>

    <!-- Delivery address -->
    <view class="section">
      <view class="section-title"><text>收件資訊</text></view>
      <view class="info-row">
        <text class="info-label">收件人</text>
        <text class="info-value">{{ detail.address?.name }}</text>
      </view>
      <view class="info-row">
        <text class="info-label">電話</text>
        <text class="info-value">{{ detail.address?.phone }}</text>
      </view>
      <view class="info-row">
        <text class="info-label">地址</text>
        <text class="info-value address-text">{{ detail.address?.region }} {{ detail.address?.detail }}</text>
      </view>
    </view>

    <!-- Price breakdown -->
    <view class="section">
      <view class="section-title"><text>費用明細</text></view>
      <view class="price-row">
        <text class="price-label">商品金額</text>
        <text class="price-value">¥{{ detail.goods_total }}</text>
      </view>
      <view class="price-row">
        <text class="price-label">運費</text>
        <text class="price-value">{{ detail.shipping_fee === '0.00' ? '免運費' : '¥' + detail.shipping_fee }}</text>
      </view>
      <view v-if="detail.coupon_discount" class="price-row">
        <text class="price-label">優惠券</text>
        <text class="price-value discount">-¥{{ detail.coupon_discount }}</text>
      </view>
      <view class="price-row total">
        <text class="price-label total-label">實付金額</text>
        <text class="price-value total-value">¥{{ detail.total_price }}</text>
      </view>
    </view>

    <!-- Order info -->
    <view class="section">
      <view class="section-title"><text>訂單資訊</text></view>
      <view class="info-row">
        <text class="info-label">訂單編號</text>
        <view class="info-value-wrap">
          <text class="info-value">{{ detail.order_no }}</text>
          <text class="copy-btn" @tap="copyText(detail.order_no)">複製</text>
        </view>
      </view>
      <view class="info-row">
        <text class="info-label">下單時間</text>
        <text class="info-value">{{ detail.created_at }}</text>
      </view>
      <view v-if="detail.paid_at" class="info-row">
        <text class="info-label">付款時間</text>
        <text class="info-value">{{ detail.paid_at }}</text>
      </view>
      <view v-if="detail.shipped_at" class="info-row">
        <text class="info-label">發貨時間</text>
        <text class="info-value">{{ detail.shipped_at }}</text>
      </view>
    </view>

    <!-- Action buttons -->
    <view v-if="showActions" class="action-bar">
      <view v-if="detail.status === 'pending'" class="action-btn outline" @tap="cancelOrder">
        <text class="action-btn-text outline-text">取消訂單</text>
      </view>
      <view v-if="detail.status === 'pending'" class="action-btn primary" @tap="payOrder">
        <text class="action-btn-text primary-text">去付款</text>
      </view>
      <view v-if="detail.status === 'shipped'" class="action-btn primary" @tap="confirmReceive">
        <text class="action-btn-text primary-text">確認收貨</text>
      </view>
      <view v-if="detail.status === 'completed'" class="action-btn outline" @tap="goReview">
        <text class="action-btn-text outline-text">去評價</text>
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
  products: [],
  address: {},
  shipping: null,
  goods_total: '0',
  shipping_fee: '0',
  coupon_discount: '',
  total_price: '0',
  created_at: '',
  paid_at: '',
  shipped_at: '',
})

const statusMap: Record<string, string> = {
  pending: t('order.status.pending'),
  paid: '待發貨',
  shipped: '待收貨',
  completed: t('order.status.completed'),
  cancelled: t('order.status.cancelled'),
}

const statusDescMap: Record<string, string> = {
  pending: '請在30分鐘內完成付款',
  paid: '商家正在準備發貨',
  shipped: '包裹正在運送中',
  completed: '訂單已完成',
  cancelled: '訂單已取消',
}

function getStatusLabel(status: string) {
  return statusMap[status] || status
}

const statusDesc = computed(() => statusDescMap[detail.value.status] || '')

const showActions = computed(() => {
  return ['pending', 'shipped', 'completed'].includes(detail.value.status)
})

async function fetchDetail() {
  try {
    const res: any = await orderApi.getShopOrderDetail(orderId.value)
    const data = res?.data ?? res
    if (data && data.id) {
      detail.value = {
        ...detail.value,
        ...data,
        products: data.products || [],
        address: data.address || {},
      }
    }
  } catch (e) {
    uni.showToast({ title: '載入失敗', icon: 'none' })
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

function goTracking() {
  if (detail.value.shipping?.tracking_no) {
    uni.navigateTo({ url: `/pages/tracking/index?no=${detail.value.shipping.tracking_no}` })
  }
}

function cancelOrder() {
  uni.showModal({
    title: '取消訂單',
    content: '確認取消此訂單？',
    success: async (res) => {
      if (res.confirm) {
        try {
          await orderApi.cancelShopOrder(orderId.value)
          uni.showToast({ title: '已取消', icon: 'success' })
          detail.value.status = 'cancelled'
        } catch (e: any) {
          uni.showToast({ title: e?.message || '取消失敗', icon: 'none' })
        }
      }
    },
  })
}

function payOrder() {
  uni.navigateTo({ url: `/pages/payment/index?order_id=${orderId.value}&type=shop` })
}

function confirmReceive() {
  uni.showModal({
    title: '確認收貨',
    content: '確認已收到商品？',
    success: async (res) => {
      if (res.confirm) {
        try {
          await orderApi.confirmShopOrder(orderId.value)
          uni.showToast({ title: '已確認收貨', icon: 'success' })
          detail.value.status = 'completed'
        } catch (e: any) {
          uni.showToast({ title: e?.message || '操作失敗', icon: 'none' })
        }
      }
    },
  })
}

function goReview() {
  uni.navigateTo({ url: `/pages/review/create?order_id=${orderId.value}&type=shop` })
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

/* Shipping card */
.shipping-section {
  cursor: pointer;
}

.shipping-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 12rpx;
}

.shipping-courier {
  font-size: 28rpx;
  font-weight: 600;
  color: #333;
}

.shipping-arrow {
  font-size: 32rpx;
  color: #ccc;
}

.shipping-info {
  display: flex;
  flex-direction: column;
  gap: 6rpx;
}

.shipping-no {
  font-size: 24rpx;
  color: #666;
}

.shipping-status {
  font-size: 26rpx;
  color: #0f3a57;
}

/* Product */
.product-item {
  display: flex;
  gap: 20rpx;
  padding: 16rpx 0;
  border-bottom: 1rpx solid #f8f8f8;
}

.product-item:last-child {
  border-bottom: none;
}

.product-img {
  width: 160rpx;
  height: 160rpx;
  border-radius: 12rpx;
  background-color: #f5f5f5;
  flex-shrink: 0;
}

.product-info {
  flex: 1;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
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

.product-bottom {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.product-price {
  font-size: 28rpx;
  color: #e74c3c;
  font-weight: 500;
}

.product-qty {
  font-size: 24rpx;
  color: #999;
}

/* Info rows */
.info-row {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  padding: 12rpx 0;
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
