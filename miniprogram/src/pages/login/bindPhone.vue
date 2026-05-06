<template>
  <view class="bind-phone-page">
    <view class="header">
      <text class="title">{{ t('login.bindPhone') }}</text>
      <text class="subtitle">請綁定您的手機號碼以完成註冊</text>
    </view>

    <view class="form-section">
      <view class="input-item">
        <text class="label">{{ t('login.phone') }}</text>
        <input
          v-model="form.phone"
          type="number"
          :placeholder="t('login.phonePlaceholder')"
          maxlength="11"
          class="input"
        />
      </view>

      <view v-if="SHOW_SMS_CODE" class="input-item">
        <text class="label">驗證碼</text>
        <view class="code-row">
          <input
            v-model="form.code"
            type="number"
            :placeholder="t('login.codePlaceholder')"
            maxlength="6"
            class="input code-input"
          />
          <view
            class="code-btn"
            :class="{ disabled: countdown > 0 }"
            @tap="handleGetCode"
          >
            <text class="code-btn-text">
              {{ countdown > 0 ? `${countdown}s` : t('login.getCode') }}
            </text>
          </view>
        </view>
      </view>
    </view>

    <view class="confirm-btn" @tap="handleConfirm">
      <text class="confirm-btn-text">{{ t('common.confirm') }}</text>
    </view>

    <view class="agreement-bar">
      <view class="agreement-row" @tap="acceptedAgreement = !acceptedAgreement">
        <view class="checkbox" :class="{ checked: acceptedAgreement }">
          <text v-if="acceptedAgreement" class="check-text">√</text>
        </view>
        <view class="agreement-text-wrap">
          <text class="agreement-text">我已閱讀並同意</text>
          <text class="agreement-link" @tap.stop="navigateTo('/pages/policy/terms')">《用戶服務協議》</text>
          <text class="agreement-text">和</text>
          <text class="agreement-link" @tap.stop="navigateTo('/pages/policy/privacy')">《隱私政策》</text>
        </view>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { t } from '@/locale'
import { userApi } from '@/api/user'

const loading = ref(false)
const countdown = ref(0)
const acceptedAgreement = ref(false)
const SHOW_SMS_CODE = false
let timer: ReturnType<typeof setInterval> | null = null

const form = reactive({
  phone: '',
  code: '',
})

function navigateTo(url: string) {
  uni.navigateTo({ url })
}

async function handleGetCode() {
  if (!acceptedAgreement.value) {
    uni.showToast({ title: '請先閱讀並同意用戶服務協議和隱私政策', icon: 'none' })
    return
  }
  if (countdown.value > 0) return
  if (!form.phone || form.phone.length < 8) {
    uni.showToast({ title: t('login.phonePlaceholder'), icon: 'none' })
    return
  }
  try {
    await userApi.sendCode({ phone: form.phone, type: 'bind' })
    countdown.value = 60
    timer = setInterval(() => {
      countdown.value--
      if (countdown.value <= 0 && timer) {
        clearInterval(timer)
        timer = null
      }
    }, 1000)
    uni.showToast({ title: '驗證碼已發送', icon: 'success' })
  } catch {
    uni.showToast({ title: '發送失敗', icon: 'none' })
  }
}

async function handleConfirm() {
  if (!acceptedAgreement.value) {
    uni.showToast({ title: '請先閱讀並同意用戶服務協議和隱私政策', icon: 'none' })
    return
  }
  if (!form.phone || !form.code) {
    uni.showToast({ title: '請填寫完整信息', icon: 'none' })
    return
  }
  if (loading.value) return
  loading.value = true

  try {
    await userApi.bindPhone({ phone: form.phone, code: form.code })
    uni.showToast({ title: '綁定成功', icon: 'success' })
    setTimeout(() => {
      uni.switchTab({ url: '/pages/index/index' })
    }, 1000)
  } catch {
    uni.showToast({ title: '綁定失敗', icon: 'none' })
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.bind-phone-page {
  min-height: 100vh;
  background: #f7f8fa;
  padding: 0 40rpx 180rpx;
  box-sizing: border-box;
}

.header {
  padding-top: 120rpx;
  margin-bottom: 60rpx;
}

.title {
  display: block;
  font-size: 44rpx;
  font-weight: 700;
  color: #0f3a57;
  margin-bottom: 16rpx;
}

.subtitle {
  font-size: 28rpx;
  color: #999;
}

.form-section {
  width: 100%;
}

.input-item {
  margin-bottom: 36rpx;
}

.label {
  display: block;
  font-size: 28rpx;
  color: #333;
  font-weight: 500;
  margin-bottom: 16rpx;
}

.input {
  width: 100%;
  height: 96rpx;
  background: #fff;
  border-radius: 16rpx;
  padding: 0 30rpx;
  font-size: 30rpx;
  color: #333;
  border: 2rpx solid #eee;
  box-sizing: border-box;
}

.code-row {
  display: flex;
  align-items: center;
  gap: 20rpx;
}

.code-input {
  flex: 1;
}

.code-btn {
  flex-shrink: 0;
  height: 96rpx;
  padding: 0 32rpx;
  background: #0f3a57;
  border-radius: 16rpx;
  display: flex;
  align-items: center;
  justify-content: center;
}

.code-btn.disabled {
  background: #ccc;
}

.code-btn-text {
  font-size: 28rpx;
  color: #fff;
  white-space: nowrap;
}

.confirm-btn {
  width: 100%;
  height: 96rpx;
  background: #0f3a57;
  border-radius: 16rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-top: 20rpx;
}

.agreement-row {
  display: flex;
  align-items: flex-start;
  width: 100%;
}

.agreement-bar {
  position: fixed;
  left: 0;
  right: 0;
  bottom: 0;
  z-index: 20;
  background: #fff;
  border-top: 1rpx solid #eee;
  padding: 24rpx 40rpx calc(24rpx + env(safe-area-inset-bottom));
  box-shadow: 0 -8rpx 24rpx rgba(0, 0, 0, 0.06);
  box-sizing: border-box;
}

.checkbox {
  width: 32rpx;
  height: 32rpx;
  border: 2rpx solid #c8c9cc;
  border-radius: 50%;
  margin-right: 14rpx;
  margin-top: 4rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.checkbox.checked {
  background: #0f3a57;
  border-color: #0f3a57;
}

.check-text {
  color: #fff;
  font-size: 22rpx;
  line-height: 1;
}

.agreement-text-wrap {
  flex: 1;
  line-height: 40rpx;
}

.agreement-text,
.agreement-link {
  font-size: 28rpx;
}

.agreement-text {
  color: #666;
}

.agreement-link {
  color: #0f3a57;
}

.confirm-btn-text {
  font-size: 32rpx;
  color: #fff;
  font-weight: 600;
}
</style>
