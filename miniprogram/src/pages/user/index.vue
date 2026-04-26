<template>
  <view class="page">
    <!-- Custom Navbar -->
    <view class="navbar" :style="{ paddingTop: statusBarHeight + 'px' }">
      <view class="navbar-content">
        <text class="navbar-title">{{ t('user.center') }}</text>
        <view class="navbar-right">
          <view class="navbar-icon" @tap="goToMessage">
            <text class="icon-text">&#x1F514;</text>
            <view v-if="messageCount > 0" class="badge">{{ messageCount > 99 ? '99+' : messageCount }}</view>
          </view>
          <view class="navbar-icon" @tap="goToSettings">
            <text class="icon-text">&#x2699;</text>
          </view>
        </view>
      </view>
    </view>

    <scroll-view
      class="scroll-content"
      scroll-y
      :style="{ paddingTop: navbarHeight + 'px' }"
    >
      <!-- User Header -->
      <view class="user-header" @tap="onHeaderTap">
        <view class="avatar-wrapper">
          <image
            v-if="userInfo?.avatar"
            class="avatar"
            :src="userInfo.avatar"
            mode="aspectFill"
          />
          <view v-else class="avatar avatar-placeholder">
            <text class="avatar-text">&#x1F464;</text>
          </view>
          <!-- VIP badge -->
          <view v-if="userInfo?.is_vip" class="vip-badge">
            <text class="vip-text">VIP{{ userInfo.vip_level || '' }}</text>
          </view>
        </view>
        <view class="user-info">
          <text class="nickname">{{ userInfo?.nickname || t('user.login') }}</text>
          <text v-if="userInfo?.phone" class="phone">{{ maskPhone(userInfo.phone) }}</text>
          <text v-if="!isLoggedIn" class="login-hint">點擊登入，享受更多服務</text>
        </view>
        <text class="header-arrow">&#x276F;</text>
      </view>

      <!-- Stats Row -->
      <view class="stats-row">
        <view class="stats-item" @tap="goToBalance">
          <text class="stats-value">{{ userInfo?.balance?.toFixed(2) || '0.00' }}</text>
          <text class="stats-label">{{ t('user.balance') }}</text>
        </view>
        <view class="stats-divider" />
        <view class="stats-item" @tap="goToPoints">
          <text class="stats-value">{{ userInfo?.points || 0 }}</text>
          <text class="stats-label">{{ t('user.points') }}</text>
        </view>
        <view class="stats-divider" />
        <view class="stats-item" @tap="goToCoupons">
          <text class="stats-value">{{ couponCount }}</text>
          <text class="stats-label">{{ t('user.coupons') }}</text>
        </view>
      </view>

      <!-- Order Section -->
      <view class="section-card">
        <view class="section-header" @tap="goToOrderList">
          <text class="section-title">{{ t('user.myOrders') }}</text>
          <view class="section-more">
            <text class="section-more-text">全部訂單</text>
            <text class="section-arrow">&#x276F;</text>
          </view>
        </view>
        <view class="order-grid">
          <view
            v-for="status in orderStatuses"
            :key="status.key"
            class="order-item"
            @tap="goToOrderListByStatus(status.key)"
          >
            <view class="order-icon-wrapper">
              <text class="order-icon">{{ status.icon }}</text>
              <view v-if="status.count > 0" class="order-badge">{{ status.count }}</view>
            </view>
            <text class="order-label">{{ status.label }}</text>
          </view>
        </view>
      </view>

      <!-- Package Section -->
      <view class="section-card">
        <view class="section-header" @tap="goToPackageList">
          <text class="section-title">{{ t('user.myPackages') }}</text>
          <view class="section-more">
            <text class="section-more-text">{{ t('common.all') }}</text>
            <text class="section-arrow">&#x276F;</text>
          </view>
        </view>
        <view class="order-grid">
          <view
            v-for="status in packageStatuses"
            :key="status.key"
            class="order-item"
            @tap="goToPackageListByStatus(status.key)"
          >
            <view class="order-icon-wrapper">
              <text class="order-icon">{{ status.icon }}</text>
            </view>
            <text class="order-label">{{ status.label }}</text>
          </view>
        </view>
      </view>

      <!-- Menu List -->
      <view class="menu-card">
        <view
          v-for="item in menuItems"
          :key="item.key"
          class="menu-item"
          @tap="onMenuTap(item)"
        >
          <view class="menu-left">
            <text class="menu-icon">{{ item.icon }}</text>
            <text class="menu-label">{{ item.label }}</text>
          </view>
          <text class="menu-arrow">&#x276F;</text>
        </view>
      </view>

      <!-- Bottom padding -->
      <view style="height: 40rpx;" />
    </scroll-view>
  </view>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { onShow } from '@dcloudio/uni-app'
