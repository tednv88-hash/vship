<template>
  <view class="page">
    <!-- Custom Navigation Bar -->
    <view class="navbar" :style="{ paddingTop: statusBarHeight + 'px', backgroundColor: pageStyle.titleBackgroundColor || '#0f3a57' }">
      <view class="navbar-content">
        <image class="navbar-logo" src="/static/logo.svg" mode="heightFix" />
        <view class="navbar-right">
          <view class="navbar-icon" @tap="goToMessage">
            <text class="iconfont icon-bell">&#x1F514;</text>
            <view v-if="messageCount > 0" class="badge">{{ messageCount > 99 ? '99+' : messageCount }}</view>
          </view>
          <view class="navbar-icon" @tap="goToSearch">
            <text class="iconfont icon-search">&#x1F50D;</text>
          </view>
        </view>
      </view>
    </view>

    <scroll-view
      class="scroll-content"
      scroll-y
      :style="{ paddingTop: navbarHeight + 'px' }"
      refresher-enabled
      :refresher-triggered="isRefreshing"
      @refresherrefresh="onRefresh"
      @scrolltolower="loadMoreGoods"
    >
      <!-- Dynamic page_data components -->
      <template v-for="(item, idx) in pageItems" :key="idx">

        <!-- Banner -->
        <view v-if="item.type === 'banner' && item.data?.length" class="banner-wrapper">
          <swiper
            class="banner-swiper"
            :autoplay="true"
            :interval="parseInt(item.params?.interval) || 4000"
            :duration="500"
            circular
            indicator-dots
            :indicator-color="'rgba(255,255,255,0.4)'"
            :indicator-active-color="item.style?.btnColor || '#ffffff'"
          >
            <swiper-item v-for="(slide, si) in item.data" :key="si" @tap="onLinkTap(slide.linkUrl)">
              <image class="banner-image" :src="slide.imgUrl" mode="aspectFill" />
            </swiper-item>
          </swiper>
        </view>

        <!-- Notice -->
        <view v-else-if="item.type === 'notice'" class="notice-bar"
              :style="{ backgroundColor: item.style?.background || '#fffbe8' }"
              @tap="goToNoticeList">
          <image v-if="item.params?.icon" class="notice-icon-img" :src="item.params.icon" mode="aspectFit" />
          <text v-else class="notice-icon">&#x1F4E2;</text>
          <text class="notice-text" :style="{ color: item.style?.textColor || '#666' }">{{ item.params?.text || '' }}</text>
          <text class="notice-arrow">&#x276F;</text>
        </view>

        <!-- NavBar Grid -->
        <view v-else-if="item.type === 'navBar' && item.data?.length" class="nav-grid"
              :style="{ backgroundColor: item.style?.background || '#fff' }">
          <view
            v-for="(nav, ni) in item.data"
            :key="ni"
            class="nav-item"
            :style="{ width: (100 / (parseInt(item.style?.rowsNum) || 5)) + '%' }"
            @tap="onLinkTap(nav.linkUrl)"
          >
            <image class="nav-icon-img" :src="nav.imgUrl" mode="aspectFit" />
            <text class="nav-label" :style="{ color: nav.color || '#333' }">{{ nav.text }}</text>
          </view>
        </view>

        <!-- Search -->
        <view v-else-if="item.type === 'search'" class="search-bar"
              :style="{ backgroundColor: item.style?.background || '#fff' }"
              @tap="goToSearch">
          <view class="search-inner">
            <text class="search-placeholder">&#x1F50D; {{ item.style?.placeholder || '搜尋商品' }}</text>
          </view>
        </view>

        <!-- Image Single -->
        <view v-else-if="item.type === 'imageSingle' && item.data?.length" class="image-single"
              :style="{ backgroundColor: item.style?.background || '#fff', padding: (item.style?.paddingTop || 0) + 'px ' + (item.style?.paddingLeft || 0) + 'px' }">
          <image
            v-for="(img, ii) in item.data"
            :key="ii"
            class="single-image"
            :src="img.imgUrl"
            mode="widthFix"
            @tap="onLinkTap(img.linkUrl)"
          />
        </view>

        <!-- Goods -->
        <view v-else-if="item.type === 'goods'" class="section"
              :style="{ backgroundColor: item.style?.background || '#f5f5f5' }">
          <view class="section-header">
            <text class="section-title">精選好物</text>
            <view class="section-more" @tap="goToGoodsList">
              <text class="section-more-text">{{ t('common.more') }}</text>
              <text class="section-more-arrow">&#x276F;</text>
            </view>
          </view>
          <view class="goods-grid" :class="'goods-col-' + (item.style?.column || '2')">
            <view
              v-for="goods in (goodsList.length ? goodsList : getGoodsDefault(item))"
              :key="goods.id || goods.goods_name"
              class="goods-card"
              @tap="goods.id ? goToGoodsDetail(goods.id) : null"
            >
              <image class="goods-image" :src="goods.image || goods.goods_image || ''" mode="aspectFill" />
              <view class="goods-info">
                <text class="goods-name">{{ goods.name || goods.goods_name || '' }}</text>
                <view class="goods-price-row">
                  <text class="goods-price">¥{{ goods.price || goods.goods_price || '0' }}</text>
                  <text v-if="goods.original_price || (goods.line_price && parseFloat(goods.line_price) > 0)" class="goods-original-price">
                    ¥{{ goods.original_price || goods.line_price }}
                  </text>
                </view>
                <text v-if="goods.sales || goods.goods_sales" class="goods-sales">已售 {{ goods.sales || goods.goods_sales || 0 }}</text>
              </view>
            </view>
          </view>

          <!-- Load more -->
          <view class="load-more">
            <text v-if="goodsLoading" class="load-more-text">{{ t('common.loading') }}</text>
            <text v-else-if="goodsFinished" class="load-more-text">- {{ t('common.noData') }} -</text>
          </view>
        </view>

        <!-- Blank -->
        <view v-else-if="item.type === 'blank'"
              :style="{ height: (item.style?.height || 20) + 'px', backgroundColor: item.style?.background || 'transparent' }">
        </view>

        <!-- Rich Text -->
        <view v-else-if="item.type === 'richText'" class="richtext-wrapper">
          <rich-text :nodes="item.params?.content || ''" />
        </view>

      </template>
    </scroll-view>
  </view>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { onShow } from '@dcloudio/uni-app'
