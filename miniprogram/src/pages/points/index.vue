<template>
  <view class="points-page">
    <!-- Points summary -->
    <view class="points-card">
      <text class="points-label">我的積分</text>
      <text class="points-total">{{ totalPoints }}</text>
    </view>

    <!-- List -->
    <scroll-view
      scroll-y
      class="list-scroll"
      :refresher-enabled="true"
      :refresher-triggered="refreshing"
      @refresherrefresh="onRefresh"
      @scrolltolower="loadMore"
    >
      <view v-if="list.length === 0 && !loading" class="empty">
        <text class="empty-text">{{ t('common.noData') }}</text>
      </view>

      <view v-for="item in list" :key="item.id" class="points-item">
        <view class="item-left">
          <text class="item-desc">{{ item.description }}</text>
          <text class="item-date">{{ item.created_at }}</text>
        </view>
        <text
          class="item-points"
          :class="{ positive: item.points > 0, negative: item.points < 0 }"
        >
          {{ item.points > 0 ? '+' : '' }}{{ item.points }}
        </text>
      </view>

      <view v-if="loading" class="loading-more">
        <text class="loading-text">{{ t('common.loading') }}</text>
      </view>
      <view v-else-if="finished && list.length > 0" class="loading-more">
        <text class="loading-text">— 沒有更多了 —</text>
      </view>
    </scroll-view>
  </view>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { t } from '@/locale'
import { userApi } from '@/api/user'

interface PointsItem {
  id: string
  description: string
  points: number
  created_at: string
}

const totalPoints = ref(0)
const list = ref<PointsItem[]>([])
const loading = ref(false)
const refreshing = ref(false)
const finished = ref(false)
const page = ref(1)
const pageSize = 20

async function loadUserInfo() {
  try {
    const res = (await userApi.getUserInfo()) as any
    const data = res?.data || res
    totalPoints.value = data?.points || 0
  } catch {
    // ignore
  }
}

async function loadList(reset = false) {
  if (loading.value) return
  if (reset) {
    page.value = 1
    finished.value = false
    list.value = []
  }
  loading.value = true

  try {
    const params = { page: page.value, page_size: pageSize }
    const res = (await userApi.getPointsList(params)) as any
    const data = res?.data || res || []
    const items = Array.isArray(data) ? data : data?.list || []

    if (reset) {
      list.value = items
    } else {
      list.value = [...list.value, ...items]
    }

    if (items.length < pageSize) {
      finished.value = true
    } else {
      page.value++
    }
  } catch {
    uni.showToast({ title: '載入失敗', icon: 'none' })
  } finally {
    loading.value = false
    refreshing.value = false
  }
}

function onRefresh() {
  refreshing.value = true
  loadUserInfo()
  loadList(true)
}

function loadMore() {
  if (!finished.value && !loading.value) {
    loadList()
  }
}

onMounted(() => {
  loadUserInfo()
  loadList(true)
})
</script>

<style scoped>
.points-page {
  min-height: 100vh;
  background: #f7f8fa;
  display: flex;
  flex-direction: column;
}

.points-card {
  background: linear-gradient(135deg, #0f3a57, #1a5f8a);
  margin: 24rpx;
  border-radius: 20rpx;
  padding: 50rpx 40rpx;
  flex-shrink: 0;
}

.points-label {
  font-size: 26rpx;
  color: rgba(255, 255, 255, 0.7);
  display: block;
  margin-bottom: 16rpx;
}

.points-total {
  font-size: 64rpx;
  font-weight: 700;
  color: #fff;
}

.list-scroll {
  flex: 1;
  height: 0;
}

.empty {
  display: flex;
  justify-content: center;
  padding-top: 150rpx;
}

.empty-text {
  font-size: 28rpx;
  color: #999;
}

.points-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  background: #fff;
  margin: 0 24rpx 16rpx;
  border-radius: 12rpx;
  padding: 28rpx 30rpx;
}

.item-left {
  flex: 1;
}

.item-desc {
  display: block;
  font-size: 28rpx;
  color: #333;
  margin-bottom: 8rpx;
}

.item-date {
  font-size: 24rpx;
  color: #999;
}

.item-points {
  font-size: 34rpx;
  font-weight: 600;
  flex-shrink: 0;
  margin-left: 20rpx;
}

.item-points.positive {
  color: #07c160;
}

.item-points.negative {
  color: #e74c3c;
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