import { t } from '@/locale'
import { userApi } from '@/api/user'
import { commonApi } from '@/api/common'
import store, { setUser } from '@/store'
import type { UserInfo } from '@/store'

// --- System info ---
const statusBarHeight = ref(20)
const navbarHeight = computed(() => statusBarHeight.value + 44)

onMounted(() => {
  const sysInfo = uni.getSystemInfoSync()
  statusBarHeight.value = sysInfo.statusBarHeight || 20
})

// --- User data ---
const isLoggedIn = computed(() => store.isLoggedIn)
const userInfo = computed(() => store.user)
const messageCount = computed(() => store.messageCount)
const couponCount = ref(0)

// --- Order statuses ---
const orderStatuses = ref([
  { key: 'pending', label: '待付款', icon: '&#x1F4B3;', count: 0 },
  { key: 'paid', label: '待發貨', icon: '&#x1F4E6;', count: 0 },
  { key: 'shipped', label: '待收貨', icon: '&#x1F69A;', count: 0 },
  { key: 'completed', label: '已完成', icon: '&#x2705;', count: 0 },
  { key: 'refund', label: '售後', icon: '&#x1F527;', count: 0 },
])

// Re-define with actual emoji rendering (the HTML entities won't render in <text>)
const orderStatusesData = [
  { key: 'pending', label: '待付款', icon: '\uD83D\uDCB3', count: 0 },
  { key: 'paid', label: '待發貨', icon: '\uD83D\uDCE6', count: 0 },
  { key: 'shipped', label: '待收貨', icon: '\uD83D\uDE9A', count: 0 },
  { key: 'completed', label: '已完成', icon: '\u2705', count: 0 },
  { key: 'refund', label: '售後', icon: '\uD83D\uDD27', count: 0 },
]

// Assign proper emoji data
orderStatuses.value = orderStatusesData

const packageStatuses = [
  { key: 'pending', label: '待入庫', icon: '📥' },
  { key: 'stored', label: '已入庫', icon: '📦' },
  { key: 'packed', label: '已打包', icon: '📮' },
  { key: 'shipped', label: '已發貨', icon: '🚚' },
  { key: 'delivered', label: '已送達', icon: '✅' },
]

// --- Menu items ---
const menuItems = [
  { key: 'packages', label: '我的包裹', icon: '📦', path: '/pages/package/index' },
  { key: 'address', label: '地址管理', icon: '📍', path: '/pages/address/list' },
  { key: 'warehouse', label: '倉庫地址', icon: '🏠', path: '/pages/warehouse/address' },
  { key: 'favorite', label: '收藏列表', icon: '❤️', path: '/pages/favorite/index' },
  { key: 'history', label: '瀏覽記錄', icon: '🕐', path: '/pages/history/index' },
  { key: 'dealer', label: '分銷中心', icon: '💰', path: '/pages/dealer/index' },
  { key: 'help', label: '幫助中心', icon: '❓', path: '/pages/help/index' },
  { key: 'feedback', label: '意見反饋', icon: '💬', path: '/pages/feedback/index' },
  { key: 'about', label: '關於我們', icon: 'ℹ️', path: '/pages/about/index' },
  { key: 'settings', label: '設置', icon: '⚙️', path: '/pages/setting/index' },
]

