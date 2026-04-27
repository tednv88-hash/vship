<template>
  <view class="goods-detail">
    <!-- Image swiper -->
    <swiper
      class="gallery-swiper"
      :indicator-dots="true"
      :autoplay="true"
      :interval="3000"
      indicator-active-color="#0f3a57"
    >
      <swiper-item v-for="(img, idx) in detail.images" :key="idx">
        <image class="gallery-img" :src="img" mode="aspectFill" @click="previewImage(idx)" />
      </swiper-item>
    </swiper>

    <!-- Product info -->
    <view class="product-info">
      <view class="price-row">
        <text class="price">¥{{ selectedSku?.price || detail.price }}</text>
        <text v-if="detail.original_price" class="original-price">¥{{ detail.original_price }}</text>
      </view>
      <text class="product-name">{{ detail.name }}</text>
      <view class="meta-row">
        <text class="meta-item">{{ t('goods.sales') }} {{ detail.sales || 0 }}</text>
        <text class="meta-item">{{ t('goods.stock') }} {{ selectedSku?.stock || detail.stock || 0 }}</text>
      </view>
    </view>

    <!-- SKU selector -->
    <view v-if="detail.skus && detail.skus.length > 0" class="section sku-section" @click="showSkuPopup = true">
      <text class="section-label">規格選擇</text>
      <view class="sku-selected">
        <text>{{ selectedSku ? selectedSku.name : '請選擇規格' }}</text>
        <uni-icons type="right" size="16" color="#999" />
      </view>
    </view>

    <!-- Quantity stepper -->
    <view class="section quantity-section">
      <text class="section-label">數量</text>
      <view class="stepper">
        <view class="stepper-btn" @click="changeQuantity(-1)">
          <text>-</text>
        </view>
        <text class="stepper-value">{{ quantity }}</text>
        <view class="stepper-btn" @click="changeQuantity(1)">
          <text>+</text>
        </view>
      </view>
    </view>

    <!-- Product description -->
    <view class="section">
      <text class="section-title">商品詳情</text>
      <rich-text class="rich-content" :nodes="detail.description || ''" />
    </view>

    <!-- Reviews section -->
    <view class="section">
      <view class="section-header">
        <text class="section-title">商品評價 ({{ reviewCount }})</text>
        <text class="section-more" @click="goReviews">{{ t('common.more') }} ></text>
      </view>
      <view v-if="reviews.length > 0" class="reviews-list">
        <view v-for="review in reviews" :key="review.id" class="review-item">
          <view class="review-header">
            <image class="reviewer-avatar" :src="review.avatar || '/static/default-avatar.png'" />
            <text class="reviewer-name">{{ review.nickname }}</text>
            <text class="review-date">{{ review.created_at }}</text>
          </view>
          <text class="review-content">{{ review.content }}</text>
          <view v-if="review.images && review.images.length" class="review-images">
            <image
              v-for="(img, i) in review.images"
              :key="i"
              class="review-img"
              :src="img"
              mode="aspectFill"
            />
          </view>
        </view>
      </view>
      <view v-else class="no-review">
        <text>暫無評價</text>
      </view>
    </view>

    <!-- Bottom padding for fixed bar -->
    <view style="height: 120rpx" />

    <!-- Bottom action bar -->
    <view class="bottom-bar">
      <view class="bar-icon" @click="contactService">
        <uni-icons type="headphones" size="22" color="#666" />
        <text class="bar-icon-text">客服</text>
      </view>
      <view class="bar-btn cart-btn" @click="addToCart">
        <text>{{ t('goods.addToCart') }}</text>
      </view>
      <view class="bar-btn buy-btn" @click="buyNow">
        <text>{{ t('goods.buyNow') }}</text>
      </view>
    </view>

    <!-- SKU popup -->
    <uni-popup ref="skuPopupRef" type="bottom" :is-mask-click="true">
      <view v-if="showSkuPopup" class="sku-popup">
        <view class="sku-popup-header">
          <image class="sku-popup-img" :src="selectedSku?.image || detail.images?.[0]" mode="aspectFill" />
          <view class="sku-popup-info">
            <text class="price">¥{{ selectedSku?.price || detail.price }}</text>
            <text class="sku-stock">{{ t('goods.stock') }}: {{ selectedSku?.stock || detail.stock || 0 }}</text>
          </view>
          <uni-icons type="closeempty" size="22" color="#999" @click="closeSkuPopup" />
        </view>
        <scroll-view scroll-y class="sku-options-scroll">
          <view v-for="group in detail.sku_groups" :key="group.name" class="sku-group">
            <text class="sku-group-title">{{ group.name }}</text>
            <view class="sku-tags">
              <view
                v-for="option in group.options"
                :key="option"
                class="sku-tag"
                :class="{ selected: selectedSkuValues[group.name] === option }"
                @click="selectSku(group.name, option)"
              >
                <text>{{ option }}</text>
              </view>
            </view>
          </view>
        </scroll-view>
        <view class="sku-popup-footer">
          <view class="bar-btn buy-btn" style="flex: 1" @click="confirmSku">
            <text>{{ t('common.confirm') }}</text>
          </view>
        </view>
      </view>
    </uni-popup>
  </view>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import { onLoad } from '@dcloudio/uni-app'
