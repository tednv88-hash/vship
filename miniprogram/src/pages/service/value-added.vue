<template>
  <view class="vas-page">
    <view class="page-desc">
      <text>為您的包裹提供更優質的服務保障</text>
    </view>

    <!-- Service cards -->
    <view class="service-list">
      <view v-for="service in serviceList" :key="service.id" class="service-card">
        <view class="service-header">
          <view class="service-icon-wrap">
            <uni-icons :type="service.icon || 'star'" size="28" color="#0f3a57" />
          </view>
          <view class="service-title-wrap">
            <text class="service-name">{{ service.name }}</text>
            <text class="service-price">¥{{ service.price }}</text>
          </view>
        </view>
        <text class="service-desc">{{ service.description }}</text>
        <view class="service-footer">
          <view
            class="select-btn"
            :class="{ selected: selectedServices.includes(service.id) }"
            @click="toggleService(service.id)"
          >
            <text>{{ selectedServices.includes(service.id) ? '已選擇' : '選擇服務' }}</text>
          </view>
        </view>
      </view>
    </view>

    <!-- Empty -->
    <view v-if="!loading && serviceList.length === 0" class="empty">
      <text>{{ t('common.noData') }}</text>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { t } from '@/locale'
import { commonApi } from '@/api/common'

const serviceList = ref<any[]>([])
const selectedServices = ref<string[]>([])
const loading = ref(true)

const defaultServices = [
  {
    id: '1',
    name: '加固包裝',
    description: '使用加厚氣泡膜及硬質紙箱進行二次加固，有效防止運輸途中碰撞損壞，適合易碎物品。',
    price: '15.00',
    icon: 'shield',
  },
  {
    id: '2',
    name: '照片確認',
    description: '倉庫收貨後拍攝包裹內物品照片，確認商品數量、外觀狀態，讓您放心。',
    price: '5.00',
    icon: 'camera',
  },
  {
    id: '3',
    name: '合箱打包',
    description: '將多個包裹合併為一個包裹發出，節省運費。按實際重量計費。',
    price: '10.00',
    icon: 'inbox',
  },
  {
    id: '4',
    name: '代購服務',
    description: '提供專業代購服務，幫您購買當地商品並寄送到倉庫。代購費用另計。',
    price: '30.00',
    icon: 'cart',
  },
  {
    id: '5',
    name: '去除包裝',
    description: '去除商品原包裝盒，僅保留產品本體包裝，減少體積重量，節省運費。',
    price: '5.00',
    icon: 'minus',
  },
  {
    id: '6',
    name: '貼標換標',
    description: '為包裹更換或加貼物流標籤，方便後續配送追蹤。',
    price: '3.00',
    icon: 'flag',
  },
]

onMounted(() => {
  loadServices()
})

async function loadServices() {
  loading.value = true
  try {
    const res = await commonApi.getValueAddedServices()
    const list = res?.data?.list || res?.data || []
    serviceList.value = list.length > 0 ? list : defaultServices
  } catch (e) {
    serviceList.value = defaultServices
    console.error(e)
  } finally {
    loading.value = false
  }
}

function toggleService(id: string) {
  const idx = selectedServices.value.indexOf(id)
  if (idx > -1) {
    selectedServices.value.splice(idx, 1)
  } else {
    selectedServices.value.push(id)
  }
}
</script>

<style scoped>
.vas-page {
  min-height: 100vh;
  background: #f5f5f5;
  padding: 24rpx;
}

.page-desc {
  text-align: center;
  padding: 20rpx 0 30rpx;
}

.page-desc text {
  font-size: 26rpx;
  color: #999;
}

.service-list {
  display: flex;
  flex-direction: column;
  gap: 20rpx;
}

.service-card {
  background: #fff;
  border-radius: 16rpx;
  padding: 28rpx 30rpx;
}

.service-header {
  display: flex;
  align-items: center;
  gap: 20rpx;
  margin-bottom: 16rpx;
}

.service-icon-wrap {
  width: 64rpx;
  height: 64rpx;
  background: rgba(15, 58, 87, 0.08);
  border-radius: 16rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.service-title-wrap {
  flex: 1;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.service-name {
  font-size: 30rpx;
  color: #333;
  font-weight: 600;
}

.service-price {
  font-size: 30rpx;
  color: #e64340;
  font-weight: 600;
}

.service-desc {
  font-size: 26rpx;
  color: #666;
  line-height: 1.6;
  display: block;
}

.service-footer {
  display: flex;
  justify-content: flex-end;
  margin-top: 20rpx;
}

.select-btn {
  padding: 12rpx 32rpx;
  border: 2rpx solid #0f3a57;
  border-radius: 24rpx;
}

.select-btn text {
  font-size: 24rpx;
  color: #0f3a57;
}

.select-btn.selected {
  background: #0f3a57;
}

.select-btn.selected text {
  color: #fff;
}

.empty {
  text-align: center;
  padding: 100rpx 0;
  color: #999;
  font-size: 28rpx;
}
</style>
