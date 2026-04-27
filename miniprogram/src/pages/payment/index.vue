<template>
  <view class="payment-page">
    <!-- Order summary -->
    <view class="summary-card">
      <text class="card-title">訂單資訊</text>
      <view class="summary-row">
        <text class="summary-label">訂單號</text>
        <text class="summary-value">{{ orderInfo.order_no }}</text>
      </view>
      <view v-if="orderInfo.goods_name" class="summary-row">
        <text class="summary-label">商品名稱</text>
        <text class="summary-value">{{ orderInfo.goods_name }}</text>
      </view>
      <view v-if="orderInfo.quantity" class="summary-row">
        <text class="summary-label">數量</text>
        <text class="summary-value">x{{ orderInfo.quantity }}</text>
      </view>
      <view class="summary-row">
        <text class="summary-label">商品金額</text>
        <text class="summary-value">¥{{ orderInfo.goods_amount || '0.00' }}</text>
      </view>
      <view v-if="orderInfo.shipping_fee" class="summary-row">
        <text class="summary-label">運費</text>
        <text class="summary-value">¥{{ orderInfo.shipping_fee }}</text>
      </view>
      <view v-if="orderInfo.discount" class="summary-row">
        <text class="summary-label">優惠</text>
        <text class="summary-value discount">-¥{{ orderInfo.discount }}</text>
      </view>
      <view class="divider" />
      <view class="summary-row total">
        <text class="summary-label">應付金額</text>
        <text class="total-amount">¥{{ orderInfo.total_amount || '0.00' }}</text>
      </view>
    </view>

    <!-- Payment method -->
    <view class="method-card">
      <text class="card-title">支付方式</text>
      <view class="method-list">
        <view
          v-for="method in paymentMethods"
          :key="method.key"
          class="method-item"
          @click="selectedMethod = method.key"
        >
          <view class="method-left">
            <view class="method-icon" :style="{ background: method.bgColor }">
              <uni-icons :type="method.icon" size="22" color="#fff" />
            </view>
            <view class="method-info">
              <text class="method-name">{{ method.label }}</text>
              <text v-if="method.desc" class="method-desc">{{ method.desc }}</text>
            </view>
          </view>
          <view class="radio" :class="{ checked: selectedMethod === method.key }" />
        </view>
      </view>
    </view>

    <!-- Bottom padding -->
    <view style="height: 140rpx" />

    <!-- Bottom bar -->
    <view class="bottom-bar">
      <view class="pay-amount">
        <text class="pay-label">需支付</text>
        <text class="pay-value">¥{{ orderInfo.total_amount || '0.00' }}</text>
      </view>
      <view class="pay-btn" :class="{ disabled: paying }" @click="confirmPay">
        <text>{{ paying ? '支付中...' : '確認支付' }}</text>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { onLoad } from '@dcloudio/uni-app'
import { t } from '@/locale'
import { orderApi } from '@/api/order'

const orderInfo = ref<any>({})
const selectedMethod = ref('wechat')
const paying = ref(false)
const orderId = ref('')
const orderType = ref<'shop' | 'consol'>('shop')

const paymentMethods = computed(() => [
  {
    key: 'wechat',
    label: '微信支付',
    desc: '',
    icon: 'chatbubble-filled',
    bgColor: '#07c160',
  },
  {
    key: 'alipay',
    label: '支付寶',
    desc: '',
    icon: 'wallet-filled',
    bgColor: '#1677ff',
  },
  {
    key: 'balance',
    label: '餘額支付',
    desc: `餘額: ¥${orderInfo.value.user_balance || '0.00'}`,
    icon: 'creditcard-filled',
    bgColor: '#0f3a57',
  },
])

onLoad((options: any) => {
  orderType.value = (options?.type === 'consol' ? 'consol' : 'shop')
  if (options?.order_id) {
    orderId.value = options.order_id
    loadOrder(options.order_id)
  } else {
    uni.showToast({ title: '訂單ID缺失', icon: 'none' })
  }
})

async function loadOrder(id: string) {
  uni.showLoading({ title: t('common.loading') })
  try {
    const res: any = orderType.value === 'shop'
      ? await orderApi.getShopOrderDetail(id)
      : await orderApi.getDetail(id)
    const data = res?.data ?? res ?? {}
    // Normalize fields used by the template
    orderInfo.value = {
      ...data,
      goods_amount: data.goods_total || data.goods_amount || data.total_amount || '0.00',
      total_amount: data.total_price || data.pay_amount || data.total_amount || '0.00',
      shipping_fee: data.shipping_fee || '0.00',
      goods_name: (data.products && data.products[0]?.name) || data.goods_name || '',
      quantity: (data.products && data.products.reduce((s: number, p: any) => s + (p.quantity || 0), 0)) || data.quantity || '',
    }
  } catch (e) {
    console.error(e)
    uni.showToast({ title: '加載失敗', icon: 'none' })
  } finally {
    uni.hideLoading()
  }
}

