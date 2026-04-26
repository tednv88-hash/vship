<template>
  <view class="refund-list-page">
    <!-- Status tabs -->
    <view class="status-tabs">
      <view
        v-for="tab in tabs"
        :key="tab.value"
        class="tab-item"
        :class="{ active: activeTab === tab.value }"
        @click="switchTab(tab.value)"
      >
        <text class="tab-text">{{ tab.label }}</text>
        <view v-if="activeTab === tab.value" class="tab-line" />
      </view>
    </view>

    <!-- Loading -->
    <view v-if="loading && list.length === 0" class="loading-wrap">
      <text class="loading-text">{{ t('common.loading') }}</text>
    </view>

    <!-- Error -->
    <view v-else-if="error && list.length === 0" class="error-wrap">
      <text class="error-text">{{ error }}</text>
      <view class="retry-btn" @click="loadRefunds(true)">
        <text class="retry-text">{{ t('common.retry') }}</text>
      </view>
    </view>

    <!-- List -->
    <view v-else class="refund-items">
      <view
        v-for="item in list"
        :key="item.id"
        class="refund-item"
        @click="goDetail(item.id)"
      >
        <view class="refund-header">
          <text class="order-no">Order: {{ item.order_no }}</text>
          <text class="refund-status" :class="'status-' + item.status">
            {{ getStatusText(item.status) }}
          </text>
        </view>
        <view class="refund-body">
          <text class="refund-amount">{{ item.amount }}</text>
          <text class="refund-date">{{ item.date }}</text>
        </view>
      </view>

      <!-- Empty -->
      <view v-if="!loading && list.length === 0" class="empty-wrap">
        <text class="empty-text">{{ t('common.noData') }}</text>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { t } from '@/locale'
import { commonApi } from '@/api/common'

interface RefundItem {
  id: string
  order_no: string
  status: string
  amount: string
  date: string
}

const tabs = [
  { label: '全部', value: '' },
  { label: '處理中', value: 'processing' },
  { label: '已完成', value: 'completed' },
  { label: '已拒絕', value: 'rejected' },
]

const activeTab = ref('')
const loading = ref(true)
const error = ref('')
const list = ref<RefundItem[]>([])

function getStatusText(status: string): string {
  const map: Record<string, string> = {
    pending: '待處理',
    processing: '處理中',
    completed: '已退款',
    rejected: '已拒絕',
  }
  return map[status] || status
}

async function loadRefunds(reset = false) {
  if (reset) {
    list.value = []
    loading.value = true
  }
  error.value = ''
  try {
    const params: any = {}
    if (activeTab.value) params.status = activeTab.value
    const res: any = await commonApi.getRefunds(params)
    const data = res?.data || res
    list.value = Array.isArray(data) ? data : data?.list || []
  } catch (e: any) {
    error.value = e?.message || t('common.retry')
  } finally {
    loading.value = false
  }
}

function switchTab(value: string) {
  activeTab.value = value
  loadRefunds(true)
}

function goDetail(id: string) {
  uni.navigateTo({ url: `/pages/refund/detail?id=${id}` })
}

onMounted(() => {
  uni.setNavigationBarTitle({ title: t('refund.title') })
  loadRefunds(true)
})
</script>

<style scoped>
.refund-list-page {
  min-height: 100vh;
  background-color: #f5f5f5;
}

.status-tabs {
  display: flex;
  background-color: #fff;
  border-bottom: 1rpx solid #eee;
}

.tab-item {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 24rpx 0;
  position: relative;
}

.tab-text {
  font-size: 28rpx;
  color: #666;
}

.tab-item.active .tab-text {
  color: #0f3a57;
  font-weight: 600;
}

.tab-line {
  position: absolute;
  bottom: 0;
  width: 48rpx;
  height: 4rpx;
  background-color: #0f3a57;
  border-radius: 2rpx;
}

.loading-wrap,
.error-wrap,
.empty-wrap {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 120rpx 0;
}

.loading-text,
.empty-text {
  font-size: 28rpx;
  color: #999;
}

.error-text {
  font-size: 28rpx;
  color: #e64340;
  margin-bottom: 24rpx;
}

.retry-btn {
  padding: 16rpx 48rpx;
  background-color: #0f3a57;
  border-radius: 8rpx;
}

.retry-text {
  font-size: 28rpx;
  color: #fff;
}

.refund-items {
  padding: 20rpx 24rpx;
}

.refund-item {
  background-color: #fff;
  border-radius: 12rpx;
  padding: 24rpx;
  margin-bottom: 16rpx;
}

.refund-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16rpx;
}

.order-no {
  font-size: 26rpx;
  color: #333;
}

.refund-status {
  font-size: 24rpx;
  font-weight: 500;
}

.status-pending,
.status-processing {
  color: #f5a623;
}

.status-completed {
  color: #07c160;
}

.status-rejected {
  color: #e64340;
}

.refund-body {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.refund-amount {
  font-size: 32rpx;
  font-weight: 600;
  color: #e64340;
}

.refund-date {
  font-size: 24rpx;
  color: #999;
}
</style>
