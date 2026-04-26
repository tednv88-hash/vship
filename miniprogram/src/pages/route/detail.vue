<template>
  <view class="route-detail">
    <!-- Header -->
    <view class="route-header">
      <text class="route-name">{{ detail.name }}</text>
      <text class="route-time">預計時效: {{ detail.transit_time }}</text>
    </view>

    <!-- Pricing table -->
    <view class="section-card">
      <text class="section-title">運費價格表</text>
      <view class="price-table">
        <view class="table-header">
          <text class="table-cell">重量</text>
          <text class="table-cell">首重價格</text>
          <text class="table-cell">續重價格</text>
        </view>
        <view v-for="(row, idx) in detail.pricing || []" :key="idx" class="table-row">
          <text class="table-cell">{{ row.weight_range }}</text>
          <text class="table-cell">¥{{ row.first_price }}/kg</text>
          <text class="table-cell">¥{{ row.extra_price }}/kg</text>
        </view>
      </view>
      <view v-if="!detail.pricing || detail.pricing.length === 0" class="empty-text">
        <text>暫無價格數據</text>
      </view>
    </view>

    <!-- Transit time -->
    <view class="section-card">
      <text class="section-title">時效說明</text>
      <view class="info-rows">
        <view class="info-row">
          <text class="info-label">預計時效</text>
          <text class="info-value">{{ detail.transit_time || '-' }}</text>
        </view>
        <view class="info-row">
          <text class="info-label">發貨週期</text>
          <text class="info-value">{{ detail.shipping_cycle || '每週發貨' }}</text>
        </view>
        <view class="info-row">
          <text class="info-label">起寄重量</text>
          <text class="info-value">{{ detail.min_weight || '0.5' }} kg</text>
        </view>
        <view class="info-row">
          <text class="info-label">體積重算法</text>
          <text class="info-value">{{ detail.volume_formula || '長×寬×高/6000' }}</text>
        </view>
      </view>
    </view>

    <!-- Restrictions -->
    <view class="section-card">
      <text class="section-title">限制說明</text>
      <view v-if="detail.restrictions && detail.restrictions.length > 0" class="restrictions-list">
        <view v-for="(item, idx) in detail.restrictions" :key="idx" class="restriction-item">
          <uni-icons type="info" size="14" color="#ff9800" />
          <text>{{ item }}</text>
        </view>
      </view>
      <view v-else class="restrictions-list">
        <view class="restriction-item">
          <uni-icons type="info" size="14" color="#ff9800" />
          <text>禁止寄送違禁品、易燃易爆品</text>
        </view>
        <view class="restriction-item">
          <uni-icons type="info" size="14" color="#ff9800" />
          <text>單件包裹重量不超過30kg</text>
        </view>
        <view class="restriction-item">
          <uni-icons type="info" size="14" color="#ff9800" />
          <text>單邊長度不超過150cm</text>
        </view>
      </view>
    </view>

    <!-- Warehouse address -->
    <view class="section-card">
      <text class="section-title">倉庫地址</text>
      <view v-if="detail.warehouse" class="warehouse-info">
        <view class="info-row">
          <text class="info-label">收件人</text>
          <text class="info-value">{{ detail.warehouse.recipient }}</text>
        </view>
        <view class="info-row">
          <text class="info-label">電話</text>
          <text class="info-value">{{ detail.warehouse.phone }}</text>
        </view>
        <view class="info-row">
          <text class="info-label">地址</text>
          <text class="info-value address">{{ detail.warehouse.address }}</text>
        </view>
        <view class="info-row">
          <text class="info-label">郵編</text>
          <text class="info-value">{{ detail.warehouse.zipcode }}</text>
        </view>
        <view class="copy-btn" @click="copyAddress">
          <text>{{ t('common.copy') }}地址</text>
        </view>
      </view>
      <view v-else class="empty-text">
        <text>暫無倉庫地址</text>
      </view>
    </view>

    <!-- Bottom bar: go to estimate -->
    <view class="bottom-bar">
      <view class="estimate-btn" @click="goEstimate">
        <text>費用估算</text>
      </view>
    </view>
    <view style="height: 120rpx" />
  </view>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { t } from '@/locale'
import { commonApi } from '@/api/common'

const detail = ref<any>({})

