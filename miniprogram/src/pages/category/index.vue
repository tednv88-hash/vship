<template>
  <view class="page">
    <view class="category-container">
      <!-- Left Sidebar -->
      <scroll-view class="sidebar" scroll-y>
        <view
          v-for="(cat, idx) in categoryList"
          :key="cat.id"
          class="sidebar-item"
          :class="{ active: activeIndex === idx }"
          @tap="onCategoryTap(idx)"
        >
          <view v-if="activeIndex === idx" class="sidebar-indicator" />
          <text class="sidebar-text">{{ cat.name }}</text>
        </view>
        <!-- Empty state for sidebar -->
        <view v-if="categoryList.length === 0 && !loading" class="sidebar-empty">
          <text class="sidebar-empty-text">{{ t('common.noData') }}</text>
        </view>
      </scroll-view>

      <!-- Right Content -->
      <scroll-view class="content" scroll-y :scroll-into-view="scrollIntoId">
        <!-- Loading -->
        <view v-if="loading" class="content-loading">
          <text class="loading-text">{{ t('common.loading') }}</text>
        </view>

        <template v-else-if="currentCategory">
          <!-- Category Banner (optional) -->
          <image
            v-if="currentCategory.image"
            class="category-banner"
            :src="currentCategory.image"
            mode="aspectFill"
          />

          <!-- Subcategories -->
          <view v-if="currentCategory.children && currentCategory.children.length > 0" class="sub-category-section">
            <view class="sub-grid">
              <view
                v-for="sub in currentCategory.children"
                :key="sub.id"
                class="sub-item"
                @tap="onSubCategoryTap(sub)"
              >
                <image
                  v-if="sub.image"
                  class="sub-icon"
                  :src="sub.image"
                  mode="aspectFill"
                />
                <view v-else class="sub-icon sub-icon-placeholder">
                  <text class="sub-icon-text">{{ sub.name.charAt(0) }}</text>
                </view>
                <text class="sub-name">{{ sub.name }}</text>
              </view>
            </view>
          </view>

          <!-- Goods list under selected category -->
          <view v-if="categoryGoods.length > 0" class="goods-section">
            <view class="goods-section-title">
              <text class="goods-section-title-text">推薦商品</text>
            </view>
            <view class="goods-list">
              <view
                v-for="goods in categoryGoods"
                :key="goods.id"
                class="goods-item"
                @tap="goToGoodsDetail(goods.id)"
              >
                <image class="goods-image" :src="goods.image" mode="aspectFill" />
                <view class="goods-info">
                  <text class="goods-name">{{ goods.name }}</text>
                  <text class="goods-desc" v-if="goods.desc">{{ goods.desc }}</text>
                  <view class="goods-bottom">
                    <text class="goods-price">¥{{ goods.price }}</text>
                    <text class="goods-sales">已售 {{ goods.sales || 0 }}</text>
                  </view>
                </view>
              </view>
            </view>
          </view>

          <!-- Empty goods -->
          <view v-if="!loading && categoryGoods.length === 0 && (!currentCategory.children || currentCategory.children.length === 0)" class="empty-state">
            <text class="empty-text">{{ t('common.noData') }}</text>
          </view>
        </template>
      </scroll-view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { onShow } from '@dcloudio/uni-app'
import { t } from '@/locale'
import { goodsApi } from '@/api/goods'

interface SubCategory {
  id: string
  name: string
  image?: string
}

interface Category {
  id: string
  name: string
  image?: string
  children?: SubCategory[]
}

interface GoodsItem {
  id: string
  name: string
  image: string
  price: string
  desc?: string
  sales?: number
}

const loading = ref(true)
const activeIndex = ref(0)
const scrollIntoId = ref('')
const categoryList = ref<Category[]>([])
const categoryGoods = ref<GoodsItem[]>([])

const currentCategory = computed(() => {
  return categoryList.value[activeIndex.value] || null
})

onShow(() => {
  if (categoryList.value.length === 0) {
    fetchCategories()
  }
})

async function fetchCategories() {
  loading.value = true
  try {
    const res: any = await goodsApi.getCategories()
    const data = res?.data || res
    const list = data?.list || data?.items || (Array.isArray(data) ? data : [])
    categoryList.value = list

    if (list.length > 0) {
      await fetchCategoryGoods(list[0].id)
    }
  } catch (e) {
    // Use placeholder categories
    categoryList.value = [
      { id: '1', name: '熱門推薦', children: [
        { id: '1-1', name: '日用品' },
        { id: '1-2', name: '數碼產品' },
        { id: '1-3', name: '服飾鞋包' },
        { id: '1-4', name: '美妝個護' },
        { id: '1-5', name: '食品飲料' },
        { id: '1-6', name: '母嬰用品' },
      ]},
      { id: '2', name: '數碼家電', children: [
        { id: '2-1', name: '手機配件' },
        { id: '2-2', name: '電腦周邊' },
        { id: '2-3', name: '智能穿戴' },
      ]},
      { id: '3', name: '服飾鞋包', children: [
        { id: '3-1', name: '男裝' },
        { id: '3-2', name: '女裝' },
        { id: '3-3', name: '運動鞋' },
        { id: '3-4', name: '箱包' },
      ]},
      { id: '4', name: '美妝個護', children: [
        { id: '4-1', name: '護膚品' },
        { id: '4-2', name: '彩妝' },
        { id: '4-3', name: '洗護用品' },
      ]},
      { id: '5', name: '食品飲料', children: [
        { id: '5-1', name: '零食' },
        { id: '5-2', name: '茶飲' },
        { id: '5-3', name: '保健品' },
      ]},
      { id: '6', name: '家居日用', children: [
        { id: '6-1', name: '清潔用品' },
        { id: '6-2', name: '廚房用品' },
        { id: '6-3', name: '收納整理' },
      ]},
      { id: '7', name: '母嬰用品', children: [
        { id: '7-1', name: '奶粉' },
        { id: '7-2', name: '紙尿褲' },
        { id: '7-3', name: '童裝' },
      ]},
      { id: '8', name: '運動戶外', children: [
        { id: '8-1', name: '運動裝備' },
        { id: '8-2', name: '戶外用品' },
      ]},
    ]
  } finally {
    loading.value = false
  }
}

