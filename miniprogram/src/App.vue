<script setup lang="ts">
import { onLaunch, onShow, onHide } from '@dcloudio/uni-app'
import { initLang } from '@/locale'
import { getToken } from '@/utils/request'
import { store, setUser } from '@/store'

onLaunch(async () => {
  console.log('App Launch')

  // Initialize language from backend
  await initLang()

  // Restore login state from storage
  const token = getToken()
  if (token) {
    store.isLoggedIn = true
    const cached = uni.getStorageSync('userInfo')
    if (cached) {
      try { setUser(typeof cached === 'string' ? JSON.parse(cached) : cached) } catch {}
    }
  }
})

onShow(() => {
  console.log('App Show')
})

onHide(() => {
  console.log('App Hide')
})
</script>

<style>
/* Global styles */
page {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', 'PingFang SC', 'Hiragino Sans GB',
    'Microsoft YaHei', 'Helvetica Neue', Helvetica, Arial, sans-serif;
  font-size: 28rpx;
  color: #333333;
  background-color: #f5f5f5;
  box-sizing: border-box;
}

/* Reset */
view, text, image, navigator, scroll-view, swiper, swiper-item, input, textarea, button {
  box-sizing: border-box;
}

/* Common utility classes */
.flex { display: flex; }
.flex-col { display: flex; flex-direction: column; }
.flex-row { display: flex; flex-direction: row; }
.flex-wrap { flex-wrap: wrap; }
.flex-1 { flex: 1; }
.items-center { align-items: center; }
.justify-center { justify-content: center; }
.justify-between { justify-content: space-between; }
.justify-around { justify-content: space-around; }
.text-center { text-align: center; }
.text-right { text-align: right; }
.text-left { text-align: left; }

/* Font sizes */
.text-xs { font-size: 20rpx; }
.text-sm { font-size: 24rpx; }
.text-base { font-size: 28rpx; }
.text-lg { font-size: 32rpx; }
.text-xl { font-size: 36rpx; }
.text-2xl { font-size: 40rpx; }

/* Colors */
.text-primary { color: #0f3a57; }
.text-gray { color: #999999; }
.text-light { color: #cccccc; }
.text-danger { color: #e74c3c; }
.text-success { color: #27ae60; }
.text-warning { color: #f39c12; }
.text-white { color: #ffffff; }
.bg-primary { background-color: #0f3a57; }
.bg-white { background-color: #ffffff; }

/* Spacing */
.p-0 { padding: 0; }
.p-1 { padding: 10rpx; }
.p-2 { padding: 20rpx; }
.p-3 { padding: 30rpx; }
.p-4 { padding: 40rpx; }
.m-0 { margin: 0; }
.m-1 { margin: 10rpx; }
.m-2 { margin: 20rpx; }
.m-3 { margin: 30rpx; }
.mt-1 { margin-top: 10rpx; }
.mt-2 { margin-top: 20rpx; }
.mt-3 { margin-top: 30rpx; }
.mb-1 { margin-bottom: 10rpx; }
.mb-2 { margin-bottom: 20rpx; }
.mb-3 { margin-bottom: 30rpx; }
.mx-auto { margin-left: auto; margin-right: auto; }

/* Border radius */
.rounded { border-radius: 8rpx; }
.rounded-lg { border-radius: 16rpx; }
.rounded-full { border-radius: 50%; }

/* Shadow */
.shadow { box-shadow: 0 2rpx 12rpx rgba(0, 0, 0, 0.1); }
.shadow-lg { box-shadow: 0 4rpx 24rpx rgba(0, 0, 0, 0.15); }

/* Card */
.card {
  background: #ffffff;
  border-radius: 16rpx;
  padding: 24rpx;
  margin: 20rpx;
  box-shadow: 0 2rpx 12rpx rgba(0, 0, 0, 0.06);
}

/* Button styles */
.btn {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 80rpx;
  border-radius: 40rpx;
  font-size: 30rpx;
  font-weight: 500;
}

.btn-primary {
  background: #0f3a57;
  color: #ffffff;
}

.btn-outline {
  background: transparent;
  border: 2rpx solid #0f3a57;
  color: #0f3a57;
}

.btn-danger {
  background: #e74c3c;
  color: #ffffff;
}

.btn-sm {
  height: 56rpx;
  border-radius: 28rpx;
  font-size: 24rpx;
  padding: 0 24rpx;
}

/* Safe area bottom padding */
.safe-bottom {
  padding-bottom: env(safe-area-inset-bottom);
}

/* Divider */
.divider {
  height: 1rpx;
  background: #eeeeee;
  margin: 20rpx 0;
}

/* Badge */
.badge {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-width: 32rpx;
  height: 32rpx;
  padding: 0 8rpx;
  border-radius: 16rpx;
  background: #e74c3c;
  color: #ffffff;
  font-size: 20rpx;
}

/* Empty state */
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 100rpx 40rpx;
  color: #999999;
}

.empty-state .empty-icon {
  font-size: 120rpx;
  margin-bottom: 20rpx;
}

.empty-state .empty-text {
  font-size: 28rpx;
  color: #999999;
}

/* Loading */
.loading-more {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 30rpx;
  color: #999999;
  font-size: 24rpx;
}

/* Status badge */
.status-badge {
  display: inline-block;
  padding: 4rpx 16rpx;
  border-radius: 20rpx;
  font-size: 22rpx;
}

.status-badge.pending { background: #fff3e0; color: #f57c00; }
.status-badge.success { background: #e8f5e9; color: #388e3c; }
.status-badge.danger { background: #ffebee; color: #d32f2f; }
.status-badge.info { background: #e3f2fd; color: #1976d2; }
.status-badge.default { background: #f5f5f5; color: #666666; }
</style>
