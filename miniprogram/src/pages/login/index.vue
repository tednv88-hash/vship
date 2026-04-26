<template>
  <view class="login-page">
    <!-- Logo -->
    <view class="logo-section">
      <image class="logo" src="/static/logo.png" mode="aspectFit" />
      <text class="brand-name">vShip</text>
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
        <view v-if="!isLogin" class="input-item">
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

const isLogin = ref(true)
const loading = ref(false)
const countdown = ref(0)
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
  if (!form.phone || !form.password) {
    uni.showToast({ title: '請填寫完整信息', icon: 'none' })
    return
  }
  if (loading.value) return
  loading.value = true

  try {
    if (isLogin.value) {
      const res = await userApi.login({ phone: form.phone, password: form.password })
      const data = res as any
      if (data?.token) {
        uni.setStorageSync('token', data.token)
      }
      uni.showToast({ title: '登入成功', icon: 'success' })
      setTimeout(() => {
        uni.switchTab({ url: '/pages/index/index' })
      }, 1000)
    } else {
      if (!form.code) {
        uni.showToast({ title: t('login.codePlaceholder'), icon: 'none' })
        loading.value = false
        return
      }
      await userApi.register({
        phone: form.phone,
        password: form.password,
        code: form.code,
      })
      uni.showToast({ title: '註冊成功', icon: 'success' })
      isLogin.value = true
    }
  } catch {
    uni.showToast({ title: isLogin.value ? '登入失敗' : '註冊失敗', icon: 'none' })
  } finally {
    loading.value = false
  }
}

async function handleWechatLogin() {
  try {
    const [err, res] = await uni.login({ provider: 'weixin' }) as any
    if (err || !res?.code) {
      uni.showToast({ title: '微信授權失敗', icon: 'none' })
      return
    }
    const data = (await userApi.wechatLogin({ code: res.code })) as any
    if (data?.token) {
      uni.setStorageSync('token', data.token)
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
