<template>
  <view class="policy-page">
    <!-- Loading -->
    <view v-if="loading" class="loading-wrap">
      <text class="loading-text">{{ t('common.loading') }}</text>
    </view>

    <!-- Error -->
    <view v-else-if="error" class="error-wrap">
      <text class="error-text">{{ error }}</text>
      <view class="retry-btn" @click="loadContent">
        <text class="retry-text">{{ t('common.retry') }}</text>
      </view>
    </view>

    <!-- Content -->
    <scroll-view v-else class="policy-scroll" scroll-y>
      <view class="policy-content">
        <rich-text :nodes="content" />
      </view>
    </scroll-view>
  </view>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { t } from '@/locale'
import { commonApi } from '@/api/common'

const loading = ref(true)
const error = ref('')
const content = ref('')

const fallbackContent = `
<div>
  <h2>用戶服務協議</h2>
  <p>歡迎使用國韻好運小程序。請您在登入、註冊或使用本服務前，仔細閱讀並充分理解本用戶服務協議。</p>
  <h3>一、服務內容</h3>
  <p>本小程序為用戶提供倉儲集運、包裹預報、訂單查詢、物流跟蹤、地址管理、支付及相關客戶服務。</p>
  <h3>二、帳號使用</h3>
  <p>您應提供真實、準確、完整的註冊及登入信息，並妥善保管帳號信息。因您保管不善造成的損失，由您自行承擔。</p>
  <h3>三、手機號收集與使用</h3>
  <p>為完成註冊登入、身份識別、訂單通知、物流聯繫及售後服務，我們需要收集並使用您的手機號碼。您未同意本協議及隱私政策前，我們不會進行登入或獲取手機號操作。</p>
  <h3>四、用戶行為規範</h3>
  <p>您不得利用本服務從事違法違規活動，不得提交虛假訂單、違規貨物信息或侵犯他人合法權益的內容。</p>
  <h3>五、服務變更與終止</h3>
  <p>我們可能根據業務需要優化或調整服務內容，並會依法保障您的合法權益。</p>
  <h3>六、聯繫我們</h3>
  <p>如您對本協議有任何疑問，可通過小程序客服或「關於我們」頁面提供的聯繫方式與我們聯繫。</p>
</div>
`

async function loadContent() {
  loading.value = true
  error.value = ''
  try {
    const res: any = await commonApi.getTerms()
    const data = res?.data || res
    content.value = data.content || data || fallbackContent
  } catch (e: any) {
    content.value = fallbackContent
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  uni.setNavigationBarTitle({ title: t('policy.terms') })
  loadContent()
})
</script>

<style scoped>
.policy-page {
  min-height: 100vh;
  background-color: #fff;
}

.policy-scroll {
  height: 100vh;
}

.loading-wrap,
.error-wrap {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 120rpx 0;
}

.loading-text {
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

.policy-content {
  padding: 32rpx;
  font-size: 28rpx;
  color: #333;
  line-height: 1.8;
}
</style>
