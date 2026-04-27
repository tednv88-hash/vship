<template>
  <view class="goods-list">
    <!-- Search bar -->
    <view class="search-bar" @click="goSearch">
      <uni-icons type="search" size="18" color="#999" />
      <text class="search-placeholder">{{ t('common.search') }}</text>
    </view>

    <!-- Sort / Filter -->
    <view class="sort-bar">
      <view
        v-for="item in sortOptions"
        :key="item.key"
        class="sort-item"
        :class="{ active: currentSort === item.key }"
        @click="onSortChange(item.key)"
      >
        <text>{{ item.label }}</text>
        <uni-icons
          v-if="item.key === 'price'"
          :type="sortOrder === 'asc' ? 'arrow-up' : 'arrow-down'"
          size="12"
          :color="currentSort === 'price' ? '#0f3a57' : '#999'"
        />
      </view>
    </view>

    <!-- 2-column grid -->
    <view class="goods-grid">
      <view
        v-for="item in goodsList"
        :key="item.id"
        class="goods-card"
        @click="goDetail(item.id)"
      >
        <image class="goods-img" :src="item.image || item.image_url || (item.images && item.images[0]) || 'https://placehold.co/600x600/0f3a57/ffffff/png?text=GUOYUN&font=roboto'" mode="aspectFill" lazy-load />
        <view class="goods-info">
          <text class="goods-name">{{ item.name }}</text>
          <view class="goods-bottom">
            <text class="goods-price">¥{{ item.price }}</text>
            <text class="goods-sales">{{ t('goods.sales') }} {{ item.sales }}</text>
          </view>
        </view>
      </view>
    </view>

    <!-- Load more -->
    <view class="load-more">
      <text v-if="loading">{{ t('common.loading') }}</text>
      <text v-else-if="finished">{{ t('common.noData') }}</text>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { t } from '@/locale'
import { goodsApi } from '@/api/goods'

const goodsList = ref<any[]>([])
const loading = ref(false)
const finished = ref(false)
const page = ref(1)
const limit = 10
const currentSort = ref('default')
const sortOrder = ref('desc')
const categoryId = ref('')

const sortOptions = [
  { key: 'default', label: '綜合' },
  { key: 'sales', label: t('goods.sales') },
  { key: 'price', label: t('goods.price') },
  { key: 'newest', label: '最新' },
]

onMounted(() => {
  const pages = getCurrentPages()
  const currentPage = pages[pages.length - 1] as any
  const options = currentPage?.$page?.options || {}
  if (options.category_id) {
    categoryId.value = options.category_id
  }
  loadGoods()
})

async function loadGoods() {
  if (loading.value || finished.value) return
  loading.value = true
  try {
    const params: any = {
      page: page.value,
      limit,
      sort: currentSort.value,
      order: sortOrder.value,
    }
    if (categoryId.value) {
      params.category_id = categoryId.value
    }
    const res = await goodsApi.getList(params)
    const list = res?.data?.list || res?.data || []
    if (list.length < limit) {
      finished.value = true
    }
    goodsList.value.push(...list)
    page.value++
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

function onSortChange(key: string) {
  if (currentSort.value === key && key === 'price') {
    sortOrder.value = sortOrder.value === 'asc' ? 'desc' : 'asc'
  } else {
    currentSort.value = key
    sortOrder.value = 'desc'
  }
  resetList()
}

function resetList() {
  page.value = 1
  finished.value = false
  goodsList.value = []
  loadGoods()
}

function goDetail(id: string) {
  uni.navigateTo({ url: `/pages/goods/detail?id=${id}` })
}

function goSearch() {
  uni.navigateTo({ url: '/pages/search/index' })
}

// Pull-down refresh
onPullDownRefresh(() => {
  resetList()
  uni.stopPullDownRefresh()
})

// Load more
onReachBottom(() => {
  loadGoods()
})
</script>

<style scoped>
.goods-list {
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

.search-placeholder {
  color: #999;
  font-size: 28rpx;
}

.sort-bar {
  display: flex;
  background: #fff;
  padding: 20rpx 0;
  margin-bottom: 16rpx;
}

.sort-item {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 26rpx;
  color: #666;
  gap: 4rpx;
}

.sort-item.active {
  color: #0f3a57;
  font-weight: 600;
}

.goods-grid {
  display: flex;
  flex-wrap: wrap;
  padding: 0 16rpx;
  gap: 16rpx;
}

.goods-card {
  width: calc(50% - 24rpx);
  background: #fff;
  border-radius: 16rpx;
  overflow: hidden;
}

.goods-img {
  width: 100%;
  height: 340rpx;
}

.goods-info {
  padding: 16rpx 20rpx;
}

.goods-name {
  font-size: 26rpx;
  color: #333;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  line-height: 1.4;
}

.goods-bottom {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-top: 12rpx;
}

.goods-price {
  font-size: 32rpx;
  color: #e64340;
  font-weight: 600;
}

.goods-sales {
  font-size: 22rpx;
  color: #999;
}

.load-more {
  text-align: center;
  padding: 30rpx 0;
  font-size: 24rpx;
  color: #999;
}
</style>
