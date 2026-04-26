<template>
  <view class="dealer-apply">
    <!-- Already applied -->
    <view v-if="applicationStatus" class="status-card">
      <uni-icons
        :type="applicationStatus === 'approved' ? 'checkmarkempty' : 'info'"
        size="48"
        :color="applicationStatus === 'approved' ? '#4caf50' : '#ff9800'"
      />
      <text class="status-text">{{ statusText }}</text>
      <text v-if="applicationStatus === 'approved'" class="status-hint">
        您已是分銷商，可前往分銷中心查看
      </text>
      <text v-else-if="applicationStatus === 'pending'" class="status-hint">
        您的申請正在審核中，請耐心等待
      </text>
      <text v-else-if="applicationStatus === 'rejected'" class="status-hint">
        {{ rejectReason || '您的申請未通過，可重新申請' }}
      </text>
      <view
        v-if="applicationStatus === 'approved'"
        class="submit-btn"
        @click="goCenter"
      >
        <text>進入分銷中心</text>
      </view>
      <view
        v-if="applicationStatus === 'rejected'"
        class="submit-btn"
        @click="applicationStatus = ''"
      >
        <text>重新申請</text>
      </view>
    </view>

    <!-- Application form -->
    <view v-else class="form-card">
      <text class="card-title">{{ t('dealer.apply') }}</text>

      <view class="form-item">
        <text class="form-label">真實姓名 <text class="required">*</text></text>
        <input
          class="form-input"
          placeholder="請輸入真實姓名"
          v-model="form.real_name"
        />
      </view>

      <view class="form-item">
        <text class="form-label">手機號 <text class="required">*</text></text>
        <input
          class="form-input"
          type="number"
          placeholder="請輸入手機號碼"
          v-model="form.phone"
          maxlength="11"
        />
      </view>

      <view class="form-item">
        <text class="form-label">微信號</text>
        <input
          class="form-input"
          placeholder="請輸入微信號"
          v-model="form.wechat"
        />
      </view>

      <view class="form-item">
        <text class="form-label">申請理由 <text class="required">*</text></text>
        <textarea
          class="form-textarea"
          placeholder="請說明您的申請理由、推廣資源等"
          v-model="form.reason"
          maxlength="500"
          :auto-height="false"
        />
        <text class="char-count">{{ form.reason.length }}/500</text>
      </view>

      <!-- Agreement -->
      <view class="agreement-row" @click="agreed = !agreed">
        <view class="checkbox" :class="{ checked: agreed }">
          <uni-icons v-if="agreed" type="checkmarkempty" size="14" color="#fff" />
        </view>
        <text class="agreement-text">
          我已閱讀並同意《分銷商合作協議》
        </text>
      </view>

      <view class="submit-btn" :class="{ disabled: !canSubmit }" @click="handleSubmit">
        <text>{{ t('common.submit') }}</text>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { t } from '@/locale'
import { commonApi } from '@/api/common'

const form = reactive({
  real_name: '',
  phone: '',
  wechat: '',
  reason: '',
})

const agreed = ref(false)
const applicationStatus = ref('')
const rejectReason = ref('')
const submitting = ref(false)

const statusText = computed(() => {
  const map: Record<string, string> = {
    pending: '審核中',
    approved: '已通過',
    rejected: '已拒絕',
  }
  return map[applicationStatus.value] || ''
})

const canSubmit = computed(() => {
  return (
    form.real_name.trim() &&
    form.phone.trim() &&
    form.reason.trim() &&
    agreed.value &&
    !submitting.value
  )
})

onMounted(() => {
  checkStatus()
})

async function checkStatus() {
  try {
    const res = await commonApi.getDealerInfo()
    const data = res?.data || {}
    if (data.status) {
      applicationStatus.value = data.status
      rejectReason.value = data.reject_reason || ''
    }
  } catch (e) {
    // Not applied yet
  }
}

async function handleSubmit() {
  if (!canSubmit.value) return
  submitting.value = true
  uni.showLoading({ title: t('common.loading') })
  try {
    await commonApi.applyDealer({
      real_name: form.real_name.trim(),
      phone: form.phone.trim(),
      wechat: form.wechat.trim(),
      reason: form.reason.trim(),
    })
    uni.showToast({ title: '申請已提交', icon: 'success' })
    applicationStatus.value = 'pending'
  } catch (e) {
    console.error(e)
    uni.showToast({ title: '提交失敗', icon: 'none' })
  } finally {
    submitting.value = false
    uni.hideLoading()
  }
}

function goCenter() {
  uni.redirectTo({ url: '/pages/dealer/index' })
}
</script>

<style scoped>
.dealer-apply {
  min-height: 100vh;
  background: #f5f5f5;
  padding: 24rpx;
}

.status-card {
  background: #fff;
  border-radius: 16rpx;
  padding: 60rpx 30rpx;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 20rpx;
}

.status-text {
  font-size: 34rpx;
  font-weight: 600;
  color: #333;
}

.status-hint {
  font-size: 26rpx;
  color: #999;
  text-align: center;
}

.form-card {
  background: #fff;
  border-radius: 16rpx;
  padding: 30rpx;
}

.card-title {
  font-size: 32rpx;
  font-weight: 600;
  color: #333;
  margin-bottom: 30rpx;
  display: block;
}

.form-item {
  margin-bottom: 28rpx;
  position: relative;
}

.form-label {
  display: block;
  font-size: 26rpx;
  color: #666;
  margin-bottom: 12rpx;
}

.required {
  color: #e64340;
}

.form-input {
  padding: 18rpx 20rpx;
  background: #f8f8f8;
  border-radius: 12rpx;
  font-size: 28rpx;
  color: #333;
}

.form-textarea {
  width: 100%;
  height: 200rpx;
  padding: 18rpx 20rpx;
  background: #f8f8f8;
  border-radius: 12rpx;
  font-size: 28rpx;
  color: #333;
  box-sizing: border-box;
}

.char-count {
  position: absolute;
  right: 16rpx;
  bottom: 12rpx;
  font-size: 22rpx;
  color: #bbb;
}

.agreement-row {
  display: flex;
  align-items: center;
  gap: 12rpx;
  margin: 30rpx 0;
}

.checkbox {
  width: 36rpx;
  height: 36rpx;
  border: 2rpx solid #ddd;
  border-radius: 6rpx;
  display: flex;
  align-items: center;
  justify-content: center;
}

.checkbox.checked {
  background: #0f3a57;
  border-color: #0f3a57;
}

.agreement-text {
  font-size: 24rpx;
  color: #666;
}

.submit-btn {
  background: #0f3a57;
  border-radius: 44rpx;
  height: 88rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-top: 20rpx;
}

.submit-btn.disabled {
  opacity: 0.5;
}

.submit-btn text {
  color: #fff;
  font-size: 30rpx;
  font-weight: 500;
}
</style>
