<template>
  <view class="history-page">
    <!-- Header with clear button -->
    <view class="page-header" v-if="groups.length > 0">
      <view class="clear-btn" @click="clearAll">
        <text class="clear-text">{{ t('common.delete') }} All</text>
      </view>
    </view>

    <!-- Loading -->
    <view v-if="loading" class="loading-wrap">
      <text class="loading-text">{{ t('common.loading') }}</text>
    </view>

    <!-- Error -->
    <view v-else-if="error" class="error-wrap">
      <text class="error-text">{{ error }}</text>
      <view class="retry-btn" @click="loadHistory">
        <text class="retry-text">{{ t('common.retry') }}</text>
      </view>
    </view>

    <!-- Empty -->
    <view v-else-if="groups.length === 0" class="empty-wrap">
      <text class="empty-text">{{ t('common.noData') }}</text>
    </view>

    <!-- History grouped by date -->
    <view v-else class="history-groups">
      <view v-for="group in groups" :key="group.date" class="history-group">
        <view class="group-header">
          <text class="group-date">{{ group.date }}</text>
        </view>
        <view class="group-items">
          <view
            v-for="item in group.items"
            :key="item.id"
            class="history-item"
            @click="goDetail(item.goods_id)"
          >
            <image class="item-image" :src="item.image" mode="aspectFill" />
            <view class="item-info">
              <text class="item-name">{{ item.name }}</text>
              <text class="item-price">{{ item.price }}</text>
            </view>
          </view>
        </view>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { t } from '@/locale'
import { commonApi } from '@/api/common'

interface HistoryItem {
  id: string
  goods_id: string
  name: string
  image: string
  price: string
  date: string
}

interface HistoryGroup {
  date: string
  items: HistoryItem[]
}

const loading = ref(true)
const error = ref('')
const groups = ref<HistoryGroup[]>([])

async function loadHistory() {
  loading.value = true
  error.value = ''
  try {
    const res: any = await commonApi.getHistory()
    const data = res?.data || res
    const items: HistoryItem[] = Array.isArray(data) ? data : data?.list || []

    // Group by date
    const map = new Map<string, HistoryItem[]>()
    for (const item of items) {
      const date = item.date || 'Unknown'
      if (!map.has(date)) {
        map.set(date, [])
      }
      map.get(date)!.push(item)
    }
    groups.value = Array.from(map.entries()).map(([date, items]) => ({
      date,
      items,
    }))
  } catch (e: any) {
    error.value = e?.message || t('common.retry')
  } finally {
    loading.value = false
  }
}

function clearAll() {
  uni.showModal({
    title: t('common.confirm'),
    content: 'Clear all browsing history?',
    success: async (res) => {
      if (res.confirm) {
        try {
          await commonApi.clearHistory()
          groups.value = []
          uni.showToast({ title: t('common.done'), icon: 'success' })
        } catch (e: any) {
          uni.showToast({ title: e?.message || t('common.retry'), icon: 'none' })
        }
      }
    },
  })
}

function goDetail(goodsId: string) {
  uni.navigateTo({ url: `/pages/goods/detail?id=${goodsId}` })
}

onMounted(() => {
  uni.setNavigationBarTitle({ title: t('user.history') })
  loadHistory()
})
</script>

<style scoped>
.history-page {
  min-height: 100vh;
  background-color: #f5f5f5;
}

.page-header {
  display: flex;
  justify-content: flex-end;
  padding: 16rpx 24rpx;
  background-color: #fff;
  border-bottom: 1rpx solid #eee;
}

.clear-btn {
  padding: 8rpx 16rpx;
}

.clear-text {
  font-size: 24rpx;
  color: #e64340;
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

.history-groups {
  padding: 20rpx 24rpx;
}

.history-group {
  margin-bottom: 24rpx;
}

.group-header {
  padding: 12rpx 0;
}

.group-date {
  font-size: 26rpx;
  color: #999;
  font-weight: 500;
}

.group-items {
  background-color: #fff;
  border-radius: 12rpx;
  overflow: hidden;
}

.history-item {
  display: flex;
  padding: 20rpx;
  border-bottom: 1rpx solid #f0f0f0;
}

.history-item:last-child {
  border-bottom: none;
}

.item-image {
  width: 140rpx;
  height: 140rpx;
  border-radius: 8rpx;
  margin-right: 20rpx;
  flex-shrink: 0;
}

.item-info {
  flex: 1;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
}

.item-name {
  font-size: 28rpx;
  color: #333;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.item-price {
  font-size: 30rpx;
  font-weight: 600;
  color: #e64340;
}
</style>