async function confirmPay() {
  if (paying.value) return
  if (!orderId.value) {
    uni.showToast({ title: '訂單信息異常', icon: 'none' })
    return
  }

  // Balance check (only relevant once balance feature is wired)
  if (selectedMethod.value === 'balance') {
    const balance = parseFloat(orderInfo.value.user_balance || '0')
    const total = parseFloat(orderInfo.value.total_amount || '0')
    if (balance < total) {
      uni.showModal({
        title: '餘額不足',
        content: '您的餘額不足，是否前往充值？',
        confirmText: '去充值',
        success(res) {
          if (res.confirm) {
            uni.navigateTo({ url: '/pages/recharge/index' })
          }
        },
      })
      return
    }
  }

  paying.value = true
  try {
    const res: any = orderType.value === 'shop'
      ? await orderApi.payShopOrder(orderId.value, { payment_method: selectedMethod.value })
      : await orderApi.payOrder(orderId.value, { payment_method: selectedMethod.value })

    const data = res?.data ?? res
    // If real wxpay params returned, invoke uni payment
    if (selectedMethod.value === 'wechat' && data?.payment_params) {
      try {
        await uni.requestPayment(data.payment_params)
      } catch (err: any) {
        if (err?.errMsg?.includes('cancel')) {
          uni.showToast({ title: '已取消支付', icon: 'none' })
          return
        }
        throw err
      }
    }

    uni.redirectTo({
      url: `/pages/payment/result?order_id=${orderId.value}&type=${orderType.value}&status=success`,
    })
  } catch (e: any) {
    console.error(e)
    uni.redirectTo({
      url: `/pages/payment/result?order_id=${orderId.value}&type=${orderType.value}&status=fail`,
    })
  } finally {
    paying.value = false
  }
}
</script>

<style scoped>
.payment-page {
  min-height: 100vh;
  background: #f5f5f5;
  padding: 24rpx;
}

.summary-card,
.method-card {
  background: #fff;
  border-radius: 16rpx;
  padding: 28rpx 30rpx;
  margin-bottom: 20rpx;
}

.card-title {
  font-size: 30rpx;
  font-weight: 600;
  color: #333;
  margin-bottom: 20rpx;
  display: block;
}

.summary-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12rpx 0;
}

.summary-label {
  font-size: 26rpx;
  color: #999;
}

.summary-value {
  font-size: 26rpx;
  color: #333;
}

.summary-value.discount {
  color: #4caf50;
}

.divider {
  height: 1rpx;
  background: #f0f0f0;
  margin: 12rpx 0;
}

.summary-row.total {
  padding-top: 16rpx;
}

.total-amount {
  font-size: 36rpx;
  color: #e64340;
  font-weight: 700;
}

.method-list {
  display: flex;
  flex-direction: column;
  gap: 16rpx;
}

.method-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 20rpx 16rpx;
  background: #f8f8f8;
  border-radius: 12rpx;
}

.method-left {
  display: flex;
  align-items: center;
  gap: 20rpx;
}

.method-icon {
  width: 52rpx;
  height: 52rpx;
  border-radius: 12rpx;
  display: flex;
  align-items: center;
  justify-content: center;
}

.method-info {
  display: flex;
  flex-direction: column;
  gap: 4rpx;
}

.method-name {
  font-size: 28rpx;
  color: #333;
}

.method-desc {
  font-size: 22rpx;
  color: #999;
}

.radio {
  width: 36rpx;
  height: 36rpx;
  border: 2rpx solid #ddd;
  border-radius: 50%;
}

.radio.checked {
  border-color: #0f3a57;
  background: #0f3a57;
  box-shadow: inset 0 0 0 4rpx #fff;
}

.bottom-bar {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  display: flex;
  align-items: center;
  background: #fff;
  padding: 16rpx 30rpx;
  padding-bottom: calc(16rpx + env(safe-area-inset-bottom));
  box-shadow: 0 -2rpx 10rpx rgba(0, 0, 0, 0.05);
  z-index: 100;
}

.pay-amount {
  flex: 1;
  display: flex;
  align-items: baseline;
  gap: 8rpx;
}

.pay-label {
  font-size: 26rpx;
  color: #666;
}

.pay-value {
  font-size: 40rpx;
  color: #e64340;
  font-weight: 700;
}

.pay-btn {
  background: #0f3a57;
  border-radius: 44rpx;
  padding: 0 60rpx;
  height: 84rpx;
  display: flex;
  align-items: center;
  justify-content: center;
}

.pay-btn.disabled {
  opacity: 0.6;
}

.pay-btn text {
  color: #fff;
  font-size: 30rpx;
  font-weight: 500;
}
</style>
