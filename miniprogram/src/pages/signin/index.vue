<template>
  <view class="signin-page">
    <!-- Loading -->
    <view v-if="loading" class="loading-wrap">
      <text class="loading-text">{{ t('common.loading') }}</text>
    </view>

    <view v-else class="signin-content">
      <!-- Consecutive days & points -->
      <view class="stats-section">
        <view class="stat-item">
          <text class="stat-value">{{ signInData.consecutive_days }}</text>
          <text class="stat-label">Consecutive Days</text>
        </view>
        <view class="stat-divider" />
        <view class="stat-item">
          <text class="stat-value">+{{ signInData.points_reward }}</text>
          <text class="stat-label">Points Today</text>
        </view>
      </view>

      <!-- Calendar -->
      <view class="calendar-section">
        <view class="calendar-header">
          <text class="calendar-month">{{ currentMonth }}</text>
        </view>
        <view class="calendar-weekdays">
          <text v-for="day in weekdays" :key="day" class="weekday">{{ day }}</text>
        </view>
        <view class="calendar-grid">
          <view
            v-for="(cell, idx) in calendarCells"
            :key="idx"
            class="calendar-cell"
            :class="{
              empty: !cell.day,
              checked: cell.checked,
              today: cell.isToday,
            }"
          >
            <text v-if="cell.day" class="cell-day">{{ cell.day }}</text>
            <view v-if="cell.checked" class="check-mark" />
          </view>
        </view>
      </view>

      <!-- Check-in button -->
      <view class="btn-section">
        <view
          class="signin-btn"
          :class="{ disabled: signInData.checked_today }"
          @click="doSignIn"
        >
          <text class="signin-btn-text">
            {{ signInData.checked_today ? 'Already Checked In' : t('signin.title') }}
          </text>
        </view>
      </view>

      <!-- Reward rules -->
      <view class="rules-section">
        <text class="rules-title">Reward Rules</text>
        <view class="rules-list">
          <text class="rule-item">- Daily check-in: +{{ signInData.daily_points || 1 }} points</text>
          <text class="rule-item">- 7 consecutive days: bonus points</text>
          <text class="rule-item">- 30 consecutive days: special reward</text>
        </view>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { t } from '@/locale'
import { commonApi } from '@/api/common'
import { userApi } from '@/api/user'

interface SignInData {
  consecutive_days: number
  points_reward: number
  daily_points: number
  checked_today: boolean
  checked_dates: string[] // 'YYYY-MM-DD' format
}

const loading = ref(true)
const signInData = ref<SignInData>({
  consecutive_days: 0,
  points_reward: 0,
  daily_points: 1,
  checked_today: false,
  checked_dates: [],
})

const weekdays = ['Sun', 'Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat']

const now = new Date()
const currentMonth = computed(() => {
  const months = [
    'January', 'February', 'March', 'April', 'May', 'June',
    'July', 'August', 'September', 'October', 'November', 'December',
  ]
  return `${months[now.getMonth()]} ${now.getFullYear()}`
})

interface CalendarCell {
  day: number | null
  checked: boolean
  isToday: boolean
}

const calendarCells = computed<CalendarCell[]>(() => {
  const year = now.getFullYear()
  const month = now.getMonth()
  const firstDay = new Date(year, month, 1).getDay()
  const daysInMonth = new Date(year, month + 1, 0).getDate()
  const today = now.getDate()

  const cells: CalendarCell[] = []

  // Empty cells before first day
  for (let i = 0; i < firstDay; i++) {
    cells.push({ day: null, checked: false, isToday: false })
  }

  // Day cells
  for (let d = 1; d <= daysInMonth; d++) {
    const dateStr = `${year}-${String(month + 1).padStart(2, '0')}-${String(d).padStart(2, '0')}`
    cells.push({
      day: d,
      checked: signInData.value.checked_dates.includes(dateStr),
      isToday: d === today,
    })
  }

  return cells
})

async function loadSignInStatus() {
  loading.value = true
  try {
    const res: any = await userApi.getSignInStatus()
    const data = res?.data || res
    signInData.value = {
      consecutive_days: data.consecutive_days || 0,
      points_reward: data.points_reward || 0,
      daily_points: data.daily_points || 1,
      checked_today: data.checked_today || false,
      checked_dates: data.checked_dates || [],
    }
  } catch (e: any) {
    uni.showToast({ title: e?.message || t('common.retry'), icon: 'none' })
  } finally {
    loading.value = false
  }
}

