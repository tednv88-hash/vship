<template>
  <view class="balance-page">
    <!-- Filter tabs -->
    <view class="filter-bar">
      <view
        v-for="tab in filterTabs"
        :key="tab.value"
        class="filter-tab"
        :class="{ active: currentFilter === tab.value }"
        @tap="switchFilter(tab.value)"
      >
        <text class="filter-text">{{ tab.label }}</text>
      </view>
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

      <view v-for="item in list" :key="item.id" class="balance-item">
        <view class="item-left">
          <text class="item-desc">{{ item.description }}</text>
          <text class="item-date">{{ item.created_at }}</text>
        </view>
        <view class="item-right">
          <text
            class="item-amount"
            :class="{ positive: item.amount > 0, negative: item.amount < 0 }"
          >
            {{ item.amount > 0 ? '+' : '' }}{{ item.amount }}
          </text>
          <view class="type-badge" :class="`type-${item.type}`">
            <text class="type-text">{{ getTypeLabel(item.type) }}</text>
          </view>
        </view>
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

interface BalanceItem {
  id: string
  description: string
  amount: number
  type: string
  created_at: string
}

const filterTabs = [
  { label: '全部', value: '' },
  { label: '充值', value: 'recharge' },
  { label: '消費', value: 'consume' },
  { label: '退款', value: 'refund' },
]

const currentFilter = ref('')
const list = ref<BalanceItem[]>([])
const loading = ref(false)
const refreshing = ref(false)
const finished = ref(false)
const page = ref(1)
const pageSize = 20

function getTypeLabel(type: string): string {
  const map: Record<string, string> = {
    recharge: '充值',
    consume: '消費',
    refund: '退款',
  }
  return map[type] || type
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
    const params: any = { page: page.value, page_size: pageSize }
    if (currentFilter.value) {
      params.type = currentFilter.value
    }
    const res = (await userApi.getBalanceList(params)) as any
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

function switchFilter(value: string) {
  currentFilter.value = value
  loadList(true)
}

function onRefresh() {
  refreshing.value = true
  loadList(true)
}

function loadMore() {
  if (!finished.value && !loading.value) {
    loadList()
  }
}

onMounted(() => {
  loadList(true)
})
</script>

<style scoped>
.balance-page {
  min-height: 100vh;
  background: #f7f8fa;
  display: flex;
  flex-direction: column;
}

.filter-bar {
  display: flex;
  background: #fff;
  padding: 0 16rpx;
  border-bottom: 1rpx solid #f0f0f0;
  flex-shrink: 0;
}

.filter-tab {
  flex: 1;
  height: 88rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
}

.filter-tab.active::after {
  content: '';
  position: absolute;
  bottom: 0;
  left: 30%;
  right: 30%;
  height: 4rpx;
  background: #0f3a57;
  border-radius: 2rpx;
}

.filter-text {
  font-size: 28rpx;
  color: #666;
}

.filter-tab.active .filter-text {
  color: #0f3a57;
  font-weight: 600;
}

.list-scroll {
  flex: 1;
  height: calc(100vh - 88rpx);
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

.balance-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  background: #fff;
  margin: 16rpx 24rpx 0;
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

.item-right {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
}

.item-amount {
  font-size: 32rpx;
  font-weight: 600;
  margin-bottom: 8rpx;
}

.item-amount.positive {
  color: #07c160;
}

.item-amount.negative {
  color: #e74c3c;
}

.type-badge {
  padding: 4rpx 14rpx;
  border-radius: 6rpx;
}

.type-recharge {
  background: rgba(15, 58, 87, 0.1);
}

.type-consume {
  background: rgba(231, 76, 60, 0.1);
}

.type-refund {
  background: rgba(7, 193, 96, 0.1);
}

.type-text {
  font-size: 20rpx;
  color: #666;
}

.type-recharge .type-text {
  color: #0f3a57;
}

.type-consume .type-text {
  color: #e74c3c;
}

.type-refund .type-text {
  color: #07c160;
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
