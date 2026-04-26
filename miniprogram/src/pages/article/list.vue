<template>
  <view class="article-list-page">
    <!-- Loading initial -->
    <view v-if="loading && list.length === 0" class="loading-wrap">
      <text class="loading-text">{{ t('common.loading') }}</text>
    </view>

    <!-- Error -->
    <view v-else-if="error && list.length === 0" class="error-wrap">
      <text class="error-text">{{ error }}</text>
      <view class="retry-btn" @click="loadArticles(true)">
        <text class="retry-text">{{ t('common.retry') }}</text>
      </view>
    </view>

    <!-- List -->
    <scroll-view
      v-else
      class="article-scroll"
      scroll-y
      refresher-enabled
      :refresher-triggered="refreshing"
      @refresherrefresh="onRefresh"
      @scrolltolower="onLoadMore"
    >
      <view class="article-items">
        <view
          v-for="item in list"
          :key="item.id"
          class="article-item"
          @click="goDetail(item.id)"
        >
          <image
            v-if="item.thumbnail"
            class="article-thumb"
            :src="item.thumbnail"
            mode="aspectFill"
          />
          <view class="article-info">
            <text class="article-title">{{ item.title }}</text>
            <view class="article-meta">
              <text class="meta-date">{{ item.date }}</text>
              <text class="meta-views">{{ item.views }} views</text>
            </view>
          </view>
        </view>
      </view>

      <!-- Load more -->
      <view v-if="loadingMore" class="loadmore-wrap">
        <text class="loadmore-text">{{ t('common.loading') }}</text>
      </view>
      <view v-else-if="noMore && list.length > 0" class="loadmore-wrap">
        <text class="loadmore-text">-- {{ t('common.noData') }} --</text>
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

interface Article {
  id: string
  title: string
  thumbnail: string
  date: string
  views: number
}

const loading = ref(true)
const loadingMore = ref(false)
const refreshing = ref(false)
const error = ref('')
const list = ref<Article[]>([])
const page = ref(1)
const noMore = ref(false)

async function loadArticles(reset = false) {
  if (reset) {
    page.value = 1
    noMore.value = false
    list.value = []
    loading.value = true
  }
  error.value = ''
  try {
    const res: any = await commonApi.getArticles({ page: page.value, per_page: 10 })
    const data = res?.data || res
    const items: Article[] = Array.isArray(data) ? data : data?.list || []
    if (items.length < 10) noMore.value = true
    if (reset) {
      list.value = items
    } else {
      list.value.push(...items)
    }
  } catch (e: any) {
    error.value = e?.message || t('common.retry')
  } finally {
    loading.value = false
    loadingMore.value = false
    refreshing.value = false
  }
}

function onRefresh() {
  refreshing.value = true
  loadArticles(true)
}

function onLoadMore() {
  if (loadingMore.value || noMore.value) return
  loadingMore.value = true
  page.value++
  loadArticles()
}

function goDetail(id: string) {
  uni.navigateTo({ url: `/pages/article/detail?id=${id}` })
}

onMounted(() => {
  uni.setNavigationBarTitle({ title: t('article.list') })
  loadArticles(true)
})
</script>

<style scoped>
.article-list-page {
  min-height: 100vh;
  background-color: #f5f5f5;
}

.article-scroll {
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

.article-items {
  padding: 20rpx 24rpx;
}

.article-item {
  display: flex;
  background-color: #fff;
  border-radius: 12rpx;
  margin-bottom: 20rpx;
  overflow: hidden;
}

.article-thumb {
  width: 200rpx;
  height: 150rpx;
  flex-shrink: 0;
}

.article-info {
  flex: 1;
  padding: 20rpx 24rpx;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
}

.article-title {
  font-size: 28rpx;
  color: #333;
  font-weight: 500;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.article-meta {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-top: 12rpx;
}

.meta-date,
.meta-views {
  font-size: 24rpx;
  color: #999;
}

.loadmore-wrap {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 32rpx 0;
}

.loadmore-text {
  font-size: 24rpx;
  color: #999;
}
</style>