import { t } from '@/locale'
import { goodsApi } from '@/api/goods'
import { commonApi } from '@/api/common'

const detail = ref<any>({})
const reviews = ref<any[]>([])
const reviewCount = ref(0)
const quantity = ref(1)
const showSkuPopup = ref(false)
const skuPopupRef = ref()
const selectedSkuValues = reactive<Record<string, string>>({})

const selectedSku = computed(() => {
  if (!detail.value.skus) return null
  return detail.value.skus.find((sku: any) => {
    return Object.entries(selectedSkuValues).every(
      ([key, val]) => sku.specs?.[key] === val
    )
  }) || null
})

const goodsId = ref('')

onLoad((options: any) => {
  const id = options?.id || options?.goods_id || ''
  if (id) {
    goodsId.value = id
    loadDetail(id)
    loadReviews(id)
  } else {
    uni.showToast({ title: '商品 ID 缺失', icon: 'none' })
  }
})

onMounted(() => {
  // Fallback for environments where onLoad options miss
  if (!goodsId.value) {
    const pages = getCurrentPages()
    const currentPage = pages[pages.length - 1] as any
    const opts = currentPage?.$page?.options || currentPage?.options || {}
    if (opts.id) {
      goodsId.value = opts.id
      loadDetail(opts.id)
      loadReviews(opts.id)
    }
  }
})

async function loadDetail(id: string) {
  try {
    const res = await goodsApi.getDetail(id)
    const d: any = res?.data || res || {}
    // Normalize: ensure images array, fallback sales field
    if (!d.images || (Array.isArray(d.images) && d.images.length === 0)) {
      d.images = d.image_url ? [d.image_url] : []
    }
    if (!d.images.length) {
      d.images = ['https://placehold.co/600x600/0f3a57/ffffff/png?text=GUOYUN&font=roboto']
    }
    if (d.sales == null) d.sales = d.sales_count || 0
    detail.value = d
    uni.setNavigationBarTitle({ title: d.name || t('goods.detail') })
  } catch (e) {
    console.error(e)
    uni.showToast({ title: '加載失敗', icon: 'none' })
  }
}

async function loadReviews(goodsId: string) {
  try {
    const res = await commonApi.getReviews({ goods_id: goodsId, page: 1, limit: 3 })
    const data = res?.data || {}
    reviews.value = data.list || data || []
    reviewCount.value = data.total || reviews.value.length
  } catch (e) {
    console.error(e)
  }
}

function changeQuantity(delta: number) {
  const next = quantity.value + delta
  if (next < 1) return
  const maxStock = selectedSku.value?.stock || detail.value.stock || 999
  if (next > maxStock) return
  quantity.value = next
}

function previewImage(index: number) {
  uni.previewImage({
    urls: detail.value.images || [],
    current: index,
  })
}

async function addToCart() {
  try {
    await commonApi.addToCart({
      goods_id: detail.value.id,
      sku_id: selectedSku.value?.id,
      quantity: quantity.value,
    })
    uni.showToast({ title: '已加入購物車', icon: 'success' })
  } catch (e) {
    console.error(e)
  }
}

async function buyNow() {
  try {
    uni.showLoading({ title: '處理中', mask: true })
    const res: any = await commonApi.addToCart({
      goods_id: detail.value.id,
      sku_id: selectedSku.value?.id,
      quantity: quantity.value,
    })
    const data = res?.data ?? res
    const cartId = data?.id
    uni.hideLoading()
    if (!cartId) {
      uni.showToast({ title: '下單失敗', icon: 'none' })
      return
    }
    uni.navigateTo({ url: `/pages/shop-order/checkout?cart_ids=${cartId}` })
  } catch (e) {
    uni.hideLoading()
    uni.showToast({ title: '下單失敗', icon: 'none' })
  }
}

function contactService() {
  uni.navigateTo({ url: '/pages/service/index' })
}

function goReviews() {
  uni.navigateTo({ url: `/pages/review/list?goods_id=${detail.value.id}` })
}

function selectSku(groupName: string, option: string) {
  selectedSkuValues[groupName] = option
}

function closeSkuPopup() {
  showSkuPopup.value = false
  skuPopupRef.value?.close()
}

function confirmSku() {
  closeSkuPopup()
}
</script>

<style scoped>
.goods-detail {
  background: #f5f5f5;
  min-height: 100vh;
}

.gallery-swiper {
  width: 100%;
  height: 750rpx;
}

.gallery-img {
  width: 100%;
  height: 100%;
}

.product-info {
  background: #fff;
  padding: 24rpx 30rpx;
  margin-bottom: 16rpx;
}

.price-row {
  display: flex;
  align-items: baseline;
  gap: 16rpx;
}

