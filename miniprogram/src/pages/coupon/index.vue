<template>
  <view class="coupon-page">
    <!-- Tabs -->
    <view class="tab-bar">
      <view
        v-for="tab in tabs"
        :key="tab.value"
        class="tab-item"
        :class="{ active: currentTab === tab.value }"
        @tap="switchTab(tab.value)"
      >
        <text class="tab-text">{{ tab.label }}</text>
      </view>
    </view>

    <!-- Coupon list -->
    <scroll-view scroll-y class="coupon-scroll">
      <view v-if="list.length === 0 && !loading" class="empty">
        <text class="empty-text">{{ t('common.noData') }}</text>
      </view>

      <view v-for="item in list" :key="item.id" class="coupon-card">
        <view class="coupon-left">
          <view class="coupon-amount-row">
            <text class="currency">¥</text>
            <text class="coupon-amount">{{ item.discount }}</text>
          </view>
          <text class="coupon-condition">滿{{ item.min_spend }}可用</text>
        </view>
        <view class="coupon-right">
          <text class="coupon-name">{{ item.name }}</text>
          <text class="coupon-period">
            {{ item.start_date }} - {{ item.end_date }}
          </text>
          <view
            v-if="currentTab === 'available'"
            class="claim-btn"
            :class="{ claimed: item.claimed }"
            @tap="handleClaim(item)"
          >
            <text class="claim-btn-text">
              {{ item.claimed ? '已領取' : '領取' }}
            </text>
          </view>
        </view>
      </view>

      <view v-if="loading" class="loading-more">
        <text class="loading-text">{{ t('common.loading') }}</text>
      </view>
    </scroll-view>
  </view>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { t } from '@/locale'
import { commonApi } from '@/api/common'

interface CouponItem {
  id: string
  name: string
  discount: number
  min_spend: number
  start_date: string
  end_date: string
  claimed?: boolean
}

const tabs = [
  { label: '可領取', value: 'available' },
  { label: '已領取', value: 'claimed' },
]

const currentTab = ref('available')
const list = ref<CouponItem[]>([])
const loading = ref(false)

async function loadList() {
  loading.value = true
  try {
    const params = { status: currentTab.value }
    const res = (await commonApi.getCoupons(params)) as any
    const data = res?.data || res || []
    list.value = Array.isArray(data) ? data : data?.list || []
  } catch {
    uni.showToast({ title: '載入失敗', icon: 'none' })
  } finally {
    loading.value = false
  }
}

function switchTab(value: string) {
  currentTab.value = value
  loadList()
}

async function handleClaim(item: CouponItem) {
  if (item.claimed) return
  try {
    await commonApi.claimCoupon(item.id)
    item.claimed = true
    uni.showToast({ title: '領取成功', icon: 'success' })
  } catch {
    uni.showToast({ title: '領取失敗', icon: 'none' })
  }
}

onMounted(() => {
  loadList()
})
</script>

<style scoped>
.coupon-page {
  min-height: 100vh;
  background: #f7f8fa;
  display: flex;
  flex-direction: column;
}

.tab-bar {
  display: flex;
  background: #fff;
  border-bottom: 1rpx solid #f0f0f0;
  flex-shrink: 0;
}

.tab-item {
  flex: 1;
  height: 88rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
}

.tab-item.active::after {
  content: '';
  position: absolute;
  bottom: 0;
  left: 30%;
  right: 30%;
  height: 4rpx;
  background: #0f3a57;
  border-radius: 2rpx;
}

.tab-text {
  font-size: 28rpx;
  color: #666;
}

.tab-item.active .tab-text {
  color: #0f3a57;
  font-weight: 600;
}

.coupon-scroll {
  flex: 1;
  height: calc(100vh - 88rpx);
  padding: 20rpx 24rpx;
}

.empty {
  display: flex;
  justify-content: center;
  padding-top: 200rpx;
}

.empty-text {
  font-size: 28rpx;
  color: #999;
}

.coupon-card {
  display: flex;
  background: #fff;
  border-radius: 16rpx;
  overflow: hidden;
  margin-bottom: 20rpx;
}

.coupon-left {
  width: 220rpx;
  background: linear-gradient(135deg, #0f3a57, #1a5f8a);
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 30rpx 16rpx;
  flex-shrink: 0;
}

.coupon-amount-row {
  display: flex;
  align-items: baseline;
}

.currency {
  font-size: 28rpx;
  color: #fff;
  font-weight: 600;
}

.coupon-amount {
  font-size: 56rpx;
  font-weight: 700;
  color: #fff;
}

.coupon-condition {
  font-size: 22rpx;
  color: rgba(255, 255, 255, 0.75);
  margin-top: 8rpx;
}

.coupon-right {
  flex: 1;
  padding: 24rpx 28rpx;
  display: flex;
  flex-direction: column;
  justify-content: center;
}

.coupon-name {
  font-size: 28rpx;
  color: #333;
  font-weight: 500;
  margin-bottom: 8rpx;
  display: block;
}

.coupon-period {
  font-size: 22rpx;
  color: #999;
  margin-bottom: 16rpx;
  display: block;
}

.claim-btn {
  align-self: flex-end;
  padding: 10rpx 32rpx;
  background: #0f3a57;
  border-radius: 30rpx;
}

.claim-btn.claimed {
  background: #ccc;
}

.claim-btn-text {
  font-size: 24rpx;
  color: #fff;
}

.loading-more {
  display: flex;
  justify-content: center;
  padding: 30rpx 0;
}

.loading-text {
  font-size: 24rpx;
  color: #999;
}
</style>
