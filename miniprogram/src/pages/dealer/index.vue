<template>
  <view class="dealer-center">
    <!-- Header stats -->
    <view class="header-card">
      <text class="header-title">{{ t('dealer.center') }}</text>
      <view class="stats-grid">
        <view class="stat-item">
          <text class="stat-value">¥{{ dealerInfo.total_earnings || '0.00' }}</text>
          <text class="stat-label">累計收益</text>
        </view>
        <view class="stat-item">
          <text class="stat-value">¥{{ dealerInfo.available_balance || '0.00' }}</text>
          <text class="stat-label">可用餘額</text>
        </view>
        <view class="stat-item">
          <text class="stat-value">{{ dealerInfo.total_orders || 0 }}</text>
          <text class="stat-label">總訂單數</text>
        </view>
        <view class="stat-item">
          <text class="stat-value">{{ dealerInfo.team_count || 0 }}</text>
          <text class="stat-label">團隊成員</text>
        </view>
      </view>
    </view>

    <!-- Quick actions -->
    <view class="actions-card">
      <view class="action-item" @click="goTo('/pages/dealer/withdraw')">
        <view class="action-icon" style="background: #ff9800">
          <uni-icons type="wallet" size="24" color="#fff" />
        </view>
        <text class="action-text">提現</text>
      </view>
      <view class="action-item" @click="goTo('/pages/dealer/orders')">
        <view class="action-icon" style="background: #4caf50">
          <uni-icons type="shop" size="24" color="#fff" />
        </view>
        <text class="action-text">訂單</text>
      </view>
      <view class="action-item" @click="goTo('/pages/dealer/team')">
        <view class="action-icon" style="background: #2196f3">
          <uni-icons type="person" size="24" color="#fff" />
        </view>
        <text class="action-text">團隊</text>
      </view>
      <view class="action-item" @click="generatePoster">
        <view class="action-icon" style="background: #9c27b0">
          <uni-icons type="image" size="24" color="#fff" />
        </view>
        <text class="action-text">海報</text>
      </view>
    </view>

    <!-- Earnings overview -->
    <view class="section-card">
      <text class="section-title">收益概覽</text>
      <view class="earnings-list">
        <view class="earnings-row">
          <text class="earnings-label">今日收益</text>
          <text class="earnings-value">¥{{ dealerInfo.today_earnings || '0.00' }}</text>
        </view>
        <view class="earnings-row">
          <text class="earnings-label">本週收益</text>
          <text class="earnings-value">¥{{ dealerInfo.week_earnings || '0.00' }}</text>
        </view>
        <view class="earnings-row">
          <text class="earnings-label">本月收益</text>
          <text class="earnings-value">¥{{ dealerInfo.month_earnings || '0.00' }}</text>
        </view>
        <view class="earnings-row">
          <text class="earnings-label">已提現金額</text>
          <text class="earnings-value">¥{{ dealerInfo.withdrawn_amount || '0.00' }}</text>
        </view>
      </view>
    </view>

    <!-- Recent orders -->
    <view class="section-card">
      <view class="section-header">
        <text class="section-title">最近訂單</text>
        <text class="section-more" @click="goTo('/pages/dealer/orders')">{{ t('common.more') }} ></text>
      </view>
      <view v-if="recentOrders.length > 0" class="orders-list">
        <view v-for="order in recentOrders" :key="order.id" class="order-item">
          <view class="order-info">
            <text class="order-no">{{ order.order_no }}</text>
            <text class="order-date">{{ order.created_at }}</text>
          </view>
          <text class="order-commission">+¥{{ order.commission }}</text>
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

const dealerInfo = ref<any>({})
const recentOrders = ref<any[]>([])

onMounted(() => {
  loadDealerInfo()
  loadRecentOrders()
})

async function loadDealerInfo() {
  try {
    const res = await commonApi.getDealerInfo()
    dealerInfo.value = res?.data || res || {}
  } catch (e) {
    console.error(e)
  }
}

async function loadRecentOrders() {
  try {
    const res = await commonApi.getDealerOrders({ page: 1, limit: 5 })
    recentOrders.value = res?.data?.list || res?.data || []
  } catch (e) {
    console.error(e)
  }
}

function goTo(url: string) {
  uni.navigateTo({ url })
}

function generatePoster() {
  uni.navigateTo({ url: '/pages/invite/index' })
}
</script>

<style scoped>
.dealer-center {
  min-height: 100vh;
  background: #f5f5f5;
}

.header-card {
  background: linear-gradient(135deg, #0f3a57 0%, #1a5c7a 100%);
  padding: 40rpx 30rpx;
  color: #fff;
}

.header-title {
  font-size: 36rpx;
  font-weight: 600;
  display: block;
  margin-bottom: 30rpx;
}

.stats-grid {
  display: flex;
  flex-wrap: wrap;
}

.stat-item {
  width: 50%;
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 16rpx 0;
}

.stat-value {
  font-size: 36rpx;
  font-weight: 700;
  color: #fff;
}

.stat-label {
  font-size: 24rpx;
  color: rgba(255, 255, 255, 0.7);
  margin-top: 8rpx;
}

.actions-card {
  display: flex;
  background: #fff;
  padding: 30rpx 0;
  margin-bottom: 16rpx;
}

.action-item {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12rpx;
}

.action-icon {
  width: 80rpx;
  height: 80rpx;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
}

.action-text {
  font-size: 24rpx;
  color: #333;
}

.section-card {
  background: #fff;
  padding: 30rpx;
  margin-bottom: 16rpx;
}

.section-title {
  font-size: 30rpx;
  color: #333;
  font-weight: 600;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20rpx;
}

.section-more {
  font-size: 24rpx;
  color: #0f3a57;
}

.earnings-list {
  margin-top: 20rpx;
}

.earnings-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 18rpx 0;
  border-bottom: 1rpx solid #f5f5f5;
}

.earnings-row:last-child {
  border-bottom: none;
}

.earnings-label {
  font-size: 28rpx;
  color: #666;
}

.earnings-value {
  font-size: 28rpx;
  color: #333;
  font-weight: 500;
}

.orders-list {
  margin-top: 10rpx;
}

.order-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20rpx 0;
  border-bottom: 1rpx solid #f5f5f5;
}

.order-item:last-child {
  border-bottom: none;
}

.order-info {
  display: flex;
  flex-direction: column;
  gap: 8rpx;
}

.order-no {
  font-size: 26rpx;
  color: #333;
}

.order-date {
  font-size: 22rpx;
  color: #999;
}

.order-commission {
  font-size: 30rpx;
  color: #e64340;
  font-weight: 600;
}

.empty-text {
  text-align: center;
  padding: 40rpx 0;
  color: #999;
  font-size: 26rpx;
}
</style>