import { t } from '@/locale'
import { commonApi } from '@/api/common'
import { goodsApi } from '@/api/goods'
import store from '@/store'

// --- Status bar & navbar ---
const statusBarHeight = ref(20)
const navbarHeight = computed(() => statusBarHeight.value + 44)

// --- Data ---
const isRefreshing = ref(false)
const messageCount = computed(() => store.messageCount)

// Page design data
const pageItems = ref<any[]>([])
const pageStyle = ref<any>({})

interface GoodsItem {
  id: string
  name: string
  image: string
  price: string
  original_price?: string
  sales?: number
  // From page_data defaultData format
  goods_name?: string
  goods_image?: string
  goods_price?: string
  goods_sales?: string
  line_price?: string
}

const goodsList = ref<GoodsItem[]>([])
const goodsPage = ref(1)
const goodsLoading = ref(false)
const goodsFinished = ref(false)

// --- System info ---
onMounted(() => {
  const sysInfo = uni.getSystemInfoSync()
  statusBarHeight.value = sysInfo.statusBarHeight || 20
})

// --- Lifecycle ---
onShow(() => {
  loadPageData()
})

// --- Data loading ---
async function loadPageData() {
  await Promise.all([
    fetchPageDesign(),
    fetchGoods(true),
  ])
}

