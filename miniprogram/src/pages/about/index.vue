<template>
  <view class="about-page">
    <!-- Loading -->
    <view v-if="loading" class="loading-wrap">
      <text class="loading-text">{{ t('common.loading') }}</text>
    </view>

    <!-- Error -->
    <view v-else-if="error" class="error-wrap">
      <text class="error-text">{{ error }}</text>
      <view class="retry-btn" @click="loadAbout">
        <text class="retry-text">{{ t('common.retry') }}</text>
      </view>
    </view>

    <!-- Content -->
    <view v-else class="about-content">
      <!-- Logo & Company info -->
      <view class="company-section">
        <image class="company-logo" src="/static/logo.svg" mode="aspectFit" />
        <text class="company-name">{{ info.company_name || '国韵好运仓储' }}</text>
        <text class="company-desc">{{ info.description || '海外仓储 · 集运转运一站式服务' }}</text>
      </view>

      <!-- Version -->
      <view class="info-section">
        <view class="info-row">
          <text class="info-label">Version</text>
          <text class="info-value">{{ info.version || '1.0.0' }}</text>
        </view>
      </view>

      <!-- Contact info -->
      <view class="info-section" v-if="info.contacts && info.contacts.length > 0">
        <view class="section-title">
          <text class="section-title-text">Contact</text>
        </view>
        <view
          v-for="(contact, idx) in info.contacts"
          :key="idx"
          class="info-row"
          @click="onContactTap(contact)"
        >
          <text class="info-label">{{ contact.label }}</text>
          <text class="info-value clickable">{{ contact.value }}</text>
        </view>
      </view>

      <!-- Links -->
      <view class="link-section">
        <view class="link-item" @click="goPage('/pages/policy/privacy')">
          <text class="link-text">{{ t('policy.privacy') }}</text>
          <uni-icons type="right" size="14" color="#ccc" />
        </view>
        <view class="link-item" @click="goPage('/pages/policy/terms')">
          <text class="link-text">{{ t('policy.terms') }}</text>
          <uni-icons type="right" size="14" color="#ccc" />
        </view>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { t } from '@/locale'
import { commonApi } from '@/api/common'

interface ContactInfo {
  label: string
  value: string
  type?: string // phone, email, website
}

interface AboutInfo {
  company_name: string
  description: string
  version: string
  contacts: ContactInfo[]
}

const loading = ref(true)
const error = ref('')
const info = ref<AboutInfo>({
  company_name: '',
  description: '',
  version: '1.0.0',
  contacts: [],
})

async function loadAbout() {
  loading.value = true
  error.value = ''
  try {
    const res: any = await commonApi.getAbout()
    const data = res?.data || res
    info.value = {
      company_name: data.company_name || data.name || '国韵好运仓储',
      description: data.description || '',
      version: data.version || '1.0.0',
      contacts: data.contacts || [],
    }
  } catch (e: any) {
    error.value = e?.message || t('common.retry')
  } finally {
    loading.value = false
  }
}

function onContactTap(contact: ContactInfo) {
  if (contact.type === 'phone') {
    uni.makePhoneCall({ phoneNumber: contact.value })
  } else {
    uni.setClipboardData({
      data: contact.value,
      success: () => {
        uni.showToast({ title: t('common.copied'), icon: 'success' })
      },
    })
  }
}

function goPage(url: string) {
  uni.navigateTo({ url })
}

onMounted(() => {
  uni.setNavigationBarTitle({ title: t('about.title') })
  loadAbout()
})
</script>

<style scoped>
.about-page {
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

.company-section {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 60rpx 32rpx 40rpx;
  background-color: #fff;
}

.company-logo {
  width: 160rpx;
  height: 160rpx;
  margin-bottom: 24rpx;
}

.company-name {
  font-size: 36rpx;
  font-weight: 600;
  color: #0f3a57;
  margin-bottom: 16rpx;
}

.company-desc {
  font-size: 26rpx;
  color: #666;
  text-align: center;
  line-height: 1.6;
  padding: 0 40rpx;
}

.info-section {
  background-color: #fff;
  margin-top: 20rpx;
}

.section-title {
  padding: 24rpx 32rpx 8rpx;
}

.section-title-text {
  font-size: 28rpx;
  font-weight: 600;
  color: #0f3a57;
}

.info-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 28rpx 32rpx;
  border-bottom: 1rpx solid #f0f0f0;
}

.info-row:last-child {
  border-bottom: none;
}

.info-label {
  font-size: 28rpx;
  color: #333;
}

.info-value {
  font-size: 28rpx;
  color: #666;
}

.info-value.clickable {
  color: #0f3a57;
}

.link-section {
  background-color: #fff;
  margin-top: 20rpx;
}

.link-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 28rpx 32rpx;
  border-bottom: 1rpx solid #f0f0f0;
}

.link-item:last-child {
  border-bottom: none;
}

.link-text {
  font-size: 28rpx;
  color: #333;
}
</style>
