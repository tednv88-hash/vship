<template>
  <view class="page">
    <!-- Cart List -->
    <scroll-view v-if="cartList.length > 0" class="cart-scroll" scroll-y>
      <view class="cart-list">
        <view
          v-for="item in cartList"
          :key="item.id"
          class="cart-item"
        >
          <!-- Swipe delete layer (simplified: show delete on long press or button) -->
          <view class="cart-item-inner">
            <!-- Checkbox -->
            <view class="cart-checkbox" @tap="toggleItem(item)">
              <view class="checkbox-icon" :class="{ checked: item.checked }">
                <text v-if="item.checked" class="check-mark">&#x2713;</text>
              </view>
            </view>

            <!-- Product Image -->
            <image
              class="cart-image"
              :src="item.image || 'https://placehold.co/600x600/0f3a57/ffffff/png?text=GUOYUN&font=roboto'"
              mode="aspectFill"
              @tap="goToGoodsDetail(item.goods_id)"
            />

            <!-- Product Info -->
            <view class="cart-info">
              <text class="cart-name" @tap="goToGoodsDetail(item.goods_id)">{{ item.name }}</text>
              <text v-if="item.sku_name" class="cart-sku">{{ item.sku_name }}</text>
              <view class="cart-bottom-row">
                <text class="cart-price">¥{{ item.price }}</text>
                <!-- Quantity Stepper -->
                <view class="stepper">
                  <view
                    class="stepper-btn"
                    :class="{ disabled: item.quantity <= 1 }"
                    @tap="changeQuantity(item, -1)"
                  >
                    <text class="stepper-btn-text">-</text>
                  </view>
                  <text class="stepper-value">{{ item.quantity }}</text>
                  <view
                    class="stepper-btn"
                    @tap="changeQuantity(item, 1)"
                  >
                    <text class="stepper-btn-text">+</text>
                  </view>
                </view>
              </view>
            </view>

            <!-- Delete button -->
            <view class="cart-delete" @tap="deleteItem(item)">
              <text class="delete-text">{{ t('common.delete') }}</text>
            </view>
          </view>
        </view>
      </view>
    </scroll-view>

    <!-- Empty Cart State -->
    <view v-else-if="!loading" class="empty-cart">
      <text class="empty-icon">&#x1F6D2;</text>
      <text class="empty-title">購物車是空的</text>
      <text class="empty-desc">去逛逛吧，發現更多好物</text>
      <view class="empty-btn" @tap="goShopping">
        <text class="empty-btn-text">去逛逛</text>
      </view>
    </view>

    <!-- Loading State -->
    <view v-else class="loading-state">
      <text class="loading-text">{{ t('common.loading') }}</text>
    </view>

    <!-- Bottom Bar -->
    <view v-if="cartList.length > 0" class="bottom-bar">
      <view class="bottom-left">
        <view class="select-all" @tap="toggleSelectAll">
          <view class="checkbox-icon" :class="{ checked: isAllSelected }">
            <text v-if="isAllSelected" class="check-mark">&#x2713;</text>
          </view>
          <text class="select-all-text">全選</text>
        </view>
      </view>
      <view class="bottom-right">
        <view class="total-info">
          <text class="total-label">合計：</text>
          <text class="total-price">¥{{ totalPrice }}</text>
        </view>
        <view
          class="checkout-btn"
          :class="{ disabled: selectedCount === 0 }"
          @tap="onCheckout"
        >
          <text class="checkout-text">結算({{ selectedCount }})</text>
        </view>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { onShow } from '@dcloudio/uni-app'
import { t } from '@/locale'
import { commonApi } from '@/api/common'
import { setCartCount } from '@/store'

interface CartItem {
  id: string
  goods_id: string
  name: string
  image: string
  price: number
  quantity: number
  sku_id?: string
  sku_name?: string
  checked: boolean
}

const loading = ref(true)
const cartList = ref<CartItem[]>([])

const selectedCount = computed(() => {
  return cartList.value.filter(item => item.checked).length
})

