<template>
  <view class="estimate-page">
    <view class="form-card">
      <text class="card-title">{{ t('estimate.title') }}</text>

      <!-- Origin country -->
      <view class="form-item">
        <text class="form-label">發貨國家/地區</text>
        <picker :range="countryNames" @change="onCountryChange">
          <view class="picker-value">
            <text :class="{ placeholder: !form.origin }">{{ form.origin || '請選擇發貨地' }}</text>
            <uni-icons type="right" size="14" color="#999" />
          </view>
        </picker>
      </view>

      <!-- Destination -->
      <view class="form-item">
        <text class="form-label">收貨地區</text>
        <picker :range="destinationNames" @change="onDestinationChange">
          <view class="picker-value">
            <text :class="{ placeholder: !form.destination }">{{ form.destination || '請選擇收貨地' }}</text>
            <uni-icons type="right" size="14" color="#999" />
          </view>
        </picker>
      </view>

      <!-- Shipping route -->
      <view class="form-item">
        <text class="form-label">運輸線路</text>
        <picker :range="routeNames" @change="onRouteChange">
          <view class="picker-value">
            <text :class="{ placeholder: !form.route }">{{ form.route || '請選擇線路' }}</text>
            <uni-icons type="right" size="14" color="#999" />
          </view>
        </picker>
      </view>

      <!-- Weight -->
      <view class="form-item">
        <text class="form-label">重量 (kg)</text>
        <input
          class="form-input"
          type="digit"
          placeholder="請輸入重量"
          v-model="form.weight"
        />
      </view>

      <!-- Dimensions -->
      <view class="form-item">
        <text class="form-label">尺寸 (cm)</text>
        <view class="dimension-row">
          <input class="dim-input" type="digit" placeholder="長" v-model="form.length" />
          <text class="dim-x">×</text>
          <input class="dim-input" type="digit" placeholder="寬" v-model="form.width" />
          <text class="dim-x">×</text>
          <input class="dim-input" type="digit" placeholder="高" v-model="form.height" />
        </view>
      </view>

      <!-- Calculate button -->
      <view class="calc-btn" @click="calculate">
        <text>開始估價</text>
      </view>
    </view>

    <!-- Result -->
    <view v-if="result" class="result-card">
      <text class="card-title">估價結果</text>

      <view class="result-row">
        <text class="result-label">實際重量</text>
        <text class="result-value">{{ result.actual_weight }} kg</text>
      </view>
      <view class="result-row">
        <text class="result-label">體積重量</text>
        <text class="result-value">{{ result.volume_weight }} kg</text>
      </view>
      <view class="result-row">
        <text class="result-label">計費重量</text>
        <text class="result-value highlight">{{ result.charge_weight }} kg</text>
      </view>

      <view class="divider" />

      <view class="result-row">
        <text class="result-label">運費</text>
        <text class="result-value">¥{{ result.shipping_fee }}</text>
      </view>
      <view v-if="result.surcharge" class="result-row">
        <text class="result-label">附加費</text>
        <text class="result-value">¥{{ result.surcharge }}</text>
      </view>
      <view v-if="result.insurance_fee" class="result-row">
        <text class="result-label">保險費</text>
        <text class="result-value">¥{{ result.insurance_fee }}</text>
      </view>

      <view class="divider" />

      <view class="result-row total">
        <text class="result-label">預估總費用</text>
        <text class="result-value total-price">¥{{ result.total }}</text>
      </view>

      <view class="result-tip">
        <text>* 以上費用僅供參考，實際費用以倉庫稱重為準</text>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { t } from '@/locale'
import { commonApi } from '@/api/common'

const countries = ['日本', '韓國', '美國', '澳洲', '德國', '英國']
const destinations = ['台灣', '香港', '中國大陸', '澳門', '新加坡', '馬來西亞']
const routes = ['海運經濟', '海運標準', '空運標準', '空運快速', 'EMS特快']
const countryNames = ref(countries)
const destinationNames = ref(destinations)
const routeNames = ref(routes)

const form = reactive({
  origin: '',
  destination: '',
  route: '',
  weight: '',
  length: '',
  width: '',
  height: '',
})

const result = ref<any>(null)

function onCountryChange(e: any) {
  form.origin = countries[e.detail.value]
}

function onDestinationChange(e: any) {
  form.destination = destinations[e.detail.value]
}

function onRouteChange(e: any) {
  form.route = routes[e.detail.value]
}

async function calculate() {
  if (!form.origin) {
    uni.showToast({ title: '請選擇發貨國家', icon: 'none' })
    return
  }
  if (!form.destination) {
    uni.showToast({ title: '請選擇收貨地區', icon: 'none' })
    return
  }
  if (!form.weight) {
    uni.showToast({ title: '請輸入重量', icon: 'none' })
    return
  }

  uni.showLoading({ title: t('common.loading') })
  try {
    const res = await commonApi.calculateEstimate({
      origin: form.origin,
      destination: form.destination,
      route: form.route,
      weight: parseFloat(form.weight),
      length: parseFloat(form.length) || 0,
      width: parseFloat(form.width) || 0,
      height: parseFloat(form.height) || 0,
    })
    result.value = res?.data || res || {}
  } catch (e) {
    console.error(e)
    uni.showToast({ title: '估價失敗，請重試', icon: 'none' })
  } finally {
    uni.hideLoading()
  }
}
</script>

<style scoped>
.estimate-page {
  min-height: 100vh;
  background: #f5f5f5;
  padding: 24rpx;
}

.form-card,
.result-card {
  background: #fff;
  border-radius: 16rpx;
  padding: 30rpx;
  margin-bottom: 24rpx;
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
}

.form-label {
  display: block;
  font-size: 26rpx;
  color: #666;
  margin-bottom: 12rpx;
}

.picker-value {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 18rpx 20rpx;
  background: #f8f8f8;
  border-radius: 12rpx;
  font-size: 28rpx;
  color: #333;
}

.placeholder {
  color: #bbb;
}

.form-input {
  padding: 18rpx 20rpx;
  background: #f8f8f8;
  border-radius: 12rpx;
  font-size: 28rpx;
  color: #333;
}

.dimension-row {
  display: flex;
  align-items: center;
  gap: 12rpx;
}

.dim-input {
  flex: 1;
  padding: 18rpx 16rpx;
  background: #f8f8f8;
  border-radius: 12rpx;
  font-size: 28rpx;
  text-align: center;
}

.dim-x {
  font-size: 28rpx;
  color: #999;
}

.calc-btn {
  margin-top: 20rpx;
  background: #0f3a57;
  border-radius: 44rpx;
  height: 88rpx;
  display: flex;
  align-items: center;
  justify-content: center;
}

.calc-btn text {
  color: #fff;
  font-size: 30rpx;
  font-weight: 500;
}

.result-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16rpx 0;
}

.result-label {
  font-size: 28rpx;
  color: #666;
}

.result-value {
  font-size: 28rpx;
  color: #333;
  font-weight: 500;
}

.result-value.highlight {
  color: #0f3a57;
  font-weight: 600;
}

.divider {
  height: 1rpx;
  background: #f0f0f0;
  margin: 12rpx 0;
}

.result-row.total {
  padding-top: 20rpx;
}

.total-price {
  font-size: 40rpx;
  color: #e64340;
  font-weight: 700;
}

.result-tip {
  margin-top: 20rpx;
}

.result-tip text {
  font-size: 22rpx;
  color: #999;
}
</style>
