<template>
  <view class="team-page">
    <!-- Stats header -->
    <view class="stats-header">
      <view class="stat-block">
        <text class="stat-num">{{ teamStats.total || 0 }}</text>
        <text class="stat-desc">總成員</text>
      </view>
      <view class="stat-divider" />
      <view class="stat-block">
        <text class="stat-num">{{ teamStats.this_month || 0 }}</text>
        <text class="stat-desc">本月新增</text>
      </view>
    </view>

    <!-- Tabs -->
    <view class="tab-bar">
      <view
        class="tab-item"
        :class="{ active: activeTab === 1 }"
        @click="switchTab(1)"
      >
        <text>一級成員</text>
      </view>
      <view
        class="tab-item"
        :class="{ active: activeTab === 2 }"
        @click="switchTab(2)"
      >
        <text>二級成員</text>
      </view>
    </view>

    <!-- Members list -->
    <view class="members-list">
      <view v-for="member in memberList" :key="member.id" class="member-card">
        <image class="member-avatar" :src="member.avatar || '/static/default-avatar.png'" />
        <view class="member-info">
          <text class="member-name">{{ member.nickname }}</text>
          <text class="member-date">加入時間: {{ member.joined_at }}</text>
        </view>
        <view class="member-stats">
          <text class="member-orders">{{ member.total_orders || 0 }} 單</text>
        </view>
      </view>
    </view>

    <!-- Empty / Loading -->
    <view v-if="!loading && memberList.length === 0" class="empty">
      <text>{{ t('common.noData') }}</text>
    </view>
    <view class="load-more">
      <text v-if="loading">{{ t('common.loading') }}</text>
      <text v-else-if="finished && memberList.length > 0">— 已經到底了 —</text>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { t } from '@/locale'
import { commonApi } from '@/api/common'

const activeTab = ref(1)
const teamStats = ref<any>({})
const memberList = ref<any[]>([])
const loading = ref(false)
const finished = ref(false)
const page = ref(1)
const limit = 15

onMounted(() => {
  loadStats()
  loadMembers()
})

async function loadStats() {
  try {
    const res = await commonApi.getDealerTeam({ stats_only: true })
    teamStats.value = res?.data || {}
  } catch (e) {
    console.error(e)
  }
}

function switchTab(tab: number) {
  activeTab.value = tab
  resetList()
}

function resetList() {
  page.value = 1
  finished.value = false
  memberList.value = []
  loadMembers()
}

async function loadMembers() {
  if (loading.value || finished.value) return
  loading.value = true
  try {
    const res = await commonApi.getDealerTeam({
      level: activeTab.value,
      page: page.value,
      limit,
    })
    const list = res?.data?.list || res?.data || []
    if (list.length < limit) {
      finished.value = true
    }
    memberList.value.push(...list)
    page.value++
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

onReachBottom(() => {
  loadMembers()
})
</script>

<style scoped>
.team-page {
  min-height: 100vh;
  background: #f5f5f5;
}

.stats-header {
  display: flex;
  background: #0f3a57;
  padding: 40rpx 30rpx;
  align-items: center;
  justify-content: center;
}

.stat-block {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8rpx;
}

.stat-num {
  font-size: 44rpx;
  font-weight: 700;
  color: #fff;
}

.stat-desc {
  font-size: 24rpx;
  color: rgba(255, 255, 255, 0.7);
}

.stat-divider {
  width: 1rpx;
  height: 60rpx;
  background: rgba(255, 255, 255, 0.3);
}

.tab-bar {
  display: flex;
  background: #fff;
  margin-bottom: 16rpx;
}

.tab-item {
  flex: 1;
  text-align: center;
  padding: 24rpx 0;
  font-size: 28rpx;
  color: #666;
  position: relative;
}

.tab-item.active {
  color: #0f3a57;
  font-weight: 600;
}

.tab-item.active::after {
  content: '';
  position: absolute;
  bottom: 0;
  left: 50%;
  transform: translateX(-50%);
  width: 60rpx;
  height: 4rpx;
  background: #0f3a57;
  border-radius: 2rpx;
}

.members-list {
  padding: 0 24rpx;
}

.member-card {
  display: flex;
  align-items: center;
  background: #fff;
  border-radius: 16rpx;
  padding: 24rpx;
  margin-bottom: 12rpx;
  gap: 20rpx;
}

.member-avatar {
  width: 80rpx;
  height: 80rpx;
  border-radius: 50%;
  flex-shrink: 0;
}

.member-info {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 8rpx;
}

.member-name {
  font-size: 28rpx;
  color: #333;
  font-weight: 500;
}

.member-date {
  font-size: 22rpx;
  color: #999;
}

.member-stats {
  flex-shrink: 0;
}

.member-orders {
  font-size: 26rpx;
  color: #0f3a57;
  font-weight: 500;
}

.empty {
  text-align: center;
  padding: 100rpx 0;
  color: #999;
  font-size: 28rpx;
}

.load-more {
  text-align: center;
  padding: 30rpx 0;
  font-size: 24rpx;
  color: #999;
}
</style>
