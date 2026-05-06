<template>
  <view class="login-page">
    <!-- Logo -->
    <view class="logo-section">
      <image class="logo" src="/static/logo.svg" mode="aspectFit" />
      <text class="brand-name">国韵好运</text>
    </view>

    <!-- Form -->
    <view class="form-section">
      <text class="form-title">{{ isLogin ? t('login.title') : t('login.register') }}</text>

      <view class="input-group">
        <view class="input-item">
          <text class="input-icon">&#xe614;</text>
          <input
            v-model="form.phone"
            type="number"
            :placeholder="t('login.phonePlaceholder')"
            maxlength="11"
            class="input"
          />
        </view>
        <view class="input-item">
          <text class="input-icon">&#xe618;</text>
          <input
            v-model="form.password"
            :password="true"
            :placeholder="t('login.passwordPlaceholder')"
            class="input"
          />
        </view>
        <view v-if="!isLogin && SHOW_SMS_CODE" class="input-item">
          <text class="input-icon">&#xe61a;</text>
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

      <!-- Login button -->
      <view class="submit-btn" @tap="handleSubmit">
        <text class="submit-btn-text">
          {{ isLogin ? t('login.loginBtn') : t('login.registerBtn') }}
        </text>
      </view>

      <!-- Links -->
      <view v-if="isLogin" class="links">
        <text class="link" @tap="navigateTo('/pages/login/index?mode=forgot')">
          {{ t('login.forgotPassword') }}
        </text>
      </view>

      <!-- Divider -->
      <view class="divider">
        <view class="divider-line" />
        <text class="divider-text">OR</text>
        <view class="divider-line" />
      </view>

      <!-- WeChat login -->
      <view class="wechat-btn" @tap="handleWechatLogin">
        <text class="wechat-icon">&#xe60e;</text>
        <text class="wechat-text">{{ t('login.wechatLogin') }}</text>
      </view>
    </view>

    <!-- Toggle mode -->
    <view class="toggle-section">
      <text class="toggle-text" @tap="toggleMode">
        {{ isLogin ? t('login.registerBtn') : t('login.loginBtn') }}
      </text>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { t } from '@/locale'
import { userApi } from '@/api/user'
import { setUser } from '@/store'

const isLogin = ref(true)
const loading = ref(false)
const countdown = ref(0)
const acceptedAgreement = ref(false)
// Toggle: SMS code temporarily disabled (provider not configured)
const SHOW_SMS_CODE = false
let timer: ReturnType<typeof setInterval> | null = null

const form = reactive({
  phone: '',
  password: '',
  code: '',
})

function toggleMode() {
  isLogin.value = !isLogin.value
  form.code = ''
}

function navigateTo(url: string) {
  uni.navigateTo({ url })
}

