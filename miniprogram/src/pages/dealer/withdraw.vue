<template>
  <view class="withdraw-page">
    <!-- Balance card -->
    <view class="balance-card">
      <text class="balance-label">可提現餘額</text>
      <text class="balance-value">¥{{ availableBalance }}</text>
    </view>

    <!-- Withdrawal form -->
    <view class="form-card">
      <text class="card-title">申請提現</text>

      <view class="form-item">
        <text class="form-label">提現金額</text>
        <view class="amount-input-wrap">
          <text class="currency-symbol">¥</text>
          <input
            class="amount-input"
            type="digit"
            placeholder="請輸入提現金額"
            v-model="withdrawAmount"
          />
        </view>
        <view class="amount-quick">
          <view
            v-for="amt in quickAmounts"
            :key="amt"
            class="quick-tag"
            @click="withdrawAmount = String(amt)"
          >
            <text>¥{{ amt }}</text>
          </view>
          <view class="quick-tag" @click="withdrawAmount = availableBalance">
            <text>全部</text>
          </view>
        </view>
      </view>

      <view class="form-item">
        <text class="form-label">提現方式</text>
        <view class="method-list">
          <view
            v-for="method in withdrawMethods"
            :key="method.key"
            class="method-item"
            :class="{ selected: selectedMethod === method.key }"
            @click="selectedMethod = method.key"
          >
            <view class="method-icon">
              <uni-icons :type="method.icon" size="22" :color="selectedMethod === method.key ? '#0f3a57' : '#666'" />
            </view>
            <text class="method-name">{{ method.label }}</text>
            <view class="radio" :class="{ checked: selectedMethod === method.key }" />
          </view>
        </view>
      </view>

      <view class="submit-btn" @click="handleWithdraw">
        <text>確認提現</text>
      </view>
    </view>

    <!-- Withdrawal history -->
    <view class="section-card">
      <text class="card-title">提現記錄</text>
      <view v-if="historyList.length > 0" class="history-list">
        <view v-for="item in historyList" :key="item.id" class="history-item">
          <view class="history-left">
            <text class="history-amount">-¥{{ item.amount }}</text>
            <text class="history-method">{{ item.method_label }}</text>
          </view>
          <view class="history-right">
            <text class="history-status" :class="item.status">{{ getStatusLabel(item.status) }}</text>
            <text class="history-date">{{ item.created_at }}</text>
          </view>
        </view>
      </view>
      <view v-else class="empty-text">
        <text>{{ t('common.noData') }}</text>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { t } from '@/locale'
import { commonApi } from '@/api/common'

const availableBalance = ref('0.00')
const withdrawAmount = ref('')
const selectedMethod = ref('wechat')
const historyList = ref<any[]>([])

const quickAmounts = [100, 500, 1000, 2000]

const withdrawMethods = [
  { key: 'wechat', label: '微信零錢', icon: 'chatbubble' },
  { key: 'alipay', label: '支付寶', icon: 'wallet' },
  { key: 'bank', label: '銀行卡', icon: 'creditcard' },
]

onMounted(() => {
  loadInfo()
  loadHistory()
})

async function loadInfo() {
  try {
    const res = await commonApi.getDealerInfo()
    const data = res?.data || {}
    availableBalance.value = data.available_balance || '0.00'
  } catch (e) {
    console.error(e)
  }
}

async function loadHistory() {
  try {
    const res = await commonApi.getDealerWithdrawals({ page: 1, limit: 20 })
    historyList.value = res?.data?.list || res?.data || []
  } catch (e) {
    console.error(e)
  }
}

async function handleWithdraw() {
  const amount = parseFloat(withdrawAmount.value)
  if (!amount || amount <= 0) {
    uni.showToast({ title: '請輸入正確金額', icon: 'none' })
    return
  }
  if (amount > parseFloat(availableBalance.value)) {
    uni.showToast({ title: '餘額不足', icon: 'none' })
    return
  }

  uni.showModal({
    title: '確認提現',
    content: `確定提現 ¥${amount.toFixed(2)} 到${withdrawMethods.find((m) => m.key === selectedMethod.value)?.label}嗎？`,
    success: async (res) => {
      if (res.confirm) {
        uni.showLoading({ title: t('common.loading') })
        try {
          await commonApi.requestWithdraw({
            amount,
            method: selectedMethod.value,
          })
          uni.showToast({ title: '提現申請已提交', icon: 'success' })
          withdrawAmount.value = ''
          loadInfo()
          loadHistory()
        } catch (e) {
          console.error(e)
          uni.showToast({ title: '提現失敗', icon: 'none' })
        } finally {
          uni.hideLoading()
        }
      }
    },
  })
}

