<template>
  <view class="article-detail-page">
    <!-- Loading -->
    <view v-if="loading" class="loading-wrap">
      <text class="loading-text">{{ t('common.loading') }}</text>
    </view>

    <!-- Error -->
    <view v-else-if="error" class="error-wrap">
      <text class="error-text">{{ error }}</text>
      <view class="retry-btn" @click="loadDetail">
        <text class="retry-text">{{ t('common.retry') }}</text>
      </view>
    </view>

    <!-- Content -->
    <view v-else class="detail-content">
      <view class="detail-header">
        <text class="detail-title">{{ detail.title }}</text>
        <view class="detail-meta">
          <text class="meta-author">{{ detail.author }}</text>
          <text class="meta-date">{{ detail.date }}</text>
        </view>
      </view>
      <view class="detail-body">
        <rich-text :nodes="detail.content" />
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, onMounted, getCurrentInstance } from 'vue'
import { t } from '@/locale'
import { commonApi } from '@/api/common'

interface ArticleDetail {
  title: string
  author: string
  date: string
  content: string
}

const loading = ref(true)
const error = ref('')
const detail = ref<ArticleDetail>({
  title: '',
  author: '',
  date: '',
  content: '',
})

let articleId = ''

async function loadDetail() {
  if (!articleId) return
  loading.value = true
  error.value = ''
  try {
    const res: any = await commonApi.getArticleDetail(articleId)
    const data = res?.data || res
    detail.value = {
      title: data.title || '',
      author: data.author || '',
      date: data.created_at || data.date || '',
      content: data.content || '',
    }
  } catch (e: any) {
    error.value = e?.message || t('common.retry')
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  uni.setNavigationBarTitle({ title: t('article.detail') })
  const instance = getCurrentInstance()
  const options = (instance?.proxy as any)?.$page?.options || {}
  articleId = options.id || ''
  loadDetail()
})
</script>

<style scoped>
.article-detail-page {
  min-height: 100vh;
  background-color: #fff;
}

.loading-wrap,
.error-wrap {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 120rpx 0;
}

.loading-text {
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

.detail-header {
  padding: 32rpx;
  border-bottom: 1rpx solid #eee;
}

.detail-title {
  font-size: 36rpx;
  font-weight: 600;
  color: #333;
  line-height: 1.5;
}

.detail-meta {
  display: flex;
  align-items: center;
  margin-top: 16rpx;
}

.meta-author {
  font-size: 24rpx;
  color: #666;
  margin-right: 24rpx;
}

.meta-date {
  font-size: 24rpx;
  color: #999;
}

.detail-body {
  padding: 32rpx;
  font-size: 28rpx;
  color: #333;
  line-height: 1.8;
}
</style>
