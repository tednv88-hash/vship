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
  <h2>隱私政策</h2>
  <p>國韻好運重視並保護您的個人信息安全。請您在使用本小程序前仔細閱讀本隱私政策。</p>
  <h3>一、我們收集的信息</h3>
  <p>為向您提供登入註冊、手機號綁定、包裹預報、訂單處理、物流通知、售後客服等服務，我們可能收集您的手機號、微信授權登入信息、收貨地址、訂單信息、包裹信息、支付狀態及客服溝通信息。</p>
  <h3>二、手機號的使用目的</h3>
  <p>我們收集手機號僅用於帳號識別、登入註冊、訂單及物流通知、異常聯繫、售後服務及法律法規要求的必要場景。未勾選同意《用戶服務協議》和《隱私政策》前，不會登入或獲取手機號。</p>
  <h3>三、信息的存儲與保護</h3>
  <p>我們會採取合理的安全措施保護您的個人信息，防止未經授權訪問、披露、使用、修改、損壞或丟失。</p>
  <h3>四、信息共享</h3>
  <p>除為完成物流配送、支付結算、售後服務或依法依規需要外，我們不會向無關第三方共享您的個人信息。</p>
  <h3>五、用戶權利</h3>
  <p>您有權查詢、更正、刪除您的個人信息，或撤回授權。您可以通過小程序客服或「關於我們」頁面提供的聯繫方式提出申請。</p>
  <h3>六、政策更新</h3>
  <p>我們可能根據法律法規或業務變化更新本政策，更新後會在本頁面展示。</p>
  <h3>七、聯繫我們</h3>
  <p>如您對個人信息保護有任何疑問、投訴或建議，可通過小程序客服與我們聯繫。</p>
</div>
`

async function loadContent() {
  loading.value = true
  error.value = ''
  try {
    const res: any = await commonApi.getPrivacy()
    const data = res?.data || res
    content.value = data.content || data || fallbackContent
  } catch (e: any) {
    content.value = fallbackContent
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  uni.setNavigationBarTitle({ title: t('policy.privacy') })
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
