<template>
  <view class="warehouse-page">
    <view v-if="list.length === 0 && !loading" class="empty">
      <text class="empty-text">{{ t('common.noData') }}</text>
    </view>

    <view v-for="item in list" :key="item.id" class="warehouse-card">
      <!-- Header -->
      <view class="card-header">
        <view class="warehouse-flag">
          <text class="flag-text">{{ item.country_flag || '🏭' }}</text>
        </view>
        <text class="warehouse-name">{{ item.name }}</text>
      </view>

      <!-- Address fields -->
      <view class="field-list">
        <view class="field-item">
          <text class="field-label">收件人</text>
          <view class="field-value-row">
            <text class="field-value">{{ item.recipient }}</text>
            <view class="copy-btn" @tap="handleCopy(item.recipient)">
              <text class="copy-text">{{ t('common.copy') }}</text>
            </view>
          </view>
        </view>
        <view class="field-item">
          <text class="field-label">聯繫電話</text>
          <view class="field-value-row">
            <text class="field-value">{{ item.phone }}</text>
            <view class="copy-btn" @tap="handleCopy(item.phone)">
              <text class="copy-text">{{ t('common.copy') }}</text>
            </view>
          </view>
        </view>
        <view class="field-item">
          <text class="field-label">郵編</text>
          <view class="field-value-row">
            <text class="field-value">{{ item.zip_code }}</text>
            <view class="copy-btn" @tap="handleCopy(item.zip_code)">
              <text class="copy-text">{{ t('common.copy') }}</text>
            </view>
          </view>
        </view>
        <view class="field-item">
          <text class="field-label">地址</text>
          <view class="field-value-row">
            <text class="field-value address-text">{{ item.full_address }}</text>
            <view class="copy-btn" @tap="handleCopy(item.full_address)">
              <text class="copy-text">{{ t('common.copy') }}</text>
            </view>
          </view>
        </view>
      </view>

      <!-- Copy all -->
      <view class="copy-all-btn" @tap="handleCopyAll(item)">
        <text class="copy-all-text">一鍵複製全部</text>
      </view>
    </view>

    <!-- Tips -->
    <view class="tips-section" v-if="list.length > 0">
      <text class="tips-title">溫馨提示</text>
      <text class="tips-content">
        請將包裹寄送至上方倉庫地址，收件人欄位請務必填寫您的會員ID以便入庫識別。
      </text>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { t } from '@/locale'
import { commonApi } from '@/api/common'

interface WarehouseItem {
  id: string
  name: string
  country_flag: string
  recipient: string
  phone: string
  zip_code: string
  full_address: string
}

const list = ref<WarehouseItem[]>([])
const loading = ref(false)

async function loadList() {
  loading.value = true
  try {
    const res = (await commonApi.getWarehouses()) as any
    const data = res?.data || res || []
    list.value = Array.isArray(data) ? data : data?.list || []
  } catch {
    uni.showToast({ title: '載入失敗', icon: 'none' })
  } finally {
    loading.value = false
  }
}

function handleCopy(text: string) {
  uni.setClipboardData({
    data: text,
    success: () => {
      uni.showToast({ title: t('common.copied'), icon: 'success' })
    },
  })
}

function handleCopyAll(item: WarehouseItem) {
  const text = `收件人: ${item.recipient}\n電話: ${item.phone}\n郵編: ${item.zip_code}\n地址: ${item.full_address}`
  uni.setClipboardData({
    data: text,
    success: () => {
      uni.showToast({ title: t('common.copied'), icon: 'success' })
    },
  })
}

onMounted(() => {
  loadList()
})
</script>

<style scoped>
.warehouse-page {
  min-height: 100vh;
  background: #f7f8fa;
  padding: 24rpx;
}

.empty {
  display: flex;
  justify-content: center;
  padding-top: 300rpx;
}

.empty-text {
  font-size: 28rpx;
  color: #999;
}

.warehouse-card {
  background: #fff;
  border-radius: 16rpx;
  margin-bottom: 24rpx;
  overflow: hidden;
}

.card-header {
  display: flex;
  align-items: center;
  padding: 28rpx 30rpx;
  background: linear-gradient(135deg, #0f3a57, #1a5f8a);
}

.warehouse-flag {
  width: 56rpx;
  height: 56rpx;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.2);
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: 20rpx;
}

.flag-text {
  font-size: 32rpx;
}

.warehouse-name {
  font-size: 32rpx;
  font-weight: 600;
  color: #fff;
}

.field-list {
  padding: 16rpx 30rpx;
}

.field-item {
  padding: 20rpx 0;
  border-bottom: 1rpx solid #f5f5f5;
}

.field-item:last-child {
  border-bottom: none;
}

.field-label {
  display: block;
  font-size: 24rpx;
  color: #999;
  margin-bottom: 8rpx;
}

.field-value-row {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
}

.field-value {
  flex: 1;
  font-size: 28rpx;
  color: #333;
  margin-right: 16rpx;
  word-break: break-all;
}

.address-text {
  line-height: 1.5;
}

.copy-btn {
  flex-shrink: 0;
  padding: 8rpx 20rpx;
  border: 2rpx solid #0f3a57;
  border-radius: 8rpx;
}

.copy-text {
  font-size: 22rpx;
  color: #0f3a57;
}

.copy-all-btn {
  margin: 0 30rpx 24rpx;
  height: 80rpx;
  background: rgba(15, 58, 87, 0.06);
  border-radius: 12rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  border: 2rpx solid #0f3a57;
}

.copy-all-text {
  font-size: 28rpx;
  color: #0f3a57;
  font-weight: 500;
}

.tips-section {
  background: #fff;
  border-radius: 16rpx;
  padding: 28rpx 30rpx;
  margin-top: 16rpx;
}

.tips-title {
  display: block;
  font-size: 28rpx;
  color: #333;
  font-weight: 600;
  margin-bottom: 12rpx;
}

.tips-content {
  font-size: 24rpx;
  color: #999;
  line-height: 1.6;
}
</style>