async function fetchPageDesign() {
  try {
    const res: any = await commonApi.getPageDesign({ is_default: true, type: 'home' })
    const rawData = res?.data || res

    // Handle list response (array or paginated)
    let pageDesign: any = null
    if (Array.isArray(rawData)) {
      pageDesign = rawData.find((d: any) => d.type === 'home' && d.is_default) || rawData[0]
    } else if (rawData?.list) {
      pageDesign = rawData.list.find((d: any) => d.type === 'home' && d.is_default) || rawData.list[0]
    } else if (rawData?.items) {
      pageDesign = rawData.items.find((d: any) => d.type === 'home' && d.is_default) || rawData.items[0]
    } else if (rawData?.page_data) {
      pageDesign = rawData
    }

    if (pageDesign?.page_data) {
      const pd = pageDesign.page_data
      pageStyle.value = pd.page?.style || {}
      pageItems.value = pd.items || []
    }
  } catch (e) {
    console.error('Failed to fetch page design:', e)
  }
}

function getGoodsDefault(item: any): any[] {
  return item.defaultData || []
}

async function fetchGoods(reset = false) {
  if (goodsLoading.value) return
  if (!reset && goodsFinished.value) return

  goodsLoading.value = true
  if (reset) {
    goodsPage.value = 1
    goodsFinished.value = false
  }

  try {
    const res: any = await goodsApi.getList({ page: goodsPage.value, limit: 10, is_featured: true })
    const data = res?.data || res
    const list: GoodsItem[] = data?.list || data?.items || (Array.isArray(data) ? data : [])

    if (reset) {
      goodsList.value = list
    } else {
      goodsList.value.push(...list)
    }

    if (list.length < 10) {
      goodsFinished.value = true
    } else {
      goodsPage.value++
    }
  } catch (e) {
    // Silent - will fall back to defaultData from page_data
  } finally {
    goodsLoading.value = false
  }
}

function loadMoreGoods() {
  fetchGoods(false)
}

async function onRefresh() {
  isRefreshing.value = true
  await loadPageData()
  isRefreshing.value = false
}

// --- Navigation ---
function onLinkTap(url?: string) {
  if (url) {
    uni.navigateTo({ url })
  }
}

function goToSearch() {
  uni.navigateTo({ url: '/pages/search/index' })
}

function goToMessage() {
  uni.navigateTo({ url: '/pages/message/index' })
}

function goToNoticeList() {
  uni.navigateTo({ url: '/pages/notice/list' })
}

function goToGoodsList() {
  uni.navigateTo({ url: '/pages/goods/list' })
}

function goToGoodsDetail(id: string) {
  uni.navigateTo({ url: `/pages/goods/detail?id=${id}` })
}
</script>

<style scoped>
.page {
  min-height: 100vh;
  background-color: #f5f5f5;
}

/* --- Navbar --- */
.navbar {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  z-index: 100;
}

.navbar-content {
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 44px;
  padding: 0 30rpx;
}

.navbar-logo {
  height: 28px;
}

.navbar-right {
  display: flex;
  align-items: center;
  gap: 24rpx;
}

.navbar-icon {
  position: relative;
  width: 44rpx;
  height: 44rpx;
  display: flex;
  align-items: center;
  justify-content: center;
}

.navbar-icon .iconfont {
  font-size: 36rpx;
  color: #ffffff;
  font-style: normal;
}

.badge {
  position: absolute;
  top: -8rpx;
  right: -12rpx;
  min-width: 28rpx;
  height: 28rpx;
  line-height: 28rpx;
  padding: 0 6rpx;
  border-radius: 28rpx;
  background-color: #ff4d4f;
  color: #ffffff;
  font-size: 18rpx;
  text-align: center;
}

/* --- Scroll Content --- */
.scroll-content {
  height: 100vh;
  box-sizing: border-box;
}

/* --- Banner --- */
.banner-wrapper {
  padding: 20rpx 24rpx 0;
}

.banner-swiper {
  width: 100%;
  height: 300rpx;
  border-radius: 16rpx;
  overflow: hidden;
}

