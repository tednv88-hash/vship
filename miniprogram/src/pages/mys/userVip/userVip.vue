<template>
  <view class="vip-page">
    <!-- Current VIP status -->
    <view class="vip-header">
      <view class="vip-card">
        <view class="vip-card-top">
          <view class="vip-avatar">
            <text class="vip-avatar-text">V</text>
          </view>
          <view class="vip-info">
            <text class="vip-level">
              {{ currentLevel ? `VIP ${currentLevel}` : '普通用戶' }}
            </text>
            <text class="vip-expire" v-if="expireDate">
              有效期至 {{ expireDate }}
            </text>
            <text class="vip-expire" v-else>尚未開通會員</text>
          </view>
        </view>
      </view>
    </view>

    <!-- VIP benefits -->
    <view class="section">
      <text class="section-title">會員專屬權益</text>
      <view class="benefits-grid">
        <view v-for="benefit in benefits" :key="benefit.title" class="benefit-item">
          <view class="benefit-icon-wrap">
            <text class="benefit-icon">{{ benefit.icon }}</text>
          </view>
          <text class="benefit-title">{{ benefit.title }}</text>
          <text class="benefit-desc">{{ benefit.desc }}</text>
        </view>
      </view>
    </view>

    <!-- Pricing plans -->
    <view class="section">
      <text class="section-title">選擇套餐</text>
      <view class="plans-list">
        <view
          v-for="plan in plans"
          :key="plan.id"
          class="plan-card"
          :class="{ active: selectedPlan === plan.id, recommended: plan.recommended }"
          @tap="selectedPlan = plan.id"
        >
          <view v-if="plan.recommended" class="recommend-tag">
            <text class="recommend-text">推薦</text>
          </view>
          <text class="plan-name">{{ plan.name }}</text>
          <view class="plan-price-row">
            <text class="plan-currency">¥</text>
            <text class="plan-price">{{ plan.price }}</text>
          </view>
          <text class="plan-original" v-if="plan.original_price">
            原價 ¥{{ plan.original_price }}
          </text>
          <text class="plan-desc">{{ plan.desc }}</text>
        </view>
      </view>
    </view>

    <!-- Purchase button -->
    <view class="bottom-bar">
      <view class="purchase-btn" @tap="handlePurchase">
        <text class="purchase-btn-text">立即開通</text>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { t } from '@/locale'
import { userApi } from '@/api/user'

const currentLevel = ref(0)
const expireDate = ref('')
const selectedPlan = ref('')
const loading = ref(false)

const benefits = [
  { icon: '💰', title: '運費折扣', desc: '最高享8折' },
  { icon: '📦', title: '免費倉儲', desc: '延長保管期' },
  { icon: '⚡', title: '優先出貨', desc: '優先處理' },
  { icon: '🎁', title: '積分加倍', desc: '雙倍積分' },
  { icon: '🎫', title: '專屬優惠', desc: '會員專享券' },
  { icon: '💎', title: '專屬客服', desc: '一對一服務' },
]

interface PricePlan {
  id: string
  name: string
  price: number
  original_price?: number
  desc: string
  recommended?: boolean
}

const plans = ref<PricePlan[]>([
  {
    id: 'monthly',
    name: '月卡',
    price: 29,
    original_price: 49,
    desc: '30天VIP會員',
  },
  {
    id: 'quarterly',
    name: '季卡',
    price: 79,
    original_price: 147,
    desc: '90天VIP會員',
    recommended: true,
  },
  {
    id: 'yearly',
    name: '年卡',
    price: 269,
    original_price: 588,
    desc: '365天VIP會員',
  },
])

async function loadVipInfo() {
  try {
    const res = (await userApi.getUserInfo()) as any
    const data = res?.data || res
    currentLevel.value = data?.vip_level || 0
    expireDate.value = data?.vip_expire_date || ''
  } catch {
    // ignore
  }
}

