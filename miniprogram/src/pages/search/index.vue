<template>
  <view class="search-page">
    <!-- Search bar -->
    <view class="search-header">
      <view class="search-input-wrap">
        <uni-icons type="search" size="18" color="#999" />
        <input
          class="search-input"
          :placeholder="t('common.search')"
          v-model="keyword"
          :focus="true"
          confirm-type="search"
          @confirm="doSearch"
        />
        <uni-icons
          v-if="keyword"
          type="clear"
          size="18"
          color="#999"
          @click="keyword = ''"
        />
      </view>
      <text class="search-cancel" @click="goBack">{{ t('common.cancel') }}</text>
    </view>

    <!-- Search history -->
    <view v-if="!keyword && !showResults" class="section">
      <view class="section-header">
        <text class="section-title">搜索歷史</text>
        <uni-icons
          v-if="searchHistory.length > 0"
          type="trash"
          size="18"
          color="#999"
          @click="clearHistory"
        />
      </view>
      <view class="tag-list">
        <view
          v-for="(item, idx) in searchHistory"
          :key="idx"
          class="tag-item"
          @click="searchByTag(item)"
        >
          <text>{{ item }}</text>
        </view>
      </view>
      <view v-if="searchHistory.length === 0" class="empty-hint">
        <text>暫無搜索歷史</text>
      </view>
    </view>

    <!-- Hot keywords -->
    <view v-if="!keyword && !showResults" class="section">
      <view class="section-header">
        <text class="section-title">熱門搜索</text>
      </view>
      <view class="tag-list">
        <view
          v-for="(item, idx) in hotKeywords"
          :key="idx"
          class="tag-item hot"
          @click="searchByTag(item)"
        >
          <text>{{ item }}</text>
        </view>
      </view>
    </view>

    <!-- Search results (2-column grid) -->
    <view v-if="showResults" class="results-section">
      <view v-if="resultList.length === 0 && !loading" class="empty-result">
        <text>{{ t('common.noData') }}</text>
      </view>
      <view class="goods-grid">
        <view
          v-for="item in resultList"
          :key="item.id"
          class="goods-card"
          @click="goDetail(item.id)"
        >
          <image class="goods-img" :src="item.image" mode="aspectFill" lazy-load />
          <view class="goods-info">
            <text class="goods-name">{{ item.name }}</text>
            <view class="goods-bottom">
              <text class="goods-price">¥{{ item.price }}</text>
              <text class="goods-sales">{{ t('goods.sales') }} {{ item.sales }}</text>
            </view>
          </view>
        </view>
      </view>
      <view class="load-more">
        <text v-if="loading">{{ t('common.loading') }}</text>
        <text v-else-if="finished">— 已經到底了 —</text>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { t } from '@/locale'
import { goodsApi } from '@/api/goods'

const keyword = ref('')
const searchHistory = ref<string[]>([])
const hotKeywords = ref<string[]>([])
const resultList = ref<any[]>([])
const showResults = ref(false)
const loading = ref(false)
const finished = ref(false)
const page = ref(1)
const limit = 10

onMounted(() => {
  loadHistory()
  loadHotKeywords()
})

function loadHistory() {
  try {
    const stored = uni.getStorageSync('search_history')
    if (stored) {
      searchHistory.value = JSON.parse(stored)
    }
  } catch (e) {
    // ignore
  }
}

function saveHistory(kw: string) {
  const list = searchHistory.value.filter((item) => item !== kw)
  list.unshift(kw)
  if (list.length > 20) list.length = 20
  searchHistory.value = list
  uni.setStorageSync('search_history', JSON.stringify(list))
}

function clearHistory() {
  uni.showModal({
    title: '提示',
    content: '確定清除搜索歷史？',
    success(res) {
      if (res.confirm) {
        searchHistory.value = []
        uni.removeStorageSync('search_history')
      }
    },
  })
}

async function loadHotKeywords() {
  // Hot keywords could come from backend or be static
  hotKeywords.value = ['日本代購', '韓國零食', '美妝護膚', '母嬰用品', '電子產品', '保健品']
}

function searchByTag(kw: string) {
  keyword.value = kw
  doSearch()
}

async function doSearch() {
  const kw = keyword.value.trim()
  if (!kw) return
  saveHistory(kw)
  showResults.value = true
  page.value = 1
  finished.value = false
  resultList.value = []
  await loadResults()
}

async function loadResults() {
  if (loading.value || finished.value) return
  loading.value = true
  try {
    const res = await goodsApi.search({
      keyword: keyword.value.trim(),
      page: page.value,
      limit,
    })
    const list = res?.data?.list || res?.data || []
    if (list.length < limit) {
      finished.value = true
    }
    resultList.value.push(...list)
    page.value++
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

function goDetail(id: string) {
  uni.navigateTo({ url: `/pages/goods/detail?id=${id}` })
}

function goBack() {
  uni.navigateBack()
}

onReachBottom(() => {
  if (showResults.value) {
    loadResults()
  }
})
</script>

<style scoped>
.search-page {
  min-height: 100vh;
  background: #f5f5f5;
}

.search-header {
  display: flex;
  align-items: center;
  padding: 16rpx 24rpx;
  background: #fff;
  gap: 16rpx;
  position: sticky;
  top: 0;
  z-index: 10;
}

.search-input-wrap {
  flex: 1;
  display: flex;
  align-items: center;
  background: #f5f5f5;
  border-radius: 32rpx;
  padding: 14rpx 20rpx;
  gap: 12rpx;
}

.search-input {
  flex: 1;
  font-size: 28rpx;
  color: #333;
}

.search-cancel {
  font-size: 28rpx;
  color: #0f3a57;
  white-space: nowrap;
}

.section {
  background: #fff;
  padding: 24rpx 30rpx;
  margin-bottom: 16rpx;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20rpx;
}

.section-title {
  font-size: 30rpx;
  color: #333;
  font-weight: 600;
}

.tag-list {
  display: flex;
  flex-wrap: wrap;
  gap: 16rpx;
}

.tag-item {
  padding: 12rpx 24rpx;
  background: #f5f5f5;
  border-radius: 24rpx;
}

.tag-item text {
  font-size: 26rpx;
  color: #666;
}

.tag-item.hot {
  background: rgba(15, 58, 87, 0.08);
}

.tag-item.hot text {
  color: #0f3a57;
}

.empty-hint {
  text-align: center;
  padding: 30rpx 0;
  color: #999;
  font-size: 26rpx;
}

.results-section {
  padding: 16rpx;
}

.empty-result {
  text-align: center;
  padding: 100rpx 0;
  color: #999;
  font-size: 28rpx;
}

.goods-grid {
  display: flex;
  flex-wrap: wrap;
  gap: 16rpx;
}

.goods-card {
  width: calc(50% - 8rpx);
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
