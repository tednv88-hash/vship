<template>
  <view class="insurance-page">
    <view class="page-desc">
      <text>為您的包裹選擇合適的保險方案，讓運輸更安心</text>
    </view>

    <!-- Insurance cards -->
    <view class="insurance-list">
      <view
        v-for="option in insuranceList"
        :key="option.id"
        class="insurance-card"
        :class="{ selected: selectedId === option.id }"
        @click="selectInsurance(option.id)"
      >
        <!-- Badge -->
        <view v-if="option.recommended" class="recommend-badge">
          <text>推薦</text>
        </view>

        <text class="tier-name">{{ option.name }}</text>

        <view class="coverage-row">
          <text class="coverage-label">保障金額</text>
          <text class="coverage-amount">¥{{ option.coverage }}</text>
        </view>

        <view class="premium-row">
          <text class="premium-label">保費</text>
          <text class="premium-amount">¥{{ option.premium }}</text>
        </view>

        <view class="divider" />

        <text class="coverage-desc">{{ option.description }}</text>

        <view class="coverage-details">
          <view v-for="(item, idx) in option.details" :key="idx" class="detail-item">
            <uni-icons type="checkmarkempty" size="14" color="#4caf50" />
            <text>{{ item }}</text>
          </view>
        </view>

        <view class="select-btn" :class="{ active: selectedId === option.id }">
          <text>{{ selectedId === option.id ? '已選擇' : '選擇方案' }}</text>
        </view>
      </view>
    </view>

    <!-- Empty -->
    <view v-if="!loading && insuranceList.length === 0" class="empty">
      <text>{{ t('common.noData') }}</text>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { t } from '@/locale'
import { commonApi } from '@/api/common'

const insuranceList = ref<any[]>([])
const selectedId = ref('')
const loading = ref(true)

const defaultOptions = [
  {
    id: '1',
    name: '基礎保障',
    coverage: '500',
    premium: '3.00',
    description: '適合低價值普通商品，提供基本運輸損壞保障。',
    recommended: false,
    details: ['運輸損壞賠付', '包裹丟失賠付', '最高賠付¥500'],
  },
  {
    id: '2',
    name: '標準保障',
    coverage: '2,000',
    premium: '10.00',
    description: '適合中等價值商品，覆蓋常見運輸風險。',
    recommended: true,
    details: ['運輸損壞賠付', '包裹丟失賠付', '海關扣留賠付', '最高賠付¥2,000'],
  },
  {
    id: '3',
    name: '全面保障',
    coverage: '5,000',
    premium: '25.00',
    description: '適合高價值商品，提供最全面的保障方案。',
    recommended: false,
    details: ['運輸損壞全額賠付', '包裹丟失全額賠付', '海關扣留賠付', '延遲送達補償', '最高賠付¥5,000'],
  },
]

onMounted(() => {
  loadInsurance()
})

async function loadInsurance() {
  loading.value = true
  try {
    const res = await commonApi.getInsuranceOptions()
    const list = res?.data?.list || res?.data || []
    insuranceList.value = list.length > 0 ? list : defaultOptions
  } catch (e) {
    insuranceList.value = defaultOptions
    console.error(e)
  } finally {
    loading.value = false
  }
}

function selectInsurance(id: string) {
  selectedId.value = selectedId.value === id ? '' : id
}
</script>

<style scoped>
.insurance-page {
  min-height: 100vh;
  background: #f5f5f5;
  padding: 24rpx;
}

.page-desc {
  text-align: center;
  padding: 20rpx 0 30rpx;
}

.page-desc text {
  font-size: 26rpx;
  color: #999;
}

.insurance-list {
  display: flex;
  flex-direction: column;
  gap: 24rpx;
}

.insurance-card {
  background: #fff;
  border-radius: 16rpx;
  padding: 30rpx;
  border: 2rpx solid transparent;
  position: relative;
  overflow: hidden;
}

.insurance-card.selected {
  border-color: #0f3a57;
}

.recommend-badge {
  position: absolute;
  top: 0;
  right: 0;
  background: #e64340;
  padding: 6rpx 20rpx;
  border-radius: 0 14rpx 0 16rpx;
}

.recommend-badge text {
  font-size: 22rpx;
  color: #fff;
}

.tier-name {
  font-size: 34rpx;
  color: #333;
  font-weight: 700;
  display: block;
  margin-bottom: 20rpx;
}

.coverage-row,
.premium-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 10rpx 0;
}

.coverage-label,
.premium-label {
  font-size: 26rpx;
  color: #999;
}

.coverage-amount {
  font-size: 36rpx;
  color: #0f3a57;
  font-weight: 700;
}

.premium-amount {
  font-size: 30rpx;
  color: #e64340;
  font-weight: 600;
}

.divider {
  height: 1rpx;
  background: #f0f0f0;
  margin: 20rpx 0;
}

.coverage-desc {
  font-size: 26rpx;
  color: #666;
  line-height: 1.6;
  display: block;
  margin-bottom: 16rpx;
}

.coverage-details {
  display: flex;
  flex-direction: column;
  gap: 12rpx;
  margin-bottom: 24rpx;
}

.detail-item {
  display: flex;
  align-items: center;
  gap: 10rpx;
}

.detail-item text {
  font-size: 24rpx;
  color: #666;
}

.select-btn {
  border: 2rpx solid #0f3a57;
  border-radius: 44rpx;
  height: 72rpx;
  display: flex;
  align-items: center;
  justify-content: center;
}

.select-btn text {
  font-size: 28rpx;
  color: #0f3a57;
}

.select-btn.active {
  background: #0f3a57;
}

.select-btn.active text {
  color: #fff;
}

.empty {
  text-align: center;
  padding: 100rpx 0;
  color: #999;
  font-size: 28rpx;
}
</style>