async function handlePurchase() {
  if (!selectedPlan.value) {
    uni.showToast({ title: '請選擇套餐', icon: 'none' })
    return
  }
  if (loading.value) return
  loading.value = true

  try {
    const res = (await userApi.recharge({
      amount: plans.value.find((p) => p.id === selectedPlan.value)?.price || 0,
      payment_method: 'wechat',
    })) as any
    const data = res?.data || res
    if (data?.payment_params) {
      await uni.requestPayment({
        provider: 'wxpay',
        ...data.payment_params,
      })
      uni.showToast({ title: '開通成功', icon: 'success' })
      loadVipInfo()
    }
  } catch {
    uni.showToast({ title: '購買失敗', icon: 'none' })
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  selectedPlan.value = 'quarterly'
  loadVipInfo()
})
</script>

<style scoped>
.vip-page {
  min-height: 100vh;
  background: #f7f8fa;
  padding-bottom: 180rpx;
}

.vip-header {
  background: linear-gradient(135deg, #0f3a57, #1a5f8a);
  padding: 60rpx 30rpx 50rpx;
}

.vip-card {
  width: 100%;
}

.vip-card-top {
  display: flex;
  align-items: center;
}

.vip-avatar {
  width: 100rpx;
  height: 100rpx;
  background: rgba(255, 255, 255, 0.2);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: 24rpx;
  border: 3rpx solid rgba(255, 255, 255, 0.4);
}

.vip-avatar-text {
  font-size: 48rpx;
  color: #fff;
  font-weight: 700;
}

.vip-info {
  flex: 1;
}

.vip-level {
  display: block;
  font-size: 36rpx;
  font-weight: 700;
  color: #fff;
  margin-bottom: 8rpx;
}

.vip-expire {
  font-size: 24rpx;
  color: rgba(255, 255, 255, 0.7);
}

.section {
  margin: 24rpx;
  background: #fff;
  border-radius: 16rpx;
  padding: 30rpx;
}

.section-title {
  font-size: 32rpx;
  font-weight: 600;
  color: #333;
  display: block;
  margin-bottom: 24rpx;
}

.benefits-grid {
  display: flex;
  flex-wrap: wrap;
  gap: 16rpx;
}

.benefit-item {
  width: calc(33.33% - 12rpx);
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 20rpx 8rpx;
  background: #f9fafb;
  border-radius: 12rpx;
  box-sizing: border-box;
}

.benefit-icon-wrap {
  width: 64rpx;
  height: 64rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 8rpx;
}

.benefit-icon {
  font-size: 40rpx;
}

.benefit-title {
  font-size: 24rpx;
  color: #333;
  font-weight: 500;
  margin-bottom: 4rpx;
}

.benefit-desc {
  font-size: 20rpx;
  color: #999;
}

.plans-list {
  display: flex;
  gap: 16rpx;
}

.plan-card {
  flex: 1;
  background: #f9fafb;
  border-radius: 16rpx;
  padding: 28rpx 16rpx;
  display: flex;
  flex-direction: column;
  align-items: center;
  border: 3rpx solid transparent;
  position: relative;
  box-sizing: border-box;
}

.plan-card.active {
  border-color: #0f3a57;
  background: rgba(15, 58, 87, 0.04);
}

.plan-card.recommended {
  border-color: #0f3a57;
}

.recommend-tag {
  position: absolute;
  top: -2rpx;
  right: -2rpx;
  background: #e74c3c;
  padding: 4rpx 14rpx;
  border-radius: 0 14rpx 0 12rpx;
}

.recommend-text {
  font-size: 18rpx;
  color: #fff;
}

.plan-name {
  font-size: 28rpx;
  color: #333;
  font-weight: 600;
  margin-bottom: 12rpx;
}

.plan-price-row {
  display: flex;
  align-items: baseline;
  margin-bottom: 6rpx;
}

.plan-currency {
  font-size: 24rpx;
  color: #0f3a57;
  font-weight: 600;
}

.plan-price {
  font-size: 48rpx;
  font-weight: 700;
  color: #0f3a57;
}

.plan-original {
  font-size: 20rpx;
  color: #999;
  text-decoration: line-through;
  margin-bottom: 8rpx;
}

.plan-desc {
  font-size: 20rpx;
  color: #666;
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

.purchase-btn {
  width: 100%;
  height: 96rpx;
  background: linear-gradient(135deg, #0f3a57, #1a5f8a);
  border-radius: 16rpx;
  display: flex;
  align-items: center;
  justify-content: center;
}

.purchase-btn-text {
  font-size: 32rpx;
  color: #fff;
  font-weight: 600;
  letter-spacing: 2rpx;
}
</style>