.banner-image {
  width: 100%;
  height: 100%;
  border-radius: 16rpx;
}

/* --- Notice Bar --- */
.notice-bar {
  display: flex;
  align-items: center;
  margin: 20rpx 24rpx 0;
  padding: 16rpx 24rpx;
  border-radius: 12rpx;
}

.notice-icon {
  font-size: 32rpx;
  margin-right: 12rpx;
  flex-shrink: 0;
}

.notice-icon-img {
  width: 36rpx;
  height: 36rpx;
  margin-right: 12rpx;
  flex-shrink: 0;
}

.notice-text {
  flex: 1;
  font-size: 24rpx;
  line-height: 40rpx;
  overflow: hidden;
  white-space: nowrap;
  text-overflow: ellipsis;
}

.notice-arrow {
  font-size: 24rpx;
  color: #999999;
  margin-left: 12rpx;
  flex-shrink: 0;
}

/* --- Nav Grid --- */
.nav-grid {
  display: flex;
  flex-wrap: wrap;
  padding: 30rpx 24rpx 10rpx;
  margin: 20rpx 24rpx 0;
  border-radius: 16rpx;
}

.nav-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  margin-bottom: 24rpx;
}

.nav-icon-img {
  width: 88rpx;
  height: 88rpx;
  margin-bottom: 12rpx;
}

.nav-label {
  font-size: 22rpx;
  text-align: center;
  line-height: 1.2;
}

/* --- Search Bar --- */
.search-bar {
  padding: 16rpx 24rpx;
}

.search-inner {
  display: flex;
  align-items: center;
  padding: 16rpx 24rpx;
  background-color: #f0f0f0;
  border-radius: 40rpx;
}

.search-placeholder {
  font-size: 26rpx;
  color: #999;
}

/* --- Image Single --- */
.image-single {
  margin: 0;
}

.single-image {
  width: 100%;
}

/* --- Section / Goods --- */
.section {
  margin-top: 20rpx;
}

.section-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 24rpx 24rpx 16rpx;
}

.section-title {
  font-size: 32rpx;
  font-weight: bold;
  color: #333333;
}

.section-more {
  display: flex;
  align-items: center;
}

.section-more-text {
  font-size: 24rpx;
  color: #999999;
}

.section-more-arrow {
  font-size: 20rpx;
  color: #999999;
  margin-left: 4rpx;
}

/* --- Goods Grid --- */
.goods-grid {
  display: flex;
  flex-wrap: wrap;
  padding: 0 24rpx;
  gap: 16rpx;
}

.goods-col-2 .goods-card {
  width: calc(50% - 8rpx);
}

.goods-col-1 .goods-card {
  width: 100%;
}

.goods-col-3 .goods-card {
  width: calc(33.33% - 11rpx);
}

.goods-card {
  background-color: #ffffff;
  border-radius: 16rpx;
  overflow: hidden;
}

.goods-image {
  width: 100%;
  height: 340rpx;
}

.goods-info {
  padding: 16rpx 20rpx 20rpx;
}

.goods-name {
  font-size: 26rpx;
  color: #333333;
  line-height: 1.4;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  text-overflow: ellipsis;
}

.goods-price-row {
  display: flex;
  align-items: baseline;
  margin-top: 10rpx;
  gap: 10rpx;
}

.goods-price {
  font-size: 32rpx;
  font-weight: bold;
  color: #e64340;
}

.goods-original-price {
  font-size: 22rpx;
  color: #cccccc;
  text-decoration: line-through;
}

.goods-sales {
  font-size: 20rpx;
  color: #999999;
  margin-top: 6rpx;
}

/* --- Rich Text --- */
.richtext-wrapper {
  padding: 20rpx 24rpx;
}

/* --- Load More --- */
.load-more {
  padding: 30rpx 0 40rpx;
  text-align: center;
}

.load-more-text {
  font-size: 24rpx;
  color: #999999;
}
</style>
