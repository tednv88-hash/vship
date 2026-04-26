<template>
  <view class="recharge-page">
    <!-- Balance card -->
    <view class="balance-card">
      <text class="balance-label">當前餘額（元）</text>
      <text class="balance-amount">{{ balance }}</text>
    </view>

    <!-- Preset amounts -->
    <view class="section">
      <text class="section-title">充值金額</text>
      <view class="amount-grid">
        <view
          v-for="item in presetAmounts"
          :key="item"
          class="amount-item"
          :class="{ active: selectedAmount === item }"
          @tap="selectAmount(item)"
        >
          <text class="amount-value">¥{{ item }}</text>
        </view>
        <view
          class="amount-item"
          :class="{ active: isCustom }"
          @tap="selectCustom"
        >
          <text class="amount-value">自定義</text>
        </view>
      </view>
      <view v-if="isCustom" class="custom-input-wrap">
        <text class="currency-symbol">¥</text>
        <input
          v-model="customAmount"
          type="digit"
          placeholder="請輸入金額"
          class="custom-input"
          @input="onCustomInput"
        />
      </view>
    </view>

    <!-- Payment method -->
    <view class="section">
      <text class="section-title">支付方式</text>
      <view class="payment-list">
        <view
          class="payment-item"
          @tap="paymentMethod = 'wechat'"
        >
          <view class="payment-left">
            <view class="payment-icon wechat-icon-bg">
              <text class="pay-icon-text">微</text>
            </view>
            <text class="payment-name">微信支付</text>
          </view>
          <view
            class="radio"
            :class="{ checked: paymentMethod === 'wechat' }"
          />
        </view>
        <view
          class="payment-item"
          @tap="paymentMethod = 'alipay'"
        >
          <view class="payment-left">
            <view class="payment-icon alipay-icon-bg">
              <text class="pay-icon-text">支</text>
            </view>
            <text class="payment-name">支付寶</text>
          </view>
          <view
            class="radio"
            :class="{ checked: paymentMethod === 'alipay' }"
          />
        </view>
      </view>
    </view>

    <!-- Recharge button -->
    <view class="bottom-bar">
      <view class="recharge-btn" @tap="handleRecharge">
        <text class="recharge-btn-text">
          立即充值 ¥{{ finalAmount || '0' }}
        </text>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { t } from '@/locale'
import { userApi } from '@/api/user'

const balance = ref('0.00')
const presetAmounts = [50, 100, 200, 500, 1000]
const selectedAmount = ref<number | null>(100)
const isCustom = ref(false)
const customAmount = ref('')
const paymentMethod = ref('wechat')
const loading = ref(false)

const finalAmount = computed(() => {
  if (isCustom.value) {
    return parseFloat(customAmount.value) || 0
  }
  return selectedAmount.value || 0
})

function selectAmount(amount: number) {
  selectedAmount.value = amount
  isCustom.value = false
  customAmount.value = ''
}

function selectCustom() {
  isCustom.value = true
  selectedAmount.value = null
}

function onCustomInput() {
  // Sanitize input
  customAmount.value = customAmount.value.replace(/[^\d.]/g, '')
}

async function loadBalance() {
  try {
    const res = (await userApi.getUserInfo()) as any
    const data = res?.data || res
    balance.value = data?.balance || '0.00'
  } catch {
    // ignore
  }
}

async function handleRecharge() {
  if (!finalAmount.value || finalAmount.value <= 0) {
    uni.showToast({ title: '請選擇充值金額', icon: 'none' })
    return
  }
  if (loading.value) return
  loading.value = true

  try {
    const res = (await userApi.recharge({
      amount: finalAmount.value,
      payment_method: paymentMethod.value,
    })) as any
    const data = res?.data || res

    // Invoke payment
    if (paymentMethod.value === 'wechat' && data?.payment_params) {
      await uni.requestPayment({
        provider: 'wxpay',
        ...data.payment_params,
      })
      uni.showToast({ title: '充值成功', icon: 'success' })
      loadBalance()
    } else {
      uni.showToast({ title: '充值請求已提交', icon: 'success' })
    }
  } catch {
    uni.showToast({ title: '充值失敗', icon: 'none' })
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadBalance()
})
</script>

<style scoped>
.recharge-page {
  min-height: 100vh;
  background: #f7f8fa;
  padding-bottom: 180rpx;
}

.balance-card {
  background: linear-gradient(135deg, #0f3a57, #1a5f8a);
  margin: 24rpx;
  border-radius: 20rpx;
  padding: 50rpx 40rpx;
}

.balance-label {
  font-size: 26rpx;
  color: rgba(255, 255, 255, 0.7);
  display: block;
  margin-bottom: 16rpx;
}

.balance-amount {
  font-size: 64rpx;
  font-weight: 700;
  color: #fff;
}

.section {
  margin: 24rpx;
  background: #fff;
  border-radius: 16rpx;
  padding: 30rpx;
}

.section-title {
  font-size: 30rpx;
  font-weight: 600;
  color: #333;
  display: block;
  margin-bottom: 24rpx;
}

.amount-grid {
  display: flex;
  flex-wrap: wrap;
  gap: 20rpx;
}

.amount-item {
  width: calc(33.33% - 14rpx);
  height: 100rpx;
  background: #f7f8fa;
  border-radius: 12rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  border: 2rpx solid transparent;
  box-sizing: border-box;
}

.amount-item.active {
  background: rgba(15, 58, 87, 0.08);
  border-color: #0f3a57;
}

.amount-value {
  font-size: 30rpx;
  color: #333;
  font-weight: 500;
}

.amount-item.active .amount-value {
  color: #0f3a57;
  font-weight: 600;
}

.custom-input-wrap {
  display: flex;
  align-items: center;
  margin-top: 24rpx;
  background: #f7f8fa;
  border-radius: 12rpx;
  padding: 0 24rpx;
  height: 88rpx;
}

.currency-symbol {
  font-size: 36rpx;
  font-weight: 600;
  color: #333;
  margin-right: 12rpx;
}

.custom-input {
  flex: 1;
  height: 88rpx;
  font-size: 32rpx;
  color: #333;
}

.payment-list {
  width: 100%;
}

.payment-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 24rpx 0;
  border-bottom: 1rpx solid #f5f5f5;
}

.payment-item:last-child {
  border-bottom: none;
}

.payment-left {
  display: flex;
  align-items: center;
}

.payment-icon {
  width: 56rpx;
  height: 56rpx;
  border-radius: 12rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: 20rpx;
}

.wechat-icon-bg {
  background: #07c160;
}

.alipay-icon-bg {
  background: #1677ff;
}

.pay-icon-text {
  font-size: 28rpx;
  color: #fff;
  font-weight: 600;
}

.payment-name {
  font-size: 30rpx;
  color: #333;
}

.radio {
  width: 40rpx;
  height: 40rpx;
  border-radius: 50%;
  border: 3rpx solid #ddd;
  box-sizing: border-box;
}

.radio.checked {
  border-color: #0f3a57;
  background: #0f3a57;
  position: relative;
}

.radio.checked::after {
  content: '';
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  width: 16rpx;
  height: 16rpx;
  border-radius: 50%;
  background: #fff;
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

.recharge-btn {
  width: 100%;
  height: 96rpx;
  background: #0f3a57;
  border-radius: 16rpx;
  display: flex;
  align-items: center;
  justify-content: center;
}

.recharge-btn-text {
  font-size: 32rpx;
  color: #fff;
  font-weight: 600;
}
</style>
