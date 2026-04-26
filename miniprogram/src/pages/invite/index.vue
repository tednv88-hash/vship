<template>
  <view class="invite-page">
    <!-- Loading -->
    <view v-if="loading" class="loading-wrap">
      <text class="loading-text">{{ t('common.loading') }}</text>
    </view>

    <view v-else class="invite-content">
      <!-- Invite code section -->
      <view class="code-section">
        <text class="code-label">Your Invite Code</text>
        <view class="code-display">
          <text class="code-text">{{ inviteInfo.code || '--' }}</text>
          <view class="copy-btn" @click="copyCode">
            <text class="copy-text">{{ t('common.copy') }}</text>
          </view>
        </view>
      </view>

      <!-- Share poster -->
      <view class="poster-section" v-if="posterUrl">
        <text class="section-title">Share Poster</text>
        <image class="poster-img" :src="posterUrl" mode="widthFix" @click="previewPoster" />
        <view class="poster-actions">
          <view class="poster-btn" @click="savePoster">
            <text class="poster-btn-text">Save to Album</text>
          </view>
        </view>
      </view>

      <!-- Reward rules -->
      <view class="rules-section">
        <text class="section-title">Reward Rules</text>
        <view class="rules-list">
          <view
            v-for="(rule, idx) in inviteInfo.rules"
            :key="idx"
            class="rule-item"
          >
            <text class="rule-number">{{ idx + 1 }}</text>
            <text class="rule-text">{{ rule }}</text>
          </view>
          <view v-if="!inviteInfo.rules || inviteInfo.rules.length === 0" class="rule-item">
            <text class="rule-number">1</text>
            <text class="rule-text">Invite friends to get reward points</text>
          </view>
        </view>
      </view>

      <!-- Invite records -->
      <view class="records-section">
        <text class="section-title">Invite Records</text>
        <view v-if="inviteInfo.records && inviteInfo.records.length > 0" class="records-list">
          <view
            v-for="record in inviteInfo.records"
            :key="record.id"
            class="record-item"
          >
            <view class="record-user">
              <image
                class="record-avatar"
                :src="record.avatar || '/static/avatar-default.png'"
                mode="aspectFill"
              />
              <text class="record-name">{{ record.nickname }}</text>
            </view>
            <view class="record-meta">
              <text class="record-reward">+{{ record.reward }}</text>
              <text class="record-date">{{ record.date }}</text>
            </view>
          </view>
        </view>
        <view v-else class="empty-records">
          <text class="empty-text">{{ t('common.noData') }}</text>
        </view>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { t } from '@/locale'
import { commonApi } from '@/api/common'

interface InviteRecord {
  id: string
  nickname: string
  avatar: string
  reward: string
  date: string
}

interface InviteInfo {
  code: string
  rules: string[]
  records: InviteRecord[]
}

const loading = ref(true)
const posterUrl = ref('')
const inviteInfo = ref<InviteInfo>({
  code: '',
  rules: [],
  records: [],
})

async function loadInviteInfo() {
  loading.value = true
  try {
    const [infoRes, posterRes]: any[] = await Promise.allSettled([
      commonApi.getInviteInfo(),
      commonApi.getInvitePoster(),
    ])

    if (infoRes.status === 'fulfilled') {
      const data = infoRes.value?.data || infoRes.value
      inviteInfo.value = {
        code: data.code || '',
        rules: data.rules || [],
        records: data.records || [],
      }
    }

    if (posterRes.status === 'fulfilled') {
      const data = posterRes.value?.data || posterRes.value
      posterUrl.value = data.url || data.poster_url || ''
    }
  } catch (e: any) {
    uni.showToast({ title: e?.message || t('common.retry'), icon: 'none' })
  } finally {
    loading.value = false
  }
}

function copyCode() {
  if (!inviteInfo.value.code) return
  uni.setClipboardData({
    data: inviteInfo.value.code,
    success: () => {
      uni.showToast({ title: t('common.copied'), icon: 'success' })
    },
  })
}