function getStatusLabel(status: string): string {
  const map: Record<string, string> = {
    pending: '處理中',
    approved: '已到賬',
    rejected: '已拒絕',
  }
  return map[status] || status
}
</script>

<style scoped>
.withdraw-page {
  min-height: 100vh;
  background: #f5f5f5;
}

.balance-card {
  background: linear-gradient(135deg, #0f3a57 0%, #1a5c7a 100%);
  padding: 50rpx 30rpx;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 16rpx;
}

.balance-label {
  font-size: 26rpx;
  color: rgba(255, 255, 255, 0.7);
}

.balance-value {
  font-size: 56rpx;
  font-weight: 700;
  color: #fff;
}

.form-card,
.section-card {
  background: #fff;
  padding: 30rpx;
  margin: 16rpx 0;
}

.card-title {
  font-size: 30rpx;
  font-weight: 600;
  color: #333;
  margin-bottom: 24rpx;
  display: block;
}

.form-item {
  margin-bottom: 30rpx;
}

.form-label {
  display: block;
  font-size: 26rpx;
  color: #666;
  margin-bottom: 16rpx;
}

.amount-input-wrap {
  display: flex;
  align-items: center;
  border-bottom: 2rpx solid #0f3a57;
  padding: 16rpx 0;
}

.currency-symbol {
  font-size: 40rpx;
  font-weight: 600;
  color: #333;
  margin-right: 12rpx;
}

.amount-input {
  flex: 1;
  font-size: 48rpx;
  font-weight: 600;
  color: #333;
}

.amount-quick {
  display: flex;
  gap: 16rpx;
  margin-top: 20rpx;
}

.quick-tag {
  padding: 10rpx 24rpx;
  background: #f5f5f5;
  border-radius: 8rpx;
}

.quick-tag text {
  font-size: 24rpx;
  color: #666;
}

.method-list {
  display: flex;
  flex-direction: column;
  gap: 16rpx;
}

.method-item {
  display: flex;
  align-items: center;
  padding: 24rpx 20rpx;
  background: #f8f8f8;
  border-radius: 12rpx;
  border: 2rpx solid transparent;
  gap: 16rpx;
}

.method-item.selected {
  border-color: #0f3a57;
  background: rgba(15, 58, 87, 0.04);
}

.method-icon {
  width: 48rpx;
  height: 48rpx;
  display: flex;
  align-items: center;
  justify-content: center;
}

.method-name {
  flex: 1;
  font-size: 28rpx;
  color: #333;
}

.radio {
  width: 32rpx;
  height: 32rpx;
  border: 2rpx solid #ddd;
  border-radius: 50%;
}

.radio.checked {
  border-color: #0f3a57;
  background: #0f3a57;
  box-shadow: inset 0 0 0 4rpx #fff;
}

.submit-btn {
  background: #0f3a57;
  border-radius: 44rpx;
  height: 88rpx;
  display: flex;
  align-items: center;
  justify-content: center;
}

.submit-btn text {
  color: #fff;
  font-size: 30rpx;
  font-weight: 500;
}

.history-list {
  margin-top: 10rpx;
}

.history-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 24rpx 0;
  border-bottom: 1rpx solid #f5f5f5;
}

.history-item:last-child {
  border-bottom: none;
}

.history-left {
  display: flex;
  flex-direction: column;
  gap: 8rpx;
}

.history-amount {
  font-size: 30rpx;
  color: #333;
  font-weight: 500;
}

.history-method {
  font-size: 22rpx;
  color: #999;
}

.history-right {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  gap: 8rpx;
}

.history-status {
  font-size: 24rpx;
  font-weight: 500;
}

.history-status.pending {
  color: #ff9800;
}

.history-status.approved {
  color: #4caf50;
}

.history-status.rejected {
  color: #e64340;
}

.history-date {
  font-size: 22rpx;
  color: #999;
}

.empty-text {
  text-align: center;
  padding: 40rpx 0;
  color: #999;
  font-size: 26rpx;
}
</style>
