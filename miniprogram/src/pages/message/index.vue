<template>
  <view class="message-page">
    <!-- Tabs -->
    <view class="msg-tabs">
      <view
        v-for="tab in tabs"
        :key="tab.value"
        class="tab-item"
        :class="{ active: activeTab === tab.value }"
        @click="switchTab(tab.value)"
      >
        <text class="tab-text">{{ tab.label }}</text>
        <view v-if="activeTab === tab.value" class="tab-line" />
      </view>
    </view>

    <!-- Mark all read -->
    <view class="toolbar" v-if="list.length > 0">
      <view class="mark-all-btn" @click="markAllRead">
        <text class="mark-all-text">Mark all read</text>
      </view>
    </view>

    <!-- Loading -->
    <view v-if="loading && list.length === 0" class="loading-wrap">
      <text class="loading-text">{{ t('common.loading') }}</text>
    </view>

    <!-- Error -->
    <view v-else-if="error && list.length === 0" class="error-wrap">
      <text class="error-text">{{ error }}</text>
      <view class="retry-btn" @click="loadMessages(true)">
        <text class="retry-text">{{ t('common.retry') }}</text>
      </view>
    </view>

    <!-- Messages -->
    <view v-else class="msg-items">
      <view
        v-for="item in list"
        :key="item.id"
        class="msg-item"
        @click="onItemTap(item)"
      >
        <view class="msg-icon-wrap">
          <view class="msg-icon" :class="'icon-' + item.type">
            <text class="icon-text">{{ getIcon(item.type) }}</text>
          </view>
          <view v-if="!item.is_read" class="msg-dot" />
        </view>
        <view class="msg-content">
          <text class="msg-title" :class="{ unread: !item.is_read }">
            {{ item.title }}
          </text>
          <text class="msg-brief">{{ item.brief }}</text>
        </view>
        <text class="msg-time">{{ item.time }}</text>
      </view>

      <!-- Empty -->
      <view v-if="!loading && list.length === 0" class="empty-wrap">
        <text class="empty-text">{{ t('common.noData') }}</text>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { t } from '@/locale'
import { commonApi } from '@/api/common'

interface MessageItem {
  id: string
  type: string
  title: string
  brief: string
  time: string
  is_read: boolean
  link?: string
}

const tabs = [
  { label: '系統通知', value: 'system' },
  { label: '物流消息', value: 'logistics' },
  { label: '訂單消息', value: 'order' },
]

const activeTab = ref('system')
const loading = ref(true)
const error = ref('')
const list = ref<MessageItem[]>([])

function getIcon(type: string): string {
  const map: Record<string, string> = {
    system: '\u{1F514}',
    logistics: '\u{1F4E6}',
    order: '\u{1F4CB}',
  }
  return map[type] || '\u{1F4E9}'
}

async function loadMessages(reset = false) {
  if (reset) {
    list.value = []
    loading.value = true
  }
  error.value = ''
  try {
    const res: any = await commonApi.getMessages({ type: activeTab.value })
    const data = res?.data || res
    list.value = Array.isArray(data) ? data : data?.list || []
  } catch (e: any) {
    error.value = e?.message || t('common.retry')
  } finally {
    loading.value = false
  }
}

function switchTab(value: string) {
  activeTab.value = value
  loadMessages(true)
}

async function markAllRead() {
  try {
    for (const item of list.value) {
      if (!item.is_read) {
        await commonApi.markRead(item.id)
        item.is_read = true
      }
    }
    uni.showToast({ title: t('common.done'), icon: 'success' })
  } catch (e: any) {
    uni.showToast({ title: e?.message || t('common.retry'), icon: 'none' })
  }
}

async function onItemTap(item: MessageItem) {
  if (!item.is_read) {
    try {
      await commonApi.markRead(item.id)
      item.is_read = true
    } catch (_) {
      // silent
    }
  }
  if (item.link) {
    uni.navigateTo({ url: item.link })
  }
}

onMounted(() => {
  uni.setNavigationBarTitle({ title: t('message.title') })
  loadMessages(true)
})
</script>

<style scoped>
.message-page {
  min-height: 100vh;
  background-color: #f5f5f5;
}

.msg-tabs {
  display: flex;
  background-color: #fff;
  border-bottom: 1rpx solid #eee;
}

.tab-item {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 24rpx 0;
  position: relative;
}

.tab-text {
  font-size: 28rpx;
  color: #666;
}

.tab-item.active .tab-text {
  color: #0f3a57;
  font-weight: 600;
}

.tab-line {
  position: absolute;
  bottom: 0;
  width: 48rpx;
  height: 4rpx;
  background-color: #0f3a57;
  border-radius: 2rpx;
}

.toolbar {
  display: flex;
  justify-content: flex-end;
  padding: 16rpx 24rpx;
}

.mark-all-btn {
  padding: 8rpx 16rpx;
}

.mark-all-text {
  font-size: 24rpx;
  color: #0f3a57;
}

.loading-wrap,
.error-wrap,
.empty-wrap {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 120rpx 0;
}

.loading-text,
.empty-text {
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

.msg-items {
  padding: 0 24rpx;
}

.msg-item {
  display: flex;
  align-items: flex-start;
  background-color: #fff;
  border-radius: 12rpx;
  padding: 24rpx;
  margin-bottom: 16rpx;
}

.msg-icon-wrap {
  position: relative;
  margin-right: 20rpx;
  flex-shrink: 0;
}

.msg-icon {
  width: 72rpx;
  height: 72rpx;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: #f0f0f0;
}

.icon-system {
  background-color: #e8f0fe;
}

.icon-logistics {
  background-color: #fef3e0;
}

.icon-order {
  background-color: #e8f5e9;
}

.icon-text {
  font-size: 32rpx;
}

.msg-dot {
  position: absolute;
  top: 0;
  right: 0;
  width: 16rpx;
  height: 16rpx;
  border-radius: 50%;
  background-color: #e64340;
}

.msg-content {
  flex: 1;
  margin-right: 16rpx;
}

.msg-title {
  font-size: 28rpx;
  color: #666;
  display: block;
  margin-bottom: 8rpx;
}

.msg-title.unread {
  color: #333;
  font-weight: 500;
}

.msg-brief {
  font-size: 24rpx;
  color: #999;
  display: block;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.msg-time {
  font-size: 22rpx;
  color: #ccc;
  flex-shrink: 0;
}
</style>
