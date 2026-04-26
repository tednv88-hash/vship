<template>
  <view class="review-list-page">
    <!-- Overall rating -->
    <view class="rating-summary" v-if="summary">
      <view class="rating-big">
        <text class="rating-score">{{ summary.avg_rating?.toFixed(1) || '0.0' }}</text>
        <text class="rating-max">/5</text>
      </view>
      <view class="rating-stars">
        <view
          v-for="i in 5"
          :key="i"
          class="star"
          :class="{ filled: i <= Math.round(summary.avg_rating || 0) }"
        >
          <text class="star-icon">{{ i <= Math.round(summary.avg_rating || 0) ? '\u2605' : '\u2606' }}</text>
        </view>
      </view>
      <text class="rating-count">{{ summary.total || 0 }} reviews</text>
    </view>

    <!-- Loading -->
    <view v-if="loading && list.length === 0" class="loading-wrap">
      <text class="loading-text">{{ t('common.loading') }}</text>
    </view>

    <!-- Error -->
    <view v-else-if="error && list.length === 0" class="error-wrap">
      <text class="error-text">{{ error }}</text>
      <view class="retry-btn" @click="loadReviews(true)">
        <text class="retry-text">{{ t('common.retry') }}</text>
      </view>
    </view>

    <!-- Reviews list -->
    <view v-else class="review-items">
      <view v-for="item in list" :key="item.id" class="review-item">
        <view class="reviewer-info">
          <image
            class="reviewer-avatar"
            :src="item.avatar || '/static/avatar-default.png'"
            mode="aspectFill"
          />
          <view class="reviewer-detail">
            <text class="reviewer-name">{{ item.nickname }}</text>
            <view class="review-stars">
              <text
                v-for="i in 5"
                :key="i"
                class="star-sm"
                :class="{ filled: i <= item.rating }"
              >{{ i <= item.rating ? '\u2605' : '\u2606' }}</text>
            </view>
          </view>
          <text class="review-date">{{ item.date }}</text>
        </view>

        <text class="review-text" v-if="item.content">{{ item.content }}</text>

        <!-- Photos -->
        <view v-if="item.photos && item.photos.length > 0" class="review-photos">
          <image
            v-for="(photo, idx) in item.photos"
            :key="idx"
            class="review-photo"
            :src="photo"
            mode="aspectFill"
            @click="previewPhoto(item.photos, idx)"
          />
        </view>
      </view>

      <!-- Empty -->
      <view v-if="!loading && list.length === 0" class="empty-wrap">
        <text class="empty-text">{{ t('common.noData') }}</text>
      </view>

      <!-- Load more -->
      <view v-if="loadingMore" class="loadmore-wrap">
        <text class="loadmore-text">{{ t('common.loading') }}</text>
      </view>
      <view v-else-if="noMore && list.length > 0" class="loadmore-wrap">
        <text class="loadmore-text">-- {{ t('common.noData') }} --</text>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, onMounted, getCurrentInstance } from 'vue'
import { t } from '@/locale'
import { commonApi } from '@/api/common'

interface ReviewItem {
  id: string
  nickname: string
  avatar: string
  rating: number
  content: string
  photos: string[]
  date: string
}

interface ReviewSummary {
  avg_rating: number
  total: number
}

const loading = ref(true)
const loadingMore = ref(false)
const error = ref('')
const list = ref<ReviewItem[]>([])
const summary = ref<ReviewSummary | null>(null)
const page = ref(1)
const noMore = ref(false)
let goodsId = ''

async function loadReviews(reset = false) {
  if (reset) {
    page.value = 1
    noMore.value = false
    list.value = []
    loading.value = true
  }
  error.value = ''
  try {
    const res: any = await commonApi.getReviews({
      goods_id: goodsId,
      page: page.value,
      per_page: 10,
    })
    const data = res?.data || res
    if (data?.summary) {
      summary.value = data.summary
    }
    const items: ReviewItem[] = Array.isArray(data) ? data : data?.list || []
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
  }
}

function previewPhoto(photos: string[], idx: number) {
  uni.previewImage({
    current: photos[idx],
    urls: photos,
  })
}

onMounted(() => {
  uni.setNavigationBarTitle({ title: t('review.list') })
  const instance = getCurrentInstance()
  const options = (instance?.proxy as any)?.$page?.options || {}
  goodsId = options.goods_id || ''
  loadReviews(true)
})
</script>

<style scoped>
.review-list-page {
  min-height: 100vh;
  background-color: #f5f5f5;
}

.rating-summary {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 40rpx 32rpx;
  background-color: #fff;
  margin-bottom: 20rpx;
}

.rating-big {
  display: flex;
  align-items: baseline;
}

.rating-score {
  font-size: 72rpx;
  font-weight: 700;
  color: #0f3a57;
}

.rating-max {
  font-size: 28rpx;
  color: #999;
  margin-left: 4rpx;
}

.rating-stars {
  display: flex;
  margin: 12rpx 0;
}

.star {
  margin: 0 4rpx;
}

.star-icon {
  font-size: 36rpx;
  color: #ddd;
}

.star.filled .star-icon {
  color: #f5a623;
}

.rating-count {
  font-size: 24rpx;
  color: #999;
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

.review-items {
  padding: 0 24rpx;
}

.review-item {
  background-color: #fff;
  border-radius: 12rpx;
  padding: 24rpx;
  margin-bottom: 16rpx;
}

.reviewer-info {
  display: flex;
  align-items: center;
}

.reviewer-avatar {
  width: 64rpx;
  height: 64rpx;
  border-radius: 50%;
  margin-right: 16rpx;
}

.reviewer-detail {
  flex: 1;
}

.reviewer-name {
  font-size: 26rpx;
  color: #333;
  font-weight: 500;
}

.review-stars {
  display: flex;
  margin-top: 4rpx;
}

.star-sm {
  font-size: 24rpx;
  color: #ddd;
}

.star-sm.filled {
  color: #f5a623;
}

.review-date {
  font-size: 22rpx;
  color: #999;
}

.review-text {
  font-size: 28rpx;
  color: #333;
  line-height: 1.6;
  margin-top: 16rpx;
  display: block;
}

.review-photos {
  display: flex;
  flex-wrap: wrap;
  margin-top: 16rpx;
  gap: 12rpx;
}

.review-photo {
  width: 160rpx;
  height: 160rpx;
  border-radius: 8rpx;
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