async function fetchCategoryGoods(categoryId: string) {
  try {
    const res: any = await goodsApi.getByCategory(categoryId, { page: 1, limit: 20 })
    const data = res?.data || res
    categoryGoods.value = data?.list || data?.items || (Array.isArray(data) ? data : [])
  } catch (e) {
    categoryGoods.value = []
  }
}

function onCategoryTap(idx: number) {
  activeIndex.value = idx
  scrollIntoId.value = ''
  const cat = categoryList.value[idx]
  if (cat) {
    fetchCategoryGoods(cat.id)
  }
}

function onSubCategoryTap(sub: SubCategory) {
  uni.navigateTo({ url: `/pages/goods/list?category_id=${sub.id}&title=${encodeURIComponent(sub.name)}` })
}

function goToGoodsDetail(id: string) {
  uni.navigateTo({ url: `/pages/goods/detail?id=${id}` })
}
</script>

<style scoped>
.page {
  height: 100vh;
  background-color: #f5f5f5;
  overflow: hidden;
}

.category-container {
  display: flex;
  height: 100%;
}

/* --- Sidebar --- */
.sidebar {
  width: 180rpx;
  height: 100%;
  background-color: #f2f2f2;
  flex-shrink: 0;
}

.sidebar-item {
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
  height: 100rpx;
  padding: 0 16rpx;
}

.sidebar-item.active {
  background-color: #ffffff;
}

.sidebar-indicator {
  position: absolute;
  left: 0;
  top: 50%;
  transform: translateY(-50%);
  width: 6rpx;
  height: 40rpx;
  background-color: #0f3a57;
  border-radius: 0 6rpx 6rpx 0;
}

.sidebar-text {
  font-size: 26rpx;
  color: #333333;
  text-align: center;
  line-height: 1.3;
}

.sidebar-item.active .sidebar-text {
  color: #0f3a57;
  font-weight: bold;
}

.sidebar-empty {
  padding: 40rpx 20rpx;
  text-align: center;
}

.sidebar-empty-text {
  font-size: 22rpx;
  color: #999999;
}

/* --- Content --- */
.content {
  flex: 1;
  height: 100%;
  background-color: #ffffff;
}

.content-loading {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 60rpx 0;
}

.loading-text {
  font-size: 26rpx;
  color: #999999;
}

/* Category banner */
.category-banner {
  width: 100%;
  height: 200rpx;
  display: block;
}

/* --- Subcategory Grid --- */
.sub-category-section {
  padding: 20rpx;
}

.sub-grid {
  display: flex;
  flex-wrap: wrap;
  gap: 10rpx;
}

.sub-item {
  width: calc(33.33% - 8rpx);
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 20rpx 0;
}

.sub-icon {
  width: 100rpx;
  height: 100rpx;
  border-radius: 16rpx;
  margin-bottom: 12rpx;
}

.sub-icon-placeholder {
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: #e8f0f5;
}

.sub-icon-text {
  font-size: 32rpx;
  color: #0f3a57;
  font-weight: bold;
}

.sub-name {
  font-size: 24rpx;
  color: #333333;
  text-align: center;
}

/* --- Goods Section --- */
.goods-section {
  padding: 0 20rpx;
}

.goods-section-title {
  padding: 24rpx 0 16rpx;
  border-top: 1rpx solid #eeeeee;
}

.goods-section-title-text {
  font-size: 28rpx;
  font-weight: bold;
  color: #333333;
}

.goods-list {
  display: flex;
  flex-direction: column;
  gap: 20rpx;
  padding-bottom: 40rpx;
}

.goods-item {
  display: flex;
  background-color: #ffffff;
  border-radius: 12rpx;
  overflow: hidden;
  border: 1rpx solid #f0f0f0;
}

.goods-image {
  width: 180rpx;
  height: 180rpx;
  flex-shrink: 0;
}

.goods-info {
  flex: 1;
  padding: 16rpx 20rpx;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
}

.goods-name {
  font-size: 26rpx;
  color: #333333;
  line-height: 1.4;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.goods-desc {
  font-size: 22rpx;
  color: #999999;
  margin-top: 6rpx;
  overflow: hidden;
  white-space: nowrap;
  text-overflow: ellipsis;
}

.goods-bottom {
  display: flex;
  align-items: baseline;
  justify-content: space-between;
  margin-top: 8rpx;
}

.goods-price {
  font-size: 30rpx;
  font-weight: bold;
  color: #e64340;
}

.goods-sales {
  font-size: 20rpx;
  color: #999999;
}

/* --- Empty --- */
.empty-state {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 80rpx 0;
}

.empty-text {
  font-size: 26rpx;
  color: #999999;
}
</style>
