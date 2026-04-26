<template>
  <view class="help-page">
    <!-- Search bar -->
    <view class="search-bar">
      <view class="search-input-wrap">
        <uni-icons type="search" size="18" color="#999" />
        <input
          class="search-input"
          v-model="keyword"
          :placeholder="t('common.search')"
          confirm-type="search"
          @confirm="onSearch"
        />
        <uni-icons
          v-if="keyword"
          type="clear"
          size="16"
          color="#ccc"
          @click="keyword = ''; onSearch()"
        />
      </view>
    </view>

    <!-- Loading -->
    <view v-if="loading" class="loading-wrap">
      <text class="loading-text">{{ t('common.loading') }}</text>
    </view>

    <!-- Error state -->
    <view v-else-if="error" class="error-wrap">
      <text class="error-text">{{ error }}</text>
      <view class="retry-btn" @click="loadData">
        <text class="retry-text">{{ t('common.retry') }}</text>
      </view>
    </view>

    <!-- FAQ categories -->
    <view v-else class="faq-list">
      <view
        v-for="category in filteredCategories"
        :key="category.id"
        class="faq-category"
      >
        <view class="category-title">
          <text class="category-name">{{ category.name }}</text>
        </view>

        <view class="faq-items">
          <view
            v-for="item in category.items"
            :key="item.id"
            class="faq-item"
          >
            <view class="faq-question" @click="toggleItem(item.id)">
              <text class="question-text">{{ item.question }}</text>
              <uni-icons
                :type="expandedIds.includes(item.id) ? 'up' : 'down'"
                size="14"
                color="#999"
              />
            </view>
            <view
              v-if="expandedIds.includes(item.id)"
              class="faq-answer"
            >
              <rich-text :nodes="item.answer" />
            </view>
          </view>
        </view>
      </view>

      <!-- Empty -->
      <view v-if="filteredCategories.length === 0" class="empty-wrap">
        <text class="empty-text">{{ t('common.noData') }}</text>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { t } from '@/locale'
import { commonApi } from '@/api/common'

interface FaqItem {
  id: string
  question: string
  answer: string
  category_id: string
}

interface FaqCategory {
  id: string
  name: string
  items: FaqItem[]
}

const loading = ref(true)
const error = ref('')
const keyword = ref('')
const categories = ref<FaqCategory[]>([])
const expandedIds = ref<string[]>([])

const filteredCategories = computed(() => {
  if (!keyword.value.trim()) return categories.value
  const kw = keyword.value.trim().toLowerCase()
  return categories.value
    .map((cat) => ({
      ...cat,
      items: cat.items.filter(
        (item) =>
          item.question.toLowerCase().includes(kw) ||
          item.answer.toLowerCase().includes(kw)
      ),
    }))
    .filter((cat) => cat.items.length > 0)
})

function toggleItem(id: string) {
  const idx = expandedIds.value.indexOf(id)
  if (idx >= 0) {
    expandedIds.value.splice(idx, 1)
  } else {
    expandedIds.value.push(id)
  }
}

function onSearch() {
  // Filtering is reactive via computed
}

async function loadData() {
  loading.value = true
  error.value = ''
  try {
    const res: any = await commonApi.getHelpList()
    const data = res?.data || res
    if (Array.isArray(data)) {
      // Flat list — group by category
      const map = new Map<string, FaqCategory>()
      for (const item of data) {
        const catId = item.category_id || 'default'
        if (!map.has(catId)) {
          map.set(catId, {
            id: catId,
            name: item.category_name || t('common.all'),
            items: [],
          })
        }
        map.get(catId)!.items.push(item)
      }
      categories.value = Array.from(map.values())
    } else if (data?.categories) {
      categories.value = data.categories
    }
  } catch (e: any) {
    error.value = e?.message || t('common.retry')
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  uni.setNavigationBarTitle({ title: t('help.title') })
  loadData()
})
</script>

<style scoped>
.help-page {
  min-height: 100vh;
  background-color: #f5f5f5;
}

.search-bar {
  padding: 20rpx 24rpx;
  background-color: #fff;
}

.search-input-wrap {
  display: flex;
  align-items: center;
  background-color: #f5f5f5;
  border-radius: 36rpx;
  padding: 0 24rpx;
  height: 72rpx;
}

.search-input {
  flex: 1;
  margin-left: 12rpx;
  font-size: 28rpx;
  color: #333;
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

.faq-category {
  margin-top: 20rpx;
}

.category-title {
  padding: 24rpx 32rpx;
  background-color: #fff;
  border-bottom: 1rpx solid #eee;
}

.category-name {
  font-size: 32rpx;
  font-weight: 600;
  color: #0f3a57;
}

.faq-items {
  background-color: #fff;
}

.faq-item {
  border-bottom: 1rpx solid #f0f0f0;
}

.faq-question {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 28rpx 32rpx;
}

.question-text {
  flex: 1;
  font-size: 28rpx;
  color: #333;
  margin-right: 16rpx;
}

.faq-answer {
  padding: 0 32rpx 28rpx;
  font-size: 26rpx;
  color: #666;
  line-height: 1.6;
}
</style>