// --- Lifecycle ---
onShow(() => {
  if (isLoggedIn.value) {
    loadUserData()
  }
})

async function loadUserData() {
  await Promise.all([
    fetchUserInfo(),
    fetchCouponCount(),
    fetchOrderCounts(),
  ])
}

async function fetchUserInfo() {
  try {
    const res: any = await userApi.getUserInfo()
    const data = res?.data || res
    if (data?.id) {
      setUser(data as UserInfo)
    }
  } catch (e) {
    // Silent
  }
}

async function fetchCouponCount() {
  try {
    const res: any = await commonApi.getMyCoupons({ page: 1, limit: 1, status: 'unused' })
    const data = res?.data || res
    couponCount.value = data?.total || data?.count || 0
  } catch (e) {
    couponCount.value = 0
  }
}

async function fetchOrderCounts() {
  try {
    const res: any = await orderApi.getList({ count_only: true })
    const data = res?.data || res
    if (data?.counts) {
      orderStatuses.value.forEach(s => {
        s.count = data.counts[s.key] || 0
      })
    }
  } catch (e) {
    // Silent
  }
}

// Import order API (lazy, to avoid circular)
import { orderApi } from '@/api/order'

// --- Helpers ---
function maskPhone(phone: string): string {
  if (!phone || phone.length < 7) return phone
  return phone.substring(0, 3) + '****' + phone.substring(phone.length - 4)
}

// --- Navigation ---
function onHeaderTap() {
  if (isLoggedIn.value) {
    // Could navigate to profile edit page
  } else {
    uni.navigateTo({ url: '/pages/login/index' })
  }
}

function onMenuTap(item: typeof menuItems[0]) {
  if (!isLoggedIn.value && ['packages', 'address', 'favorite', 'history', 'dealer', 'feedback'].includes(item.key)) {
    uni.navigateTo({ url: '/pages/login/index' })
    return
  }
  uni.navigateTo({ url: item.path })
}

function goToMessage() {
  uni.navigateTo({ url: '/pages/message/index' })
}

function goToSettings() {
  uni.navigateTo({ url: '/pages/setting/index' })
}

function goToBalance() {
  if (!isLoggedIn.value) { uni.navigateTo({ url: '/pages/login/index' }); return }
  uni.navigateTo({ url: '/pages/balance/index' })
}

function goToPoints() {
  if (!isLoggedIn.value) { uni.navigateTo({ url: '/pages/login/index' }); return }
  uni.navigateTo({ url: '/pages/points/index' })
}

function goToCoupons() {
  if (!isLoggedIn.value) { uni.navigateTo({ url: '/pages/login/index' }); return }
  uni.navigateTo({ url: '/pages/coupon/my' })
}

function goToOrderList() {
  if (!isLoggedIn.value) { uni.navigateTo({ url: '/pages/login/index' }); return }
  uni.navigateTo({ url: '/pages/order/list' })
}

function goToOrderListByStatus(status: string) {
  if (!isLoggedIn.value) { uni.navigateTo({ url: '/pages/login/index' }); return }
  if (status === 'refund') {
    uni.navigateTo({ url: '/pages/refund/index' })
  } else {
    uni.navigateTo({ url: `/pages/order/list?status=${status}` })
  }
}

function goToPackageList() {
  if (!isLoggedIn.value) { uni.navigateTo({ url: '/pages/login/index' }); return }
  uni.navigateTo({ url: '/pages/package/index' })
}

function goToPackageListByStatus(status: string) {
  if (!isLoggedIn.value) { uni.navigateTo({ url: '/pages/login/index' }); return }
  uni.navigateTo({ url: `/pages/package/index?status=${status}` })
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
  background-color: #0f3a57;
}

.navbar-content {
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 44px;
  padding: 0 30rpx;
}

