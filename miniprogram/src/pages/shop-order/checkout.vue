<template>
  <view class="page">
    <!-- Address -->
    <view class="addr-card" @tap="goPickAddress">
      <view v-if="address" class="addr-info">
        <view class="addr-row">
          <text class="addr-name">{{ address.recipient_name }}</text>
          <text class="addr-phone">{{ address.phone }}</text>
          <view v-if="address.is_default" class="addr-default">默認</view>
        </view>
        <text class="addr-detail">{{ fullAddress }}</text>
      </view>
      <view v-else class="addr-empty">
        <text>請選擇收貨地址</text>
      </view>
      <text class="addr-arrow">&#x276F;</text>
    </view>

    <!-- Goods -->
    <view class="goods-card">
      <view class="shop-header">
        <text>国韵好运商城</text>
      </view>
      <view v-for="item in items" :key="item.id" class="goods-row">
        <image class="goods-img" :src="item.image" mode="aspectFill" />
        <view class="goods-info">
          <text class="goods-name">{{ item.name }}</text>
          <text v-if="item.sku_name" class="goods-sku">{{ item.sku_name }}</text>
          <view class="goods-bottom">
            <text class="goods-price">¥{{ item.price.toFixed(2) }}</text>
            <text class="goods-qty">x{{ item.quantity }}</text>
          </view>
        </view>
      </view>
    </view>

    <!-- Remark -->
    <view class="remark-card">
      <text class="remark-label">訂單備註</text>
      <input class="remark-input" v-model="remark" placeholder="選填，請輸入備註" maxlength="100" />
    </view>

    <!-- Summary -->
    <view class="sum-card">
      <view class="sum-row">
        <text>商品金額</text>
        <text>¥{{ totalPrice.toFixed(2) }}</text>
      </view>
      <view class="sum-row">
        <text>運費</text>
        <text>¥0.00</text>
      </view>
    </view>

    <!-- Bottom bar -->
    <view class="bottom-bar">
      <view class="total-info">
        <text class="total-label">合計：</text>
        <text class="total-price">¥{{ totalPrice.toFixed(2) }}</text>
      </view>
      <view class="submit-btn" :class="{ disabled: submitting || !address }" @tap="submitOrder">
        <text class="submit-text">{{ submitting ? '提交中…' : '提交訂單' }}</text>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { onLoad, onShow } from '@dcloudio/uni-app'
import { commonApi } from '@/api/common'
import { orderApi } from '@/api/order'

interface Item {
  id: string
  goods_id: string
  sku_id?: string
  sku_name?: string
  name: string
  image: string
  price: number
  quantity: number
}

const cartIds = ref<string[]>([])
const items = ref<Item[]>([])
const address = ref<any>(null)
const remark = ref('')
const submitting = ref(false)

const totalPrice = computed(() =>
  items.value.reduce((s, i) => s + i.price * i.quantity, 0)
)

const fullAddress = computed(() => {
  if (!address.value) return ''
  const a = address.value
  return [a.province, a.city, a.district, a.address].filter(Boolean).join('')
})

onLoad((options: any) => {
  const ids = (options?.cart_ids || '').split(',').filter(Boolean)
  cartIds.value = ids
  loadCart(ids)
})

onShow(() => {
  // re-load default address each time (in case user picks new one)
  if (!address.value) loadDefaultAddress()
})

async function loadCart(ids: string[]) {
  try {
    const res: any = await commonApi.getCart()
    const data = res?.data ?? res
    const list = Array.isArray(data) ? data : (data?.list || [])
    items.value = list
      .filter((it: any) => ids.includes(it.id))
      .map((it: any) => ({
        id: it.id,
        goods_id: it.goods_id,
        sku_id: it.sku_id,
        sku_name: it.sku_name || '',
        name: it.name || it.goods_name || '商品',
        image:
          it.image ||
          it.goods_image_url ||
          it.goods_image ||
          'https://placehold.co/600x600/0f3a57/ffffff/png?text=GUOYUN&font=roboto',
        price: Number(it.price ?? it.goods_price ?? 0),
        quantity: Number(it.quantity) || 1,
      }))
    if (items.value.length === 0) {
      uni.showToast({ title: '購物車項目失效', icon: 'none' })
    }
  } catch (e) {
    uni.showToast({ title: '載入失敗', icon: 'none' })
  }
}