const isAllSelected = computed(() => {
  return cartList.value.length > 0 && cartList.value.every(item => item.checked)
})

const totalPrice = computed(() => {
  const total = cartList.value
    .filter(item => item.checked)
    .reduce((sum, item) => sum + item.price * item.quantity, 0)
  return total.toFixed(2)
})

onShow(() => {
  fetchCart()
})

async function fetchCart() {
  loading.value = true
  try {
    const res: any = await commonApi.getCart()
    const data = res?.data ?? res
    const list = Array.isArray(data) ? data : (data?.list || data?.items || [])
    cartList.value = list.map((item: any) => ({
      id: item.id,
      goods_id: item.goods_id,
      sku_id: item.sku_id,
      sku_name: item.sku_name || item.goods_sku_name || '',
      name: item.name || item.goods_name || '商品',
      image:
        item.image ||
        item.goods_image_url ||
        item.goods_image ||
        'https://placehold.co/600x600/0f3a57/ffffff/png?text=GUOYUN&font=roboto',
      price: Number(item.price ?? item.goods_price ?? 0),
      quantity: Number(item.quantity) || 1,
      checked: true,
    }))
  } catch (e) {
    cartList.value = []
  } finally {
    loading.value = false
  }
  updateCartBadge()
}

function updateCartBadge() {
  const total = cartList.value.reduce((sum, item) => sum + item.quantity, 0)
  setCartCount(total)
}

function toggleItem(item: CartItem) {
  item.checked = !item.checked
}

function toggleSelectAll() {
  const newVal = !isAllSelected.value
  cartList.value.forEach(item => {
    item.checked = newVal
  })
}

async function changeQuantity(item: CartItem, delta: number) {
  const newQty = item.quantity + delta
  if (newQty < 1) return

  item.quantity = newQty
  try {
    await commonApi.updateCartItem(item.id, { quantity: newQty })
  } catch (e) {
    item.quantity = newQty - delta
    uni.showToast({ title: '更新失敗', icon: 'none' })
  }
  updateCartBadge()
}

async function deleteItem(item: CartItem) {
  uni.showModal({
    title: '提示',
    content: '確認刪除該商品嗎？',
    success: async (res) => {
      if (res.confirm) {
        try {
          await commonApi.deleteCartItem(item.id)
          cartList.value = cartList.value.filter(i => i.id !== item.id)
          updateCartBadge()
          uni.showToast({ title: '已刪除', icon: 'success' })
        } catch (e) {
          uni.showToast({ title: '刪除失敗', icon: 'none' })
        }
      }
    },
  })
}

function onCheckout() {
  if (selectedCount.value === 0) {
    uni.showToast({ title: '請選擇商品', icon: 'none' })
    return
  }

  const selectedIds = cartList.value
    .filter(item => item.checked)
    .map(item => item.id)
    .join(',')

  uni.navigateTo({ url: `/pages/shop-order/checkout?cart_ids=${selectedIds}` })
}

function goToGoodsDetail(id: string) {
  uni.navigateTo({ url: `/pages/goods/detail?id=${id}` })
}

function goShopping() {
  uni.switchTab({ url: '/pages/index/index' })
}
</script>

<style scoped>
.page {
  min-height: 100vh;
  background-color: #f5f5f5;
  padding-bottom: 120rpx;
}

/* --- Cart Scroll --- */
.cart-scroll {
  height: calc(100vh - 120rpx);
}

.cart-list {
  padding: 20rpx 24rpx;
}

/* --- Cart Item --- */
.cart-item {
  margin-bottom: 20rpx;
  background-color: #ffffff;
  border-radius: 16rpx;
  overflow: hidden;
}

.cart-item-inner {
  display: flex;
  align-items: center;
  padding: 24rpx;
}

/* --- Checkbox --- */
.cart-checkbox {
  flex-shrink: 0;
  padding-right: 20rpx;
}

.checkbox-icon {
  width: 40rpx;
  height: 40rpx;
  border-radius: 50%;
  border: 3rpx solid #cccccc;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: #ffffff;
}

.checkbox-icon.checked {
  background-color: #0f3a57;
  border-color: #0f3a57;
}