.navbar-title {
  font-size: 34rpx;
  font-weight: bold;
  color: #ffffff;
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

.icon-text {
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

/* --- Scroll --- */
.scroll-content {
  height: 100vh;
  box-sizing: border-box;
}

/* --- User Header --- */
.user-header {
  display: flex;
  align-items: center;
  padding: 40rpx 30rpx 30rpx;
  background: linear-gradient(180deg, #0f3a57 0%, #164d72 60%, #f5f5f5 100%);
}

.avatar-wrapper {
  position: relative;
  flex-shrink: 0;
  margin-right: 24rpx;
}

.avatar {
  width: 120rpx;
  height: 120rpx;
  border-radius: 50%;
  border: 4rpx solid rgba(255, 255, 255, 0.4);
}

.avatar-placeholder {
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: rgba(255, 255, 255, 0.2);
}

.avatar-text {
  font-size: 60rpx;
}

.vip-badge {
  position: absolute;
  bottom: -4rpx;
  left: 50%;
  transform: translateX(-50%);
  padding: 2rpx 16rpx;
  background: linear-gradient(135deg, #f6d365, #fda085);
  border-radius: 20rpx;
  white-space: nowrap;
}

.vip-text {
  font-size: 18rpx;
  color: #7a4100;
  font-weight: bold;
}

.user-info {
  flex: 1;
  min-width: 0;
}

.nickname {
  font-size: 34rpx;
  font-weight: bold;
  color: #ffffff;
  display: block;
}

.phone {
  font-size: 24rpx;
  color: rgba(255, 255, 255, 0.7);
  margin-top: 6rpx;
  display: block;
}

.login-hint {
  font-size: 24rpx;
  color: rgba(255, 255, 255, 0.6);
  margin-top: 6rpx;
  display: block;
}

.header-arrow {
  font-size: 28rpx;
  color: rgba(255, 255, 255, 0.6);
  flex-shrink: 0;
  margin-left: 12rpx;
}

/* --- Stats Row --- */
.stats-row {
  display: flex;
  align-items: center;
  background-color: #ffffff;
  margin: 0 24rpx;
  border-radius: 16rpx;
  padding: 30rpx 0;
  margin-top: -20rpx;
  position: relative;
  z-index: 2;
  box-shadow: 0 4rpx 16rpx rgba(0, 0, 0, 0.06);
}

.stats-item {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
}

.stats-value {
  font-size: 36rpx;
  font-weight: bold;
  color: #0f3a57;
}

.stats-label {
  font-size: 22rpx;
  color: #999999;
  margin-top: 6rpx;
}

.stats-divider {
  width: 1rpx;
  height: 50rpx;
  background-color: #eeeeee;
}

/* --- Section Card --- */
.section-card {
  background-color: #ffffff;
  margin: 20rpx 24rpx 0;
  border-radius: 16rpx;
  padding: 24rpx;
}

.section-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 20rpx;
}

.section-title {
  font-size: 30rpx;
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

.section-arrow {
  font-size: 22rpx;
  color: #cccccc;
  margin-left: 4rpx;
}

/* --- Order Grid --- */
.order-grid {
  display: flex;
  justify-content: space-around;
}

.order-item {
  display: flex;
  flex-direction: column;
  align-items: center;
}

.order-icon-wrapper {
  position: relative;
  width: 64rpx;
  height: 64rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 10rpx;
}

.order-icon {
  font-size: 44rpx;
}

.order-badge {
  position: absolute;
  top: -6rpx;
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

.order-label {
  font-size: 22rpx;
  color: #666666;
}

/* --- Menu List --- */
.menu-card {
  background-color: #ffffff;
  margin: 20rpx 24rpx 0;
  border-radius: 16rpx;
  overflow: hidden;
}

.menu-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 28rpx 24rpx;
  border-bottom: 1rpx solid #f5f5f5;
}

.menu-item:last-child {
  border-bottom: none;
}

.menu-left {
  display: flex;
  align-items: center;
  gap: 16rpx;
}

.menu-icon {
  font-size: 36rpx;
  width: 44rpx;
  text-align: center;
}

.menu-label {
  font-size: 28rpx;
  color: #333333;
}

.menu-arrow {
  font-size: 24rpx;
  color: #cccccc;
}
</style>
