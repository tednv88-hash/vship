<template>
  <view class="refund-detail-page">
    <!-- Loading -->
    <view v-if="loading" class="loading-wrap">
      <text class="loading-text">{{ t('common.loading') }}</text>
    </view>

    <!-- Error -->
    <view v-else-if="error" class="error-wrap">
      <text class="error-text">{{ error }}</text>
      <view class="retry-btn" @click="loadDetail">
        <text class="retry-text">{{ t('common.retry') }}</text>
      </view>
    </view>

    <!-- Content -->
    <view v-else class="detail-content">
      <!-- Status header -->
      <view class="status-header" :class="'bg-' + detail.status">
        <text class="status-text">{{ getStatusText(detail.status) }}</text>
        <text class="status-desc">{{ detail.status_desc }}</text>
      </view>

      <!-- Order info -->
      <view class="info-card">
        <view class="card-title">
          <text class="card-title-text">Order Info</text>
        </view>
        <view class="info-row">
          <text class="info-label">Order No</text>
          <text class="info-value">{{ detail.order_no }}</text>
        </view>
        <view class="info-row">
          <text class="info-label">Refund Amount</text>
          <text class="info-value amount">{{ detail.amount }}</text>
        </view>
        <view class="info-row">
          <text class="info-label">Refund Reason</text>
          <text class="info-value">{{ detail.reason }}</text>
        </view>
        <view class="info-row" v-if="detail.remark">
          <text class="info-label">Remark</text>
          <text class="info-value">{{ detail.remark }}</text>
        </view>
      </view>

      <!-- Timeline -->
      <view class="timeline-card" v-if="detail.timeline && detail.timeline.length > 0">
        <view class="card-title">
          <text class="card-title-text">Progress</text>
        </view>
        <view class="timeline">
          <view
            v-for="(step, idx) in detail.timeline"
            :key="idx"
            class="timeline-item"
            :class="{ first: idx === 0 }"
          >
            <view class="timeline-dot" :class="{ active: idx === 0 }" />
            <view class="timeline-line" v-if="idx < detail.timeline.length - 1" />
            <view class="timeline-content">
              <text class="timeline-title">{{ step.title }}</text>
              <text class="timeline-time">{{ step.time }}</text>
            </view>
          </view>
        </view>
      </view>

      <!-- Actions -->
      <view class="action-bar" v-if="detail.actions && detail.actions.length > 0">
        <view
          v-for="action in detail.actions"
          :key="action.type"
          class="action-btn"
          :class="action.type === 'primary' ? 'primary' : 'secondary'"
          @click="onAction(action)"
        >
          <text class="action-text">{{ action.label }}</text>
        </view>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, onMounted, getCurrentInstance } from 'vue'
import { t } from '@/locale'
import { commonApi } from '@/api/common'

interface TimelineStep {
  title: string
  time: string
}

interface ActionItem {
  type: string
  label: string
  action: string
}

interface RefundDetail {
  id: string
  order_no: string
  status: string
  status_desc: string
  amount: string
  reason: string
  remark: string
  timeline: TimelineStep[]
  actions: ActionItem[]
}

const loading = ref(true)
const error = ref('')
const detail = ref<RefundDetail>({
  id: '',
  order_no: '',
  status: '',
  status_desc: '',
  amount: '',
  reason: '',
  remark: '',
  timeline: [],
  actions: [],
})

let refundId = ''

function getStatusText(status: string): string {
  const map: Record<string, string> = {
    pending: '待處理',
    processing: '處理中',
    completed: '已退款',
    rejected: '已拒絕',
  }
  return map[status] || status
}

async function loadDetail() {
  if (!refundId) return
  loading.value = true
  error.value = ''
  try {
    const res: any = await commonApi.getRefundDetail(refundId)
    const data = res?.data || res
    detail.value = {
      id: data.id || '',
      order_no: data.order_no || '',
      status: data.status || '',
      status_desc: data.status_desc || '',
      amount: data.amount || '',
      reason: data.reason || '',
      remark: data.remark || '',
      timeline: data.timeline || [],
      actions: data.actions || [],
    }
  } catch (e: any) {
    error.value = e?.message || t('common.retry')
  } finally {
    loading.value = false
  }
}

