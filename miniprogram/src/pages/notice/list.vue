<template>
  <view class="notice-list-page">
    <!-- Loading -->
    <view v-if="loading && list.length === 0" class="loading-wrap">
      <text class="loading-text">{{ t('common.loading') }}</text>
    </view>

    <!-- Error -->
    <view v-else-if="error && list.length === 0" class="error-wrap">
      <text class="error-text">{{ error }}</text>
      <view class="retry-btn" @click="loadNotices(true)">
        <text class="retry-text">{{ t('common.retry') }}</text>
      </view>
    </view>

    <!-- List -->
    <scroll-view
      v-else
      class="notice-scroll"
      scroll-y
      refresher-enabled
      :refresher-triggered="refreshing"
      @refresherrefresh="onRefresh"
    >
      <view class="notice-items">
        <view
          v-for="item in list"
          :key="item.id"
          class="notice-item"
          @click="goDetail(item.id)"
        >
          <view class="notice-left">
            <view v-if="!item.is_read" class="unread-dot" />
            <text class="notice-title" :class="{ unread: !item.is_read }">
              {{ item.title }}
            </text>
          </view>
          <text class="notice-date">{{ item.date }}</text>
        </view>
      </view>

      <!-- Empty -->
      <view v-if="!loading && list.length === 0" class="empty-wrap">
        <text class="empty-text">{{ t('common.noData') }}</text>
      </view>
    </scroll-view>
  </view>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { t } from '@/locale'
import { commonApi } from '@/api/common'

interface Notice {
  id: string
  title: string
  date: string
  is_read: boolean
}

const loading = ref(true)
const refreshing = ref(false)
const error = ref('')
const list = ref<Notice[]>([])

async function loadNotices(reset = false) {
  if (reset) {
    list.value = []
    loading.value = true
  }
  error.value = ''
  try {
    const res: any = await commonApi.getNotices()
    const data = res?.data || res
    list.value = Array.isArray(data) ? data : data?.list || []
  } catch (e: any) {
    error.value = e?.message || t('common.retry')
  } finally {
    loading.value = false
    refreshing.value = false
  }
}

function onRefresh() {
  refreshing.value = true
  loadNotices(true)
}

function goDetail(id: string) {
  uni.navigateTo({ url: `/pages/notice/detail?id=${id}` })
}

onMounted(() => {
  uni.setNavigationBarTitle({ title: t('notice.list') })
  loadNotices(true)
})
</script>

<style scoped>
.notice-list-page {
  min-height: 100vh;
  background-color: #f5f5f5;
}

.notice-scroll {
  height: 100vh;
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

.notice-items {
  padding: 20rpx 24rpx;
}

.notice-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  background-color: #fff;
  border-radius: 12rpx;
  padding: 28rpx 24rpx;
  margin-bottom: 16rpx;
}

.notice-left {
  display: flex;
  align-items: center;
  flex: 1;
  margin-right: 16rpx;
}

.unread-dot {
  width: 16rpx;
  height: 16rpx;
  border-radius: 50%;
  background-color: #e64340;
  margin-right: 16rpx;
  flex-shrink: 0;
}

.notice-title {
  font-size: 28rpx;
  color: #666;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.notice-title.unread {
  color: #333;
  font-weight: 500;
}

.notice-date {
  font-size: 24rpx;
  color: #999;
  flex-shrink: 0;
}
</style>
