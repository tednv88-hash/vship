<template>
  <view class="service-page">
    <!-- Header -->
    <view class="header-card">
      <uni-icons type="headphones" size="48" color="#fff" />
      <text class="header-title">{{ t('service.title') }}</text>
      <text class="header-desc">有任何問題，歡迎隨時聯繫我們</text>
    </view>

    <!-- Contact phone -->
    <view class="contact-card">
      <view class="contact-item" @click="callPhone">
        <view class="contact-icon">
          <uni-icons type="phone" size="24" color="#0f3a57" />
        </view>
        <view class="contact-info">
          <text class="contact-label">客服電話</text>
          <text class="contact-value">{{ phoneNumber }}</text>
        </view>
        <view class="call-btn">
          <text>撥打</text>
        </view>
      </view>
    </view>

    <!-- WeChat QR code -->
    <view class="contact-card">
      <view class="contact-item">
        <view class="contact-icon">
          <uni-icons type="chatbubble" size="24" color="#0f3a57" />
        </view>
        <view class="contact-info">
          <text class="contact-label">微信客服</text>
          <text class="contact-value">長按識別二維碼添加客服</text>
        </view>
      </view>
      <view class="qr-wrap">
        <image class="qr-image" src="/static/wechat-qr.png" mode="aspectFit" show-menu-by-longpress />
        <text class="qr-hint">長按保存或識別二維碼</text>
      </view>
    </view>

    <!-- Working hours -->
    <view class="contact-card">
      <view class="contact-item">
        <view class="contact-icon">
          <uni-icons type="calendar" size="24" color="#0f3a57" />
        </view>
        <view class="contact-info">
          <text class="contact-label">工作時間</text>
          <text class="contact-value">{{ workingHours }}</text>
        </view>
      </view>
    </view>

    <!-- FAQ quick links -->
    <view class="faq-card">
      <text class="faq-title">常見問題</text>
      <view v-for="faq in faqList" :key="faq.id" class="faq-item" @click="goFaq(faq.id)">
        <text class="faq-text">{{ faq.title }}</text>
        <uni-icons type="right" size="14" color="#999" />
      </view>
    </view>

    <!-- Online chat button -->
    <view class="chat-btn" @click="openChat">
      <uni-icons type="chatbubble-filled" size="24" color="#fff" />
      <text>線上客服</text>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { t } from '@/locale'
import { commonApi } from '@/api/common'

const phoneNumber = ref('400-888-8888')
const workingHours = ref('週一至週五 9:00 - 18:00')
const faqList = ref<any[]>([])

onMounted(() => {
  loadServiceInfo()
  loadFaq()
})

async function loadServiceInfo() {
  try {
    const res = await commonApi.getAppSettings('customer_service')
    const data = res?.data?.[0]?.config || res?.data?.config || {}
    if (data.phone) phoneNumber.value = data.phone
    if (data.working_hours) workingHours.value = data.working_hours
  } catch (e) {
    console.error(e)
  }
}

async function loadFaq() {
  try {
    const res = await commonApi.getHelpList({ page: 1, limit: 8, is_faq: true })
    faqList.value = res?.data?.list || res?.data || []
  } catch (e) {
    console.error(e)
  }
}

function callPhone() {
  uni.makePhoneCall({ phoneNumber: phoneNumber.value })
}

function goFaq(id: string) {
  uni.navigateTo({ url: `/pages/help/detail?id=${id}` })
}

function openChat() {
  // Open built-in customer service if available
  // #ifdef MP-WEIXIN
  // WeChat miniprogram has built-in contact feature via <button open-type="contact">
  // #endif
  uni.showToast({ title: '正在連接客服...', icon: 'none' })
}
</script>

<style scoped>
.service-page {
  min-height: 100vh;
  background: #f5f5f5;
  padding-bottom: 120rpx;
}

.header-card {
  background: linear-gradient(135deg, #0f3a57 0%, #1a5c7a 100%);
  padding: 60rpx 30rpx;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 16rpx;
}

.header-title {
  font-size: 36rpx;
  font-weight: 600;
  color: #fff;
}

.header-desc {
  font-size: 26rpx;
  color: rgba(255, 255, 255, 0.7);
}

.contact-card {
  background: #fff;
  margin: 16rpx 24rpx 0;
  border-radius: 16rpx;
  padding: 24rpx 30rpx;
}

.contact-item {
  display: flex;
  align-items: center;
  gap: 20rpx;
}

.contact-icon {
  width: 60rpx;
  height: 60rpx;
  background: rgba(15, 58, 87, 0.08);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.contact-info {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 6rpx;
}

.contact-label {
  font-size: 28rpx;
  color: #333;
  font-weight: 500;
}

.contact-value {
  font-size: 24rpx;
  color: #999;
}

.call-btn {
  background: #0f3a57;
  border-radius: 24rpx;
  padding: 10rpx 28rpx;
}

.call-btn text {
  font-size: 24rpx;
  color: #fff;
}

.qr-wrap {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 30rpx 0 10rpx;
}

.qr-image {
  width: 300rpx;
  height: 300rpx;
}

.qr-hint {
  font-size: 22rpx;
  color: #999;
  margin-top: 16rpx;
}

.faq-card {
  background: #fff;
  margin: 16rpx 24rpx 0;
  border-radius: 16rpx;
  padding: 24rpx 30rpx;
}

.faq-title {
  font-size: 30rpx;
  font-weight: 600;
  color: #333;
  margin-bottom: 16rpx;
  display: block;
}

.faq-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20rpx 0;
  border-bottom: 1rpx solid #f5f5f5;
}

.faq-item:last-child {
  border-bottom: none;
}

.faq-text {
  font-size: 26rpx;
  color: #666;
  flex: 1;
}

.chat-btn {
  position: fixed;
  bottom: 60rpx;
  right: 40rpx;
  background: #0f3a57;
  border-radius: 44rpx;
  padding: 20rpx 36rpx;
  display: flex;
  align-items: center;
  gap: 10rpx;
  box-shadow: 0 8rpx 24rpx rgba(15, 58, 87, 0.3);
  z-index: 10;
}

.chat-btn text {
  color: #fff;
  font-size: 28rpx;
  font-weight: 500;
}
</style>