async function doSignIn() {
  if (signInData.value.checked_today) return
  try {
    const res: any = await userApi.signIn()
    const data = res?.data || res
    signInData.value.checked_today = true
    signInData.value.consecutive_days = data.consecutive_days || signInData.value.consecutive_days + 1
    signInData.value.points_reward = data.points_reward || signInData.value.daily_points

    // Add today's date to checked dates
    const year = now.getFullYear()
    const month = String(now.getMonth() + 1).padStart(2, '0')
    const day = String(now.getDate()).padStart(2, '0')
    const todayStr = `${year}-${month}-${day}`
    if (!signInData.value.checked_dates.includes(todayStr)) {
      signInData.value.checked_dates.push(todayStr)
    }

    uni.showToast({ title: `+${signInData.value.points_reward} points!`, icon: 'success' })
  } catch (e: any) {
    uni.showToast({ title: e?.message || t('common.retry'), icon: 'none' })
  }
}

onMounted(() => {
  uni.setNavigationBarTitle({ title: t('signin.title') })
  loadSignInStatus()
})
</script>

<style scoped>
.signin-page {
  min-height: 100vh;
  background-color: #f5f5f5;
}

.loading-wrap {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 120rpx 0;
}

.loading-text {
  font-size: 28rpx;
  color: #999;
}

.stats-section {
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: #0f3a57;
  padding: 48rpx 32rpx;
}

.stat-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  flex: 1;
}

.stat-value {
  font-size: 48rpx;
  font-weight: 700;
  color: #fff;
}

.stat-label {
  font-size: 24rpx;
  color: rgba(255, 255, 255, 0.7);
  margin-top: 8rpx;
}

.stat-divider {
  width: 1rpx;
  height: 60rpx;
  background-color: rgba(255, 255, 255, 0.3);
}

.calendar-section {
  background-color: #fff;
  margin: 20rpx 24rpx;
  border-radius: 16rpx;
  padding: 24rpx;
}

.calendar-header {
  text-align: center;
  margin-bottom: 20rpx;
}

.calendar-month {
  font-size: 30rpx;
  font-weight: 600;
  color: #333;
}

.calendar-weekdays {
  display: flex;
}

.weekday {
  flex: 1;
  text-align: center;
  font-size: 24rpx;
  color: #999;
  padding: 12rpx 0;
}

.calendar-grid {
  display: flex;
  flex-wrap: wrap;
}

.calendar-cell {
  width: calc(100% / 7);
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 16rpx 0;
  position: relative;
}

.cell-day {
  font-size: 28rpx;
  color: #333;
}

.calendar-cell.empty .cell-day {
  color: transparent;
}

.calendar-cell.today {
  background-color: rgba(15, 58, 87, 0.1);
  border-radius: 8rpx;
}

.calendar-cell.checked {
  background-color: #0f3a57;
  border-radius: 8rpx;
}

.calendar-cell.checked .cell-day {
  color: #fff;
}

.check-mark {
  width: 8rpx;
  height: 8rpx;
  border-radius: 50%;
  background-color: #fff;
  margin-top: 4rpx;
}

.btn-section {
  padding: 40rpx 48rpx;
}

.signin-btn {
  background-color: #0f3a57;
  border-radius: 50rpx;
  height: 100rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 8rpx 24rpx rgba(15, 58, 87, 0.3);
}

.signin-btn.disabled {
  background-color: #ccc;
  box-shadow: none;
}

.signin-btn-text {
  font-size: 34rpx;
  color: #fff;
  font-weight: 600;
}

.rules-section {
  background-color: #fff;
  margin: 0 24rpx 40rpx;
  border-radius: 16rpx;
  padding: 24rpx;
}

.rules-title {
  font-size: 28rpx;
  font-weight: 600;
  color: #333;
  margin-bottom: 16rpx;
  display: block;
}

.rules-list {
  display: flex;
  flex-direction: column;
  gap: 12rpx;
}

.rule-item {
  font-size: 26rpx;
  color: #666;
  line-height: 1.6;
}
</style>