async function handleGetCode() {
  if (countdown.value > 0) return
  if (!form.phone || form.phone.length < 8) {
    uni.showToast({ title: t('login.phonePlaceholder'), icon: 'none' })
    return
  }
  try {
    await userApi.sendCode({ phone: form.phone, type: 'register' })
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

async function handleSubmit() {
  if (!acceptedAgreement.value) {
    uni.showToast({ title: '請先閱讀並同意用戶服務協議和隱私政策', icon: 'none' })
    return
  }
  if (!form.phone || !form.password) {
    uni.showToast({ title: '請填寫完整信息', icon: 'none' })
    return
  }
  if (loading.value) return
  loading.value = true

  try {
    if (isLogin.value) {
      const res = await userApi.login({ phone: form.phone, password: form.password })
      const payload = (res as any)?.data || res
      const token = payload?.token
      if (token) {
        uni.setStorageSync('token', token)
        if (payload.user) {
          uni.setStorageSync('userInfo', payload.user)
          setUser(payload.user as any)
        }
      } else {
        throw new Error('登入回應冇 token')
      }
      uni.showToast({ title: '登入成功', icon: 'success' })
      setTimeout(() => {
        uni.switchTab({ url: '/pages/index/index' })
      }, 1000)
    } else {
      if (SHOW_SMS_CODE && !form.code) {
        uni.showToast({ title: t('login.codePlaceholder'), icon: 'none' })
        loading.value = false
        return
      }
      await userApi.register({
        phone: form.phone,
        password: form.password,
        code: form.code || '000000',
      })
      uni.showToast({ title: '註冊成功', icon: 'success' })
      isLogin.value = true
    }
  } catch (e: any) {
    const msg = e?.message || e?.data?.error || (isLogin.value ? '登入失敗' : '註冊失敗')
    uni.showToast({ title: String(msg).slice(0, 30), icon: 'none' })
  } finally {
    loading.value = false
  }
}

async function handleWechatLogin() {
  if (!acceptedAgreement.value) {
    uni.showToast({ title: '請先閱讀並同意用戶服務協議和隱私政策', icon: 'none' })
    return
  }
  try {
    const [err, res] = await uni.login({ provider: 'weixin' }) as any
    if (err || !res?.code) {
      uni.showToast({ title: '微信授權失敗', icon: 'none' })
      return
    }
    const resp = (await userApi.wechatLogin({ code: res.code })) as any
    const data = resp?.data || resp
    if (data?.token) {
      uni.setStorageSync('token', data.token)
      if (data.user) {
        uni.setStorageSync('userInfo', data.user)
        setUser(data.user as any)
      }
      if (data.needBindPhone) {
        uni.navigateTo({ url: '/pages/login/bindPhone' })
      } else {
        uni.switchTab({ url: '/pages/index/index' })
      }
    }
  } catch {
    uni.showToast({ title: '微信登入失敗', icon: 'none' })
  }
}
</script>

<style scoped>
.login-page {
  min-height: 100vh;
  background: #f7f8fa;
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 0 60rpx;
  box-sizing: border-box;
}

.logo-section {
  display: flex;
  flex-direction: column;
  align-items: center;
  margin-top: 120rpx;
  margin-bottom: 60rpx;
}

.logo {
  width: 160rpx;
  height: 160rpx;
  border-radius: 24rpx;
}

.brand-name {
  font-size: 44rpx;
  font-weight: 700;
  color: #0f3a57;
  margin-top: 20rpx;
  letter-spacing: 2rpx;
}

.form-section {
  width: 100%;
}

.form-title {
  font-size: 40rpx;
  font-weight: 600;
  color: #333;
  margin-bottom: 40rpx;
}

.input-group {
  width: 100%;
}

.input-item {
  display: flex;
  align-items: center;
  background: #fff;
  border-radius: 16rpx;
  padding: 0 30rpx;
  height: 100rpx;
  margin-bottom: 24rpx;
  border: 2rpx solid #eee;
}

.input-icon {
  font-size: 36rpx;
  color: #999;
  margin-right: 20rpx;
  width: 40rpx;
  text-align: center;
}

.input {
  flex: 1;
  height: 100rpx;
  font-size: 30rpx;
  color: #333;
}

.code-input {
  flex: 1;
}

.code-btn {
  padding: 12rpx 24rpx;
  background: #0f3a57;
  border-radius: 8rpx;
  flex-shrink: 0;
}

.code-btn.disabled {
  background: #ccc;
}

.code-btn-text {
  font-size: 24rpx;
  color: #fff;
}

.submit-btn {
  width: 100%;
  height: 96rpx;
  background: #0f3a57;
  border-radius: 16rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-top: 16rpx;
}

.agreement-row {
  display: flex;
  align-items: flex-start;
  width: 100%;
  margin: 8rpx 0 28rpx;
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

.submit-btn-text {
  font-size: 32rpx;
  color: #fff;
  font-weight: 600;
}

.links {
  display: flex;
  justify-content: flex-end;
  margin-top: 20rpx;
}

.link {
  font-size: 26rpx;
  color: #0f3a57;
}

.divider {
  display: flex;
  align-items: center;
  margin: 50rpx 0 40rpx;
}

.divider-line {
  flex: 1;
  height: 1rpx;
  background: #ddd;
}

.divider-text {
  font-size: 24rpx;
  color: #999;
  padding: 0 30rpx;
}

.wechat-btn {
  width: 100%;
  height: 96rpx;
  background: #07c160;
  border-radius: 16rpx;
  display: flex;
  align-items: center;
  justify-content: center;
}

.wechat-icon {
  font-size: 40rpx;
  color: #fff;
  margin-right: 12rpx;
}

.wechat-text {
  font-size: 32rpx;
  color: #fff;
  font-weight: 600;
}

.toggle-section {
  margin-top: 40rpx;
}

.toggle-text {
  font-size: 28rpx;
  color: #0f3a57;
  text-decoration: underline;
}
</style>
