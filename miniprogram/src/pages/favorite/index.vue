<template>
  <view class="favorite-page">
    <!-- Loading -->
    <view v-if="loading && list.length === 0" class="loading-wrap">
      <text class="loading-text">{{ t('common.loading') }}</text>
    </view>

    <!-- Error -->
    <view v-else-if="error && list.length === 0" class="error-wrap">
      <text class="error-text">{{ error }}</text>
      <view class="retry-btn" @click="loadFavorites">
        <text class="retry-text">{{ t('common.retry') }}</text>
      </view>
    </view>

    <!-- Empty -->
    <view v-else-if="!loading && list.length === 0" class="empty-wrap">
      <text class="empty-icon">\u2661</text>
      <text class="empty-text">{{ t('common.noData') }}</text>
      <view class="empty-btn" @click="goHome">
        <text class="empty-btn-text">Go Shopping</text>
      </view>
    </view>

    <!-- Favorites list -->
    <view v-else class="fav-list">
      <view
        v-for="item in list"
        :key="item.id"
        class="fav-item"
      >
        <view class="fav-card" @click="goDetail(item.goods_id)">
          <image class="fav-image" :src="item.image" mode="aspectFill" />
          <view class="fav-info">
            <text class="fav-name">{{ item.name }}</text>
            <text class="fav-price">{{ item.price }}</text>
          </view>
        </view>
        <view class="fav-remove" @click="removeFavorite(item.id, item)">
          <text class="remove-text">{{ t('common.delete') }}</text>
        </view>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { t } from '@/locale'
import { commonApi } from '@/api/common'

interface FavoriteItem {
  id: string
  goods_id: string
  name: string
  image: string
  price: string
}

const loading = ref(true)
const error = ref('')
const list = ref<FavoriteItem[]>([])

async function loadFavorites() {
  loading.value = true
  error.value = ''
  try {
    const res: any = await commonApi.getFavorites()
    const data = res?.data || res
    list.value = Array.isArray(data) ? data : data?.list || []
  } catch (e: any) {
    error.value = e?.message || t('common.retry')
  } finally {
    loading.value = false
  }
}

async function removeFavorite(id: string, item: FavoriteItem) {
  uni.showModal({
    title: t('common.confirm'),
    content: `Remove "${item.name}" from favorites?`,
    success: async (res) => {
      if (res.confirm) {
        try {
          await commonApi.removeFavorite(id)
          list.value = list.value.filter((f) => f.id !== id)
          uni.showToast({ title: t('common.done'), icon: 'success' })
        } catch (e: any) {
          uni.showToast({ title: e?.message || t('common.retry'), icon: 'none' })
        }
      }
    },
  })
}

function goDetail(goodsId: string) {
  uni.navigateTo({ url: `/pages/goods/detail?id=${goodsId}` })
}

function goHome() {
  uni.switchTab({ url: '/pages/index/index' })
}

onMounted(() => {
  uni.setNavigationBarTitle({ title: t('user.favorite') })
  loadFavorites()
})
</script>

<style scoped>
.favorite-page {
  min-height: 100vh;
  background-color: #f5f5f5;
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

.empty-wrap {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 160rpx 0;
}

.empty-icon {
  font-size: 96rpx;
  color: #ddd;
  margin-bottom: 24rpx;
}

.empty-text {
  font-size: 28rpx;
  color: #999;
  margin-bottom: 32rpx;
}

.empty-btn {
  padding: 16rpx 48rpx;
  background-color: #0f3a57;
  border-radius: 44rpx;
}

.empty-btn-text {
  font-size: 28rpx;
  color: #fff;
}

.fav-list {
  padding: 20rpx 24rpx;
}

.fav-item {
  display: flex;
  background-color: #fff;
  border-radius: 12rpx;
  margin-bottom: 16rpx;
  overflow: hidden;
}

.fav-card {
  flex: 1;
  display: flex;
  padding: 20rpx;
}

.fav-image {
  width: 160rpx;
  height: 160rpx;
  border-radius: 8rpx;
  margin-right: 20rpx;
  flex-shrink: 0;
}

.fav-info {
  flex: 1;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
}

.fav-name {
  font-size: 28rpx;
  color: #333;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.fav-price {
  font-size: 30rpx;
  font-weight: 600;
  color: #e64340;
}

.fav-remove {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 120rpx;
  background-color: #e64340;
}

.remove-text {
  font-size: 24rpx;
  color: #fff;
}
</style>
