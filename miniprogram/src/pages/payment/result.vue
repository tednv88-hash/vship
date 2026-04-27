<template>
  <view class="result-page">
    <view class="result-card">
      <!-- Icon -->
      <view class="result-icon" :class="isSuccess ? 'success' : 'fail'">
        <uni-icons
          :type="isSuccess ? 'checkmarkempty' : 'closeempty'"
          size="48"
          color="#fff"
        />
      </view>

      <!-- Message -->
      <text class="result-title">{{ isSuccess ? '支付成功' : '支付失敗' }}</text>
      <text class="result-desc">
        {{ isSuccess ? '您的訂單已支付成功' : '支付未完成，請重新嘗試' }}
      </text>

      <!-- Order info -->
      <view class="info-section">
        <view class="info-row">
          <text class="info-label">訂單號</text>
          <text class="info-value">{{ orderInfo.order_no || orderId }}</text>
        </view>
        <view v-if="isSuccess && orderInfo.total_amount" class="info-row">
          <text class="info-label">支付金額</text>
          <text class="info-value amount">¥{{ orderInfo.total_amount }}</text>
        </view>
        <view v-if="isSuccess && orderInfo.payment_method" class="info-row">
          <text class="info-label">支付方式</text>
          <text class="info-value">{{ getMethodLabel(orderInfo.payment_method) }}</text>
        </view>
        <view v-if="isSuccess && orderInfo.paid_at" class="info-row">
          <text class="info-label">支付時間</text>
          <text class="info-value">{{ orderInfo.paid_at }}</text>
        </view>
      </view>
    </view>

    <!-- Action buttons -->
    <view class="actions">
      <view class="action-btn primary" @click="goOrderDetail">
        <text>查看訂單</text>
      </view>
      <view class="action-btn secondary" @click="goHome">
        <text>返回首頁</text>
      </view>
      <view v-if="!isSuccess" class="action-btn primary" @click="retryPay">
        <text>重新支付</text>
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
const status = ref('')
const orderType = ref<'shop' | 'consol'>('shop')
const orderInfo = ref<any>({})

const isSuccess = computed(() => status.value === 'success')

onLoad((options: any) => {
  orderId.value = options?.order_id || ''
  status.value = options?.status || 'fail'
  orderType.value = options?.type === 'consol' ? 'consol' : 'shop'
  if (orderId.value) {
    loadOrderInfo()
  }
})

async function loadOrderInfo() {
  try {
    const res: any = orderType.value === 'shop'
      ? await orderApi.getShopOrderDetail(orderId.value)
      : await orderApi.getDetail(orderId.value)
    const data = res?.data ?? res ?? {}
    orderInfo.value = {
      ...data,
      total_amount: data.total_price || data.pay_amount || data.total_amount || '',
      payment_method: data.pay_method || data.payment_method || '',
    }
  } catch (e) {
    console.error(e)
  }
}

function getMethodLabel(method: string): string {
  const map: Record<string, string> = {
    wechat: '微信支付',
    wxpay: '微信支付',
    alipay: '支付寶',
    balance: '餘額支付',
  }
  return map[method] || method
}

function goOrderDetail() {
  const path = orderType.value === 'shop' ? '/pages/shop-order/detail' : '/pages/order/detail'
  uni.redirectTo({ url: `${path}?id=${orderId.value}` })
}

function goHome() {
  uni.switchTab({ url: '/pages/index/index' })
}

function retryPay() {
  uni.redirectTo({ url: `/pages/payment/index?order_id=${orderId.value}&type=${orderType.value}` })
}
</script>

<style scoped>
.result-page {
  min-height: 100vh;
  background: #f5f5f5;
  padding: 40rpx 24rpx;
}

.result-card {
  background: #fff;
  border-radius: 16rpx;
  padding: 60rpx 30rpx 40rpx;
  display: flex;
  flex-direction: column;
  align-items: center;
}

.result-icon {
  width: 120rpx;
  height: 120rpx;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 30rpx;
}

.result-icon.success {
  background: #4caf50;
}

.result-icon.fail {
  background: #e64340;
}

.result-title {
  font-size: 36rpx;
  font-weight: 600;
  color: #333;
  margin-bottom: 12rpx;
}

.result-desc {
  font-size: 26rpx;
  color: #999;
  margin-bottom: 40rpx;
}

.info-section {
  width: 100%;
  background: #f8f8f8;
  border-radius: 12rpx;
  padding: 24rpx 30rpx;
}

.info-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 14rpx 0;
}

.info-label {
  font-size: 26rpx;
  color: #999;
}

.info-value {
  font-size: 26rpx;
  color: #333;
}

.info-value.amount {
  font-size: 30rpx;
  color: #e64340;
  font-weight: 600;
}

.actions {
  margin-top: 40rpx;
  display: flex;
  flex-direction: column;
  gap: 20rpx;
}

.action-btn {
  height: 88rpx;
  border-radius: 44rpx;
  display: flex;
  align-items: center;
  justify-content: center;
}

.action-btn.primary {
  background: #0f3a57;
}

.action-btn.primary text {
  color: #fff;
  font-size: 30rpx;
  font-weight: 500;
}

.action-btn.secondary {
  background: #fff;
  border: 2rpx solid #ddd;
}

.action-btn.secondary text {
  color: #666;
  font-size: 30rpx;
}
</style>