function onAction(action: ActionItem) {
  if (action.action === 'cancel') {
    uni.showModal({
      title: t('common.confirm'),
      content: 'Cancel this refund?',
      success: (res) => {
        if (res.confirm) {
          // Handle cancel refund
          uni.navigateBack()
        }
      },
    })
  }
}

onMounted(() => {
  uni.setNavigationBarTitle({ title: t('refund.detail') })
  const instance = getCurrentInstance()
  const options = (instance?.proxy as any)?.$page?.options || {}
  refundId = options.id || ''
  loadDetail()
})
</script>

<style scoped>
.refund-detail-page {
  min-height: 100vh;
  background-color: #f5f5f5;
}

.loading-wrap,
.error-wrap {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 120rpx 0;
}

.loading-text {
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

.status-header {
  padding: 40rpx 32rpx;
  background-color: #0f3a57;
}

.bg-processing {
  background-color: #f5a623;
}

.bg-completed {
  background-color: #07c160;
}

.bg-rejected {
  background-color: #e64340;
}

.status-text {
  font-size: 36rpx;
  font-weight: 600;
  color: #fff;
  display: block;
}

.status-desc {
  font-size: 26rpx;
  color: rgba(255, 255, 255, 0.8);
  margin-top: 8rpx;
  display: block;
}

.info-card,
.timeline-card {
  background-color: #fff;
  margin-top: 20rpx;
  padding: 24rpx 32rpx;
}

.card-title {
  margin-bottom: 16rpx;
  padding-bottom: 16rpx;
  border-bottom: 1rpx solid #f0f0f0;
}

.card-title-text {
  font-size: 30rpx;
  font-weight: 600;
  color: #333;
}

.info-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16rpx 0;
}

.info-label {
  font-size: 28rpx;
  color: #666;
}

.info-value {
  font-size: 28rpx;
  color: #333;
}

.info-value.amount {
  color: #e64340;
  font-weight: 600;
}

.timeline {
  padding: 8rpx 0;
}

.timeline-item {
  display: flex;
  position: relative;
  padding-bottom: 32rpx;
}

.timeline-item:last-child {
  padding-bottom: 0;
}

.timeline-dot {
  width: 20rpx;
  height: 20rpx;
  border-radius: 50%;
  background-color: #ddd;
  margin-right: 20rpx;
  margin-top: 6rpx;
  flex-shrink: 0;
  z-index: 1;
}

.timeline-dot.active {
  background-color: #0f3a57;
}

.timeline-line {
  position: absolute;
  left: 9rpx;
  top: 26rpx;
  width: 2rpx;
  bottom: 0;
  background-color: #eee;
}

.timeline-content {
  flex: 1;
}

.timeline-title {
  font-size: 28rpx;
  color: #333;
  display: block;
}

.timeline-item.first .timeline-title {
  color: #0f3a57;
  font-weight: 500;
}

.timeline-time {
  font-size: 24rpx;
  color: #999;
  margin-top: 4rpx;
  display: block;
}

.action-bar {
  display: flex;
  justify-content: flex-end;
  gap: 20rpx;
  padding: 32rpx;
  background-color: #fff;
  margin-top: 20rpx;
}

.action-btn {
  padding: 16rpx 40rpx;
  border-radius: 8rpx;
}

.action-btn.primary {
  background-color: #0f3a57;
}

.action-btn.secondary {
  background-color: #fff;
  border: 1rpx solid #ddd;
}

.action-btn.primary .action-text {
  color: #fff;
  font-size: 28rpx;
}

.action-btn.secondary .action-text {
  color: #333;
  font-size: 28rpx;
}
</style>