async function loadDefaultAddress() {
  try {
    const res: any = await commonApi.getAddresses()
    const data = res?.data ?? res
    const list = Array.isArray(data) ? data : (data?.list || [])
    address.value = list.find((a: any) => a.is_default) || list[0] || null
  } catch (e) {
    // ignore
  }
}

function goPickAddress() {
  uni.navigateTo({ url: '/pages/address/list?picker=1' })
}

async function submitOrder() {
  if (submitting.value) return
  if (!address.value) {
    uni.showToast({ title: '請選擇地址', icon: 'none' })
    return
  }
  if (items.value.length === 0) {
    uni.showToast({ title: '無商品', icon: 'none' })
    return
  }
  submitting.value = true
  try {
    const res: any = await orderApi.checkoutCart({
      cart_ids: cartIds.value,
      address_id: address.value.id,
      remark: remark.value,
      pay_method: 'wxpay',
    })
    const data = res?.data ?? res
    uni.showToast({ title: '下單成功', icon: 'success' })
    setTimeout(() => {
      uni.redirectTo({ url: `/pages/shop-order/detail?id=${data.id}` })
    }, 800)
  } catch (e: any) {
    uni.showToast({ title: e?.message || '下單失敗', icon: 'none' })
  } finally {
    submitting.value = false
  }
}
</script>

<style scoped>
.page {
  min-height: 100vh;
  background: #f5f5f5;
  padding: 20rpx 24rpx 140rpx;
}
.addr-card, .goods-card, .remark-card, .sum-card {
  background: #fff;
  border-radius: 16rpx;
  padding: 28rpx;
  margin-bottom: 20rpx;
  display: flex;
}
.addr-card { align-items: center; }
.addr-info { flex: 1; min-width: 0; display: flex; flex-direction: column; gap: 8rpx; }
.addr-row { display: flex; align-items: center; gap: 16rpx; }
.addr-name { font-size: 30rpx; font-weight: bold; color: #333; }
.addr-phone { font-size: 26rpx; color: #666; }
.addr-default { font-size: 20rpx; background: #fbbf24; color: #fff; padding: 2rpx 10rpx; border-radius: 6rpx; }
.addr-detail { font-size: 26rpx; color: #555; line-height: 1.5; }
.addr-empty { flex: 1; font-size: 28rpx; color: #999; }
.addr-arrow { color: #ccc; font-size: 28rpx; }

.goods-card { flex-direction: column; }
.shop-header { font-size: 28rpx; color: #333; padding-bottom: 16rpx; border-bottom: 1rpx solid #f0f0f0; }
.goods-row { display: flex; gap: 16rpx; padding: 20rpx 0; border-bottom: 1rpx solid #f5f5f5; }
.goods-row:last-child { border-bottom: none; }
.goods-img { width: 140rpx; height: 140rpx; border-radius: 12rpx; }
.goods-info { flex: 1; display: flex; flex-direction: column; justify-content: space-between; }
.goods-name { font-size: 26rpx; color: #333; }
.goods-sku { font-size: 22rpx; color: #999; background: #f5f5f5; padding: 4rpx 10rpx; border-radius: 4rpx; align-self: flex-start; }
.goods-bottom { display: flex; justify-content: space-between; align-items: center; }
.goods-price { font-size: 28rpx; color: #e64340; font-weight: bold; }
.goods-qty { font-size: 24rpx; color: #999; }

.remark-card { align-items: center; gap: 16rpx; }
.remark-label { font-size: 28rpx; color: #333; flex-shrink: 0; }
.remark-input { flex: 1; font-size: 26rpx; color: #333; }

.sum-card { flex-direction: column; gap: 12rpx; }
.sum-row { display: flex; justify-content: space-between; font-size: 26rpx; color: #555; }

.bottom-bar {
  position: fixed; bottom: 0; left: 0; right: 0;
  height: 110rpx; background: #fff;
  display: flex; align-items: center; justify-content: space-between;
  padding: 0 24rpx; box-shadow: 0 -2rpx 10rpx rgba(0,0,0,0.05);
  padding-bottom: env(safe-area-inset-bottom, 0);
}
.total-info { display: flex; align-items: baseline; }
.total-label { font-size: 26rpx; color: #333; }
.total-price { font-size: 36rpx; font-weight: bold; color: #e64340; }
.submit-btn { padding: 18rpx 50rpx; background: #0f3a57; border-radius: 44rpx; }
.submit-btn.disabled { opacity: 0.5; }
.submit-text { font-size: 28rpx; color: #fff; font-weight: bold; }
</style>
