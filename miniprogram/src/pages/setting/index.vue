<template>
  <view class="setting-page">
    <!-- Menu list -->
    <view class="menu-group">
      <view class="menu-item" @tap="navigateTo('/pages/user/index')">
        <text class="menu-label">個人資料</text>
        <text class="menu-arrow">›</text>
      </view>
      <view class="menu-item" @tap="showLangPicker = true">
        <text class="menu-label">語言設置</text>
        <view class="menu-right">
          <text class="menu-value">{{ currentLangLabel }}</text>
          <text class="menu-arrow">›</text>
        </view>
      </view>
    </view>

    <view class="menu-group">
      <view class="menu-item" @tap="handleClearCache">
        <text class="menu-label">清除緩存</text>
        <view class="menu-right">
          <text class="menu-value">{{ cacheSize }}</text>
          <text class="menu-arrow">›</text>
        </view>
      </view>
    </view>

    <view class="menu-group">
      <view class="menu-item" @tap="navigateTo('/pages/policy/privacy')">
        <text class="menu-label">隱私政策</text>
        <text class="menu-arrow">›</text>
      </view>
      <view class="menu-item" @tap="navigateTo('/pages/policy/terms')">
        <text class="menu-label">用戶服務協議</text>
        <text class="menu-arrow">›</text>
      </view>
      <view class="menu-item" @tap="navigateTo('/pages/about/index')">
        <text class="menu-label">關於我們</text>
        <text class="menu-arrow">›</text>
      </view>
    </view>

    <!-- Logout -->
    <view class="logout-btn" @tap="handleLogout">
      <text class="logout-text">{{ t('user.logout') }}</text>
    </view>

    <!-- Language picker action sheet -->
    <view v-if="showLangPicker" class="picker-mask" @tap="showLangPicker = false">
      <view class="picker-sheet" @tap.stop>
        <view class="picker-header">
          <text class="picker-title">選擇語言</text>
        </view>
        <view
          v-for="lang in langOptions"
          :key="lang.value"
          class="picker-option"
          :class="{ active: currentLang === lang.value }"
          @tap="handleSetLang(lang.value)"
        >
          <text class="option-text">{{ lang.label }}</text>
          <text v-if="currentLang === lang.value" class="check-mark">✓</text>
        </view>
        <view class="picker-cancel" @tap="showLangPicker = false">
          <text class="cancel-text">{{ t('common.cancel') }}</text>
        </view>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { t, getLang, setLang } from '@/locale'
import type { Lang } from '@/locale'
import { clearUser } from '@/store'

const showLangPicker = ref(false)
const cacheSize = ref('0 KB')
const currentLang = ref<Lang>(getLang())

const langOptions = [
  { label: '繁體中文', value: 'zhHant' as Lang },
  { label: '简体中文', value: 'zhHans' as Lang },
  { label: 'English', value: 'en' as Lang },
]

const currentLangLabel = computed(() => {
  const found = langOptions.find((l) => l.value === currentLang.value)
  return found?.label || '繁體中文'
})

function navigateTo(url: string) {
  uni.navigateTo({ url })
}

function handleSetLang(lang: Lang) {
  currentLang.value = lang
  setLang(lang)
  showLangPicker.value = false
  uni.showToast({ title: '語言已切換', icon: 'success' })
  // Reload pages after a brief delay
  setTimeout(() => {
    uni.reLaunch({ url: '/pages/index/index' })
  }, 800)
}

async function getCacheSize() {
  try {
    const res = await uni.getStorageInfo()
    const sizeKB = res.currentSize || 0
    if (sizeKB > 1024) {
      cacheSize.value = `${(sizeKB / 1024).toFixed(1)} MB`
    } else {
      cacheSize.value = `${sizeKB} KB`
    }
  } catch {
    cacheSize.value = '0 KB'
  }
}

function handleClearCache() {
  uni.showModal({
    title: '提示',
    content: '確定清除所有緩存？',
    success: (res) => {
      if (!res.confirm) return
      try {
        const token = uni.getStorageSync('token')
        const lang = uni.getStorageSync('lang')
        uni.clearStorageSync()
        if (token) uni.setStorageSync('token', token)
        if (lang) uni.setStorageSync('lang', lang)
        cacheSize.value = '0 KB'
        uni.showToast({ title: '清除成功', icon: 'success' })
      } catch {
        uni.showToast({ title: '清除失敗', icon: 'none' })
      }
    },
  })
}

async function handleLogout() {
  uni.showModal({
    title: '提示',
    content: '確定退出登入？',
    success: (res) => {
      if (res.confirm) {
        try {
          uni.removeStorageSync('token')
          uni.removeStorageSync('userInfo')
          uni.removeStorageSync('user')
        } catch (e) {}
        try {
          clearUser()
        } catch (e) {}
        uni.reLaunch({ url: '/pages/login/index' })
      }
    },
  })
}

onMounted(() => {
  getCacheSize()
})
</script>

<style scoped>
.setting-page {
  min-height: 100vh;
  background: #f7f8fa;
  padding: 24rpx;
}

.menu-group {
  background: #fff;
  border-radius: 16rpx;
  margin-bottom: 24rpx;
  overflow: hidden;
}

.menu-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 32rpx 30rpx;
  border-bottom: 1rpx solid #f5f5f5;
}

.menu-item:last-child {
  border-bottom: none;
}

.menu-label {
  font-size: 30rpx;
  color: #333;
}

.menu-right {
  display: flex;
  align-items: center;
}

.menu-value {
  font-size: 26rpx;
  color: #999;
  margin-right: 12rpx;
}

.menu-arrow {
  font-size: 32rpx;
  color: #ccc;
}

.logout-btn {
  margin-top: 60rpx;
  background: #fff;
  border-radius: 16rpx;
  height: 96rpx;
  display: flex;
  align-items: center;
  justify-content: center;
}

.logout-text {
  font-size: 32rpx;
  color: #e74c3c;
  font-weight: 600;
}

/* Language picker sheet */
.picker-mask {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  z-index: 999;
  display: flex;
  align-items: flex-end;
}

.picker-sheet {
  width: 100%;
  background: #fff;
  border-radius: 24rpx 24rpx 0 0;
  padding-bottom: env(safe-area-inset-bottom);
}

.picker-header {
  padding: 30rpx;
  text-align: center;
  border-bottom: 1rpx solid #f0f0f0;
}

.picker-title {
  font-size: 32rpx;
  font-weight: 600;
  color: #333;
}

.picker-option {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 32rpx 40rpx;
  border-bottom: 1rpx solid #f5f5f5;
}

.picker-option.active {
  background: rgba(15, 58, 87, 0.04);
}

.option-text {
  font-size: 30rpx;
  color: #333;
}

.picker-option.active .option-text {
  color: #0f3a57;
  font-weight: 600;
}

.check-mark {
  font-size: 32rpx;
  color: #0f3a57;
  font-weight: 600;
}

.picker-cancel {
  padding: 30rpx;
  text-align: center;
  border-top: 12rpx solid #f7f8fa;
}

.cancel-text {
  font-size: 30rpx;
  color: #999;
}
</style>