.check-mark {
  font-size: 24rpx;
  color: #ffffff;
  font-weight: bold;
}

/* --- Image --- */
.cart-image {
  width: 180rpx;
  height: 180rpx;
  border-radius: 12rpx;
  flex-shrink: 0;
  margin-right: 20rpx;
}

/* --- Info --- */
.cart-info {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  height: 180rpx;
}

.cart-name {
  font-size: 26rpx;
  color: #333333;
  line-height: 1.4;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  text-overflow: ellipsis;
}

.cart-sku {
  font-size: 22rpx;
  color: #999999;
  margin-top: 6rpx;
  background-color: #f5f5f5;
  padding: 4rpx 12rpx;
  border-radius: 4rpx;
  align-self: flex-start;
}

.cart-bottom-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-top: 8rpx;
}

.cart-price {
  font-size: 30rpx;
  font-weight: bold;
  color: #e64340;
}

/* --- Stepper --- */
.stepper {
  display: flex;
  align-items: center;
  border: 1rpx solid #e0e0e0;
  border-radius: 8rpx;
  overflow: hidden;
}

.stepper-btn {
  width: 52rpx;
  height: 48rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: #f8f8f8;
}

.stepper-btn.disabled {
  opacity: 0.4;
}

.stepper-btn-text {
  font-size: 28rpx;
  color: #333333;
  font-weight: bold;
}

.stepper-value {
  min-width: 60rpx;
  height: 48rpx;
  line-height: 48rpx;
  text-align: center;
  font-size: 26rpx;
  color: #333333;
  background-color: #ffffff;
  border-left: 1rpx solid #e0e0e0;
  border-right: 1rpx solid #e0e0e0;
}

/* --- Delete --- */
.cart-delete {
  flex-shrink: 0;
  margin-left: 16rpx;
  padding: 12rpx 16rpx;
}

.delete-text {
  font-size: 22rpx;
  color: #e64340;
}

/* --- Empty Cart --- */
.empty-cart {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding-top: 240rpx;
}

.empty-icon {
  font-size: 120rpx;
  margin-bottom: 30rpx;
}

.empty-title {
  font-size: 32rpx;
  color: #333333;
  font-weight: bold;
  margin-bottom: 12rpx;
}

.empty-desc {
  font-size: 26rpx;
  color: #999999;
  margin-bottom: 40rpx;
}

.empty-btn {
  padding: 20rpx 80rpx;
  background-color: #0f3a57;
  border-radius: 44rpx;
}

.empty-btn-text {
  font-size: 28rpx;
  color: #ffffff;
}

/* --- Loading --- */
.loading-state {
  display: flex;
  align-items: center;
  justify-content: center;
  padding-top: 240rpx;
}

.loading-text {
  font-size: 28rpx;
  color: #999999;
}

/* --- Bottom Bar --- */
.bottom-bar {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  height: 110rpx;
  background-color: #ffffff;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 24rpx;
  box-shadow: 0 -2rpx 10rpx rgba(0, 0, 0, 0.05);
  z-index: 50;
  /* safe area for iPhone */
  padding-bottom: env(safe-area-inset-bottom, 0);
}

.bottom-left {
  display: flex;
  align-items: center;
}

.select-all {
  display: flex;
  align-items: center;
  gap: 12rpx;
}

.select-all-text {
  font-size: 26rpx;
  color: #333333;
}

.bottom-right {
  display: flex;
  align-items: center;
  gap: 20rpx;
}

.total-info {
  display: flex;
  align-items: baseline;
}

.total-label {
  font-size: 26rpx;
  color: #333333;
}

.total-price {
  font-size: 36rpx;
  font-weight: bold;
  color: #e64340;
}

.checkout-btn {
  padding: 16rpx 40rpx;
  background-color: #0f3a57;
  border-radius: 44rpx;
}

.checkout-btn.disabled {
  opacity: 0.5;
}

.checkout-text {
  font-size: 28rpx;
  color: #ffffff;
  font-weight: bold;
}
</style>
