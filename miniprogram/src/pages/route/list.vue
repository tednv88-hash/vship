<template>
  <view class="route-list">
    <!-- Search bar -->
    <view class="search-bar">
      <uni-icons type="search" size="18" color="#999" />
      <input
        class="search-input"
        :placeholder="'搜索線路'"
        v-model="searchKeyword"
        @input="onSearch"
      />
    </view>

    <!-- Filter tags -->
    <scroll-view scroll-x class="filter-scroll">
      <view class="filter-tags">
        <view
          v-for="tag in filterTags"
          :key="tag"
          class="filter-tag"
          :class="{ active: activeFilter === tag }"
          @click="onFilterChange(tag)"
        >
          <text>{{ tag }}</text>
        </view>
      </view>
    </scroll-view>

    <!-- Route cards -->
    <view class="routes-list">
      <view
        v-for="route in filteredRoutes"
        :key="route.id"
        class="route-card"
        @click="goDetail(route.id)"
      >
        <view class="route-header">
          <text class="route-name">{{ route.name }}</text>
          <uni-icons type="right" size="16" color="#999" />
        </view>

        <view class="route-info-grid">
          <view class="info-item">
            <text class="info-label">時效</text>
            <text class="info-value">{{ route.transit_time }}</text>
          </view>
          <view class="info-item">
            <text class="info-label">價格</text>
            <text class="info-value price">{{ route.price_range }}</text>
          </view>
        </view>

        <view class="route-tags">
          <view v-for="(type, idx) in route.item_types" :key="idx" class="item-tag">
            <text>{{ type }}</text>
          </view>
        </view>
      </view>
    </view>

    <!-- Empty -->
    <view v-if="!loading && filteredRoutes.length === 0" class="empty">
      <text>{{ t('common.noData') }}</text>
    </view>
    <view v-if="loading" class="load-more">
      <text>{{ t('common.loading') }}</text>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { t } from '@/locale'
import { commonApi } from '@/api/common'

const routeList = ref<any[]>([])
const searchKeyword = ref('')
const activeFilter = ref('全部')
const loading = ref(true)

const filterTags = ['全部', '日本', '韓國', '美國', '澳洲', '歐洲']

const filteredRoutes = computed(() => {
  let list = routeList.value
  if (activeFilter.value !== '全部') {
    list = list.filter((r) => r.name?.includes(activeFilter.value) || r.origin?.includes(activeFilter.value))
  }
  if (searchKeyword.value.trim()) {
    const kw = searchKeyword.value.trim().toLowerCase()
    list = list.filter(
      (r) =>
        r.name?.toLowerCase().includes(kw) ||
        r.origin?.toLowerCase().includes(kw) ||
        r.destination?.toLowerCase().includes(kw)
    )
  }
  return list
})

onMounted(() => {
  loadRoutes()
})

async function loadRoutes() {
  loading.value = true
  try {
    const res = await commonApi.getRoutes()
    routeList.value = res?.data?.list || res?.data || []
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

function onFilterChange(tag: string) {
  activeFilter.value = tag
}

function onSearch() {
  // Reactive filtering via computed
}

function goDetail(id: string) {
  uni.navigateTo({ url: `/pages/route/detail?id=${id}` })
}
</script>

<style scoped>
.route-list {
  min-height: 100vh;
  background: #f5f5f5;
}

.search-bar {
  display: flex;
  align-items: center;
  margin: 20rpx 24rpx;
  padding: 16rpx 24rpx;
  background: #fff;
  border-radius: 32rpx;
  gap: 12rpx;
}

.search-input {
  flex: 1;
  font-size: 28rpx;
  color: #333;
}

.filter-scroll {
  white-space: nowrap;
  padding: 0 24rpx;
  margin-bottom: 16rpx;
}

.filter-tags {
  display: inline-flex;
  gap: 16rpx;
}

.filter-tag {
  padding: 12rpx 28rpx;
  background: #fff;
  border-radius: 24rpx;
  display: inline-flex;
}

.filter-tag text {
  font-size: 26rpx;
  color: #666;
}

.filter-tag.active {
  background: #0f3a57;
}

.filter-tag.active text {
  color: #fff;
}

.routes-list {
  padding: 0 24rpx;
}

.route-card {
  background: #fff;
  border-radius: 16rpx;
  padding: 28rpx 30rpx;
  margin-bottom: 16rpx;
}

.route-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20rpx;
}

.route-name {
  font-size: 32rpx;
  color: #333;
  font-weight: 600;
}

.route-info-grid {
  display: flex;
  gap: 40rpx;
  margin-bottom: 16rpx;
}

.info-item {
  display: flex;
  flex-direction: column;
  gap: 6rpx;
}

.info-label {
  font-size: 24rpx;
  color: #999;
}

.info-value {
  font-size: 28rpx;
  color: #333;
  font-weight: 500;
}

.info-value.price {
  color: #e64340;
}

.route-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 12rpx;
}

.item-tag {
  padding: 6rpx 16rpx;
  background: rgba(15, 58, 87, 0.06);
  border-radius: 6rpx;
}

.item-tag text {
  font-size: 22rpx;
  color: #0f3a57;
}

.empty {
  text-align: center;
  padding: 100rpx 0;
  color: #999;
  font-size: 28rpx;
}

.load-more {
  text-align: center;
  padding: 30rpx 0;
  font-size: 24rpx;
  color: #999;
}
</style>