.price {
  font-size: 40rpx;
  color: #e64340;
  font-weight: 700;
}

.original-price {
  font-size: 26rpx;
  color: #999;
  text-decoration: line-through;
}

.product-name {
  display: block;
  font-size: 30rpx;
  color: #333;
  margin-top: 16rpx;
  line-height: 1.5;
}

.meta-row {
  display: flex;
  gap: 30rpx;
  margin-top: 16rpx;
}

.meta-item {
  font-size: 24rpx;
  color: #999;
}

.section {
  background: #fff;
  padding: 24rpx 30rpx;
  margin-bottom: 16rpx;
}

.section-label {
  font-size: 28rpx;
  color: #333;
}

.section-title {
  font-size: 30rpx;
  color: #333;
  font-weight: 600;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20rpx;
}

.section-more {
  font-size: 24rpx;
  color: #0f3a57;
}

.sku-section {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.sku-selected {
  display: flex;
  align-items: center;
  gap: 8rpx;
  font-size: 26rpx;
  color: #666;
}

.quantity-section {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.stepper {
  display: flex;
  align-items: center;
  border: 1rpx solid #ddd;
  border-radius: 8rpx;
  overflow: hidden;
}

.stepper-btn {
  width: 60rpx;
  height: 52rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #f5f5f5;
}

.stepper-btn text {
  font-size: 32rpx;
  color: #333;
}

.stepper-value {
  width: 80rpx;
  text-align: center;
  font-size: 28rpx;
}

.rich-content {
  margin-top: 20rpx;
}

/* Reviews */
.reviews-list {
  margin-top: 10rpx;
}

.review-item {
  padding: 20rpx 0;
  border-bottom: 1rpx solid #f0f0f0;
}

.review-item:last-child {
  border-bottom: none;
}

.review-header {
  display: flex;
  align-items: center;
  gap: 12rpx;
  margin-bottom: 12rpx;
}

.reviewer-avatar {
  width: 48rpx;
  height: 48rpx;
  border-radius: 50%;
}

.reviewer-name {
  font-size: 24rpx;
  color: #333;
  flex: 1;
}

.review-date {
  font-size: 22rpx;
  color: #999;
}

.review-content {
  font-size: 26rpx;
  color: #333;
  line-height: 1.5;
}

.review-images {
  display: flex;
  gap: 12rpx;
  margin-top: 12rpx;
}

.review-img {
  width: 160rpx;
  height: 160rpx;
  border-radius: 8rpx;
}

.no-review {
  text-align: center;
  padding: 40rpx 0;
  color: #999;
  font-size: 26rpx;
}

/* Bottom bar */
.bottom-bar {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  display: flex;
  align-items: center;
  background: #fff;
  padding: 16rpx 24rpx;
  padding-bottom: calc(16rpx + env(safe-area-inset-bottom));
  box-shadow: 0 -2rpx 10rpx rgba(0, 0, 0, 0.05);
  z-index: 100;
}

.bar-icon {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 0 24rpx;
}

.bar-icon-text {
  font-size: 20rpx;
  color: #666;
  margin-top: 4rpx;
}

.bar-btn {
  flex: 1;
  height: 80rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 40rpx;
  margin-left: 16rpx;
}

.bar-btn text {
  font-size: 28rpx;
  color: #fff;
  font-weight: 500;
}

.cart-btn {
  background: #ff9800;
}

.buy-btn {
  background: #0f3a57;
}

/* SKU popup */
.sku-popup {
  background: #fff;
  border-radius: 24rpx 24rpx 0 0;
  max-height: 70vh;
  padding: 30rpx;
  padding-bottom: calc(30rpx + env(safe-area-inset-bottom));
}

.sku-popup-header {
  display: flex;
  gap: 20rpx;
  padding-bottom: 24rpx;
  border-bottom: 1rpx solid #f0f0f0;
}

.sku-popup-img {
  width: 160rpx;
  height: 160rpx;
  border-radius: 12rpx;
}

.sku-popup-info {
  flex: 1;
  display: flex;
  flex-direction: column;
  justify-content: flex-end;
  gap: 12rpx;
}

.sku-stock {
  font-size: 24rpx;
  color: #999;
}

.sku-options-scroll {
  max-height: 40vh;
  padding: 20rpx 0;
}

.sku-group {
  margin-bottom: 24rpx;
}

.sku-group-title {
  font-size: 28rpx;
  color: #333;
  font-weight: 500;
  margin-bottom: 16rpx;
}

.sku-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 16rpx;
}

.sku-tag {
  padding: 12rpx 28rpx;
  background: #f5f5f5;
  border-radius: 8rpx;
  border: 2rpx solid transparent;
}

.sku-tag.selected {
  border-color: #0f3a57;
  background: rgba(15, 58, 87, 0.08);
  color: #0f3a57;
}

.sku-tag text {
  font-size: 26rpx;
}

.sku-popup-footer {
  display: flex;
  padding-top: 20rpx;
  border-top: 1rpx solid #f0f0f0;
}
</style>