function previewPoster() {
  if (!posterUrl.value) return
  uni.previewImage({
    current: posterUrl.value,
    urls: [posterUrl.value],
  })
}

function savePoster() {
  if (!posterUrl.value) return
  uni.downloadFile({
    url: posterUrl.value,
    success: (downloadRes) => {
      if (downloadRes.statusCode === 200) {
        uni.saveImageToPhotosAlbum({
          filePath: downloadRes.tempFilePath,
          success: () => {
            uni.showToast({ title: t('common.done'), icon: 'success' })
          },
          fail: () => {
            uni.showToast({ title: 'Save failed', icon: 'none' })
          },
        })
      }
    },
  })
}

onMounted(() => {
  uni.setNavigationBarTitle({ title: t('invite.title') })
  loadInviteInfo()
})
</script>

<style scoped>
.invite-page {
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

.code-section {
  background-color: #0f3a57;
  padding: 48rpx 32rpx;
  text-align: center;
}

.code-label {
  font-size: 28rpx;
  color: rgba(255, 255, 255, 0.7);
  display: block;
  margin-bottom: 20rpx;
}

.code-display {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 24rpx;
}

.code-text {
  font-size: 48rpx;
  font-weight: 700;
  color: #fff;
  letter-spacing: 8rpx;
}

.copy-btn {
  padding: 12rpx 28rpx;
  background-color: rgba(255, 255, 255, 0.2);
  border-radius: 24rpx;
  border: 1rpx solid rgba(255, 255, 255, 0.5);
}

.copy-text {
  font-size: 24rpx;
  color: #fff;
}

.section-title {
  font-size: 30rpx;
  font-weight: 600;
  color: #333;
  display: block;
  margin-bottom: 20rpx;
}

.poster-section {
  background-color: #fff;
  margin: 20rpx 24rpx;
  border-radius: 16rpx;
  padding: 24rpx;
}

.poster-img {
  width: 100%;
  border-radius: 12rpx;
}

.poster-actions {
  display: flex;
  justify-content: center;
  margin-top: 20rpx;
}

.poster-btn {
  padding: 16rpx 48rpx;
  background-color: #0f3a57;
  border-radius: 44rpx;
}

.poster-btn-text {
  font-size: 28rpx;
  color: #fff;
}

.rules-section {
  background-color: #fff;
  margin: 20rpx 24rpx;
  border-radius: 16rpx;
  padding: 24rpx;
}

.rules-list {
  display: flex;
  flex-direction: column;
  gap: 16rpx;
}

.rule-item {
  display: flex;
  align-items: flex-start;
}

.rule-number {
  width: 40rpx;
  height: 40rpx;
  border-radius: 50%;
  background-color: #0f3a57;
  color: #fff;
  font-size: 22rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: 16rpx;
  flex-shrink: 0;
}

.rule-text {
  font-size: 26rpx;
  color: #666;
  line-height: 1.6;
  flex: 1;
}

.records-section {
  background-color: #fff;
  margin: 20rpx 24rpx;
  border-radius: 16rpx;
  padding: 24rpx;
}

.record-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16rpx 0;
  border-bottom: 1rpx solid #f0f0f0;
}

.record-item:last-child {
  border-bottom: none;
}

.record-user {
  display: flex;
  align-items: center;
}

.record-avatar {
  width: 56rpx;
  height: 56rpx;
  border-radius: 50%;
  margin-right: 16rpx;
}

.record-name {
  font-size: 28rpx;
  color: #333;
}

.record-meta {
  text-align: right;
}

.record-reward {
  font-size: 28rpx;
  color: #e64340;
  font-weight: 500;
  display: block;
}

.record-date {
  font-size: 22rpx;
  color: #999;
  margin-top: 4rpx;
  display: block;
}

.empty-records {
  padding: 40rpx 0;
  text-align: center;
}

.empty-text {
  font-size: 26rpx;
  color: #999;
}
</style>