onMounted(() => {
  const pages = getCurrentPages()
  const currentPage = pages[pages.length - 1] as any
  const options = currentPage?.$page?.options || {}
  if (options.id) {
    loadDetail(options.id)
  }
})

async function loadDetail(id: string) {
  uni.showLoading({ title: t('common.loading') })
  try {
    const res = await commonApi.getRouteDetail(id)
    detail.value = res?.data || res || {}
    uni.setNavigationBarTitle({ title: detail.value.name || t('route.detail') })
  } catch (e) {
    console.error(e)
  } finally {
    uni.hideLoading()
  }
}

function copyAddress() {
  const wh = detail.value.warehouse
  if (!wh) return
  const text = `${wh.recipient}\n${wh.phone}\n${wh.address}\n${wh.zipcode}`
  uni.setClipboardData({
    data: text,
    success() {
      uni.showToast({ title: t('common.copied'), icon: 'success' })
    },
  })
}

function goEstimate() {
  uni.navigateTo({ url: '/pages/estimate/index' })
}
</script>

<style scoped>
.route-detail {
  min-height: 100vh;
  background: #f5f5f5;
}

.route-header {
  background: linear-gradient(135deg, #0f3a57 0%, #1a5c7a 100%);
  padding: 40rpx 30rpx;
  display: flex;
  flex-direction: column;
  gap: 12rpx;
}

.route-name {
  font-size: 36rpx;
  color: #fff;
  font-weight: 700;
}

.route-time {
  font-size: 26rpx;
  color: rgba(255, 255, 255, 0.7);
}

.section-card {
  background: #fff;
  margin: 16rpx 24rpx 0;
  border-radius: 16rpx;
  padding: 28rpx 30rpx;
}

.section-title {
  font-size: 30rpx;
  color: #333;
  font-weight: 600;
  margin-bottom: 20rpx;
  display: block;
}

.price-table {
  border: 1rpx solid #eee;
  border-radius: 8rpx;
  overflow: hidden;
}

.table-header {
  display: flex;
  background: rgba(15, 58, 87, 0.06);
}

.table-header .table-cell {
  font-weight: 600;
  color: #333;
}

.table-row {
  display: flex;
  border-top: 1rpx solid #eee;
}

.table-cell {
  flex: 1;
  padding: 18rpx 16rpx;
  font-size: 24rpx;
  color: #666;
  text-align: center;
}

.info-rows {
  display: flex;
  flex-direction: column;
}

.info-row {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  padding: 16rpx 0;
  border-bottom: 1rpx solid #f5f5f5;
}

.info-row:last-child {
  border-bottom: none;
}

.info-label {
  font-size: 26rpx;
  color: #999;
  flex-shrink: 0;
  width: 160rpx;
}

.info-value {
  font-size: 26rpx;
  color: #333;
  text-align: right;
  flex: 1;
}

.info-value.address {
  line-height: 1.5;
}

.copy-btn {
  margin-top: 20rpx;
  border: 2rpx solid #0f3a57;
  border-radius: 24rpx;
  padding: 14rpx 0;
  text-align: center;
}

.copy-btn text {
  font-size: 26rpx;
  color: #0f3a57;
}

.restrictions-list {
  display: flex;
  flex-direction: column;
  gap: 16rpx;
}

.restriction-item {
  display: flex;
  align-items: flex-start;
  gap: 10rpx;
}

.restriction-item text {
  font-size: 26rpx;
  color: #666;
  line-height: 1.5;
}

.warehouse-info {
  margin-top: 4rpx;
}

.empty-text {
  text-align: center;
  padding: 30rpx 0;
  color: #999;
  font-size: 26rpx;
}

.bottom-bar {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  background: #fff;
  padding: 16rpx 30rpx;
  padding-bottom: calc(16rpx + env(safe-area-inset-bottom));
  box-shadow: 0 -2rpx 10rpx rgba(0, 0, 0, 0.05);
  z-index: 100;
}

.estimate-btn {
  background: #0f3a57;
  border-radius: 44rpx;
  height: 84rpx;
  display: flex;
  align-items: center;
  justify-content: center;
}

.estimate-btn text {
  color: #fff;
  font-size: 30rpx;
  font-weight: 500;
}
</style>
