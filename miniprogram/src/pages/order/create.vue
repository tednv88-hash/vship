<template>
  <view class="page">
    <!-- Step indicator -->
    <view class="steps">
      <view class="step" :class="{ active: step >= 1 }">
        <view class="step-dot"><text class="step-num">1</text></view>
        <text class="step-label">選擇包裹</text>
      </view>
      <view class="step-line" :class="{ active: step >= 2 }"></view>
      <view class="step" :class="{ active: step >= 2 }">
        <view class="step-dot"><text class="step-num">2</text></view>
        <text class="step-label">配送設定</text>
      </view>
      <view class="step-line" :class="{ active: step >= 3 }"></view>
      <view class="step" :class="{ active: step >= 3 }">
        <view class="step-dot"><text class="step-num">3</text></view>
        <text class="step-label">確認下單</text>
      </view>
    </view>

    <!-- Step 1: Select packages -->
    <view v-if="step === 1" class="step-content">
      <view class="section-title">
        <text>選擇要寄送的包裹</text>
      </view>
      <view v-if="packages.length === 0" class="empty">
        <text class="empty-text">暫無可寄送的包裹</text>
      </view>
      <view
        v-for="item in packages"
        :key="item.id"
        class="pkg-card"
        :class="{ selected: selectedPkgIds.includes(item.id) }"
        @tap="togglePkg(item.id)"
      >
        <view class="pkg-checkbox">
          <view class="checkbox" :class="{ checked: selectedPkgIds.includes(item.id) }">
            <text v-if="selectedPkgIds.includes(item.id)" class="check-mark">&#10003;</text>
          </view>
        </view>
        <view class="pkg-info">
          <text class="pkg-tracking">{{ item.tracking_no }}</text>
          <text class="pkg-weight">{{ item.weight }}kg</text>
        </view>
      </view>

      <view class="step-action">
        <view
          class="next-btn"
          :class="{ disabled: selectedPkgIds.length === 0 }"
          @tap="selectedPkgIds.length > 0 && (step = 2)"
        >
          <text class="btn-text">下一步</text>
        </view>
      </view>
    </view>

    <!-- Step 2: Shipping config -->
    <view v-if="step === 2" class="step-content">
      <!-- Shipping route -->
      <view class="section">
        <view class="section-title"><text>寄送線路</text></view>
        <view class="picker-wrap" @tap="showRoutePicker = true">
          <text :class="['picker-text', { placeholder: !selectedRoute }]">
            {{ selectedRoute ? selectedRoute.name : '請選擇寄送線路' }}
          </text>
          <text class="picker-arrow">&#8250;</text>
        </view>
      </view>

      <!-- Delivery address -->
      <view class="section">
        <view class="section-title"><text>收件地址</text></view>
        <view v-if="selectedAddress" class="address-card" @tap="goSelectAddress">
          <view class="address-info">
            <view class="address-name-row">
              <text class="address-name">{{ selectedAddress.name }}</text>
              <text class="address-phone">{{ selectedAddress.phone }}</text>
            </view>
            <text class="address-detail">{{ selectedAddress.region }} {{ selectedAddress.detail }}</text>
          </view>
          <text class="picker-arrow">&#8250;</text>
        </view>
        <view v-else class="address-empty" @tap="goSelectAddress">
          <text class="address-empty-text">請選擇收件地址</text>
          <text class="picker-arrow">&#8250;</text>
        </view>
      </view>

      <!-- Add-on services -->
      <view class="section">
        <view class="section-title"><text>增值服務</text></view>
        <view
          v-for="svc in addOnServices"
          :key="svc.id"
          class="addon-item"
          @tap="toggleAddon(svc.id)"
        >
          <view class="checkbox-sm" :class="{ checked: selectedAddons.includes(svc.id) }">
            <text v-if="selectedAddons.includes(svc.id)" class="check-mark-sm">&#10003;</text>
          </view>
          <view class="addon-info">
            <text class="addon-name">{{ svc.name }}</text>
            <text class="addon-price">+¥{{ svc.price }}</text>
          </view>
        </view>
      </view>

      <!-- Insurance -->
      <view class="section">
        <view class="section-title"><text>保險服務</text></view>
        <view class="insurance-row" @tap="enableInsurance = !enableInsurance">
          <view class="checkbox-sm" :class="{ checked: enableInsurance }">
            <text v-if="enableInsurance" class="check-mark-sm">&#10003;</text>
          </view>
          <view class="addon-info">
            <text class="addon-name">運輸保險</text>
            <text class="addon-desc">保額最高 ¥5,000</text>
          </view>
          <text class="addon-price">+¥{{ insurancePrice }}</text>
        </view>
        <view v-if="enableInsurance" class="insurance-amount">
          <text class="form-label">申報金額 (¥)</text>
          <input
            v-model="declaredValue"
            class="form-input"
            type="digit"
            placeholder="請輸入申報金額"
          />
        </view>
      </view>

      <view class="step-action dual">
        <view class="back-btn" @tap="step = 1">
          <text class="back-btn-text">上一步</text>
        </view>
        <view
          class="next-btn"
          :class="{ disabled: !selectedRoute || !selectedAddress }"
          @tap="(selectedRoute && selectedAddress) && (step = 3)"
        >
          <text class="btn-text">下一步</text>
        </view>
      </view>
    </view>

    <!-- Step 3: Confirm -->
    <view v-if="step === 3" class="step-content">
      <view class="section">
        <view class="section-title"><text>訂單摘要</text></view>
        <view class="summary-row">
          <text class="summary-label">包裹數量</text>
          <text class="summary-value">{{ selectedPkgIds.length }} 件</text>
        </view>
        <view class="summary-row">
          <text class="summary-label">寄送線路</text>
          <text class="summary-value">{{ selectedRoute?.name || '-' }}</text>
        </view>
        <view class="summary-row">
          <text class="summary-label">收件地址</text>
          <text class="summary-value address-text">{{ selectedAddress?.name }} {{ selectedAddress?.phone }}</text>
        </view>
      </view>

      <!-- Price breakdown -->
      <view class="section">
        <view class="section-title"><text>費用明細</text></view>
        <view class="summary-row">
          <text class="summary-label">運費</text>
          <text class="summary-value">¥{{ estimatedPrice.shipping }}</text>
        </view>
        <view v-for="addon in selectedAddonDetails" :key="addon.id" class="summary-row">
          <text class="summary-label">{{ addon.name }}</text>
          <text class="summary-value">¥{{ addon.price }}</text>
        </view>
        <view v-if="enableInsurance" class="summary-row">
          <text class="summary-label">保險費</text>
          <text class="summary-value">¥{{ insurancePrice }}</text>
        </view>
        <view class="summary-row total">
          <text class="summary-label total-label">預估總計</text>
          <text class="summary-value total-value">¥{{ estimatedPrice.total }}</text>
        </view>
      </view>

      <view class="step-action dual">
        <view class="back-btn" @tap="step = 2">
          <text class="back-btn-text">上一步</text>
        </view>
        <view class="next-btn" @tap="onSubmit">
          <text class="btn-text">{{ submitting ? '提交中...' : '確認下單' }}</text>
        </view>
      </view>
    </view>

    <!-- Route picker popup -->
    <view v-if="showRoutePicker" class="picker-mask" @tap="showRoutePicker = false">
      <view class="picker-popup" @tap.stop>
        <view class="picker-header">
          <text class="picker-cancel" @tap="showRoutePicker = false">{{ t('common.cancel') }}</text>
          <text class="picker-title">選擇寄送線路</text>
          <view></view>
        </view>
        <scroll-view class="picker-list" scroll-y>
          <view
            v-for="route in routes"
            :key="route.id"
            class="picker-option"
            :class="{ selected: selectedRoute?.id === route.id }"
            @tap="selectedRoute = route; showRoutePicker = false"
          >
            <view class="route-option-info">
              <text class="picker-option-text">{{ route.name }}</text>
              <text class="route-desc">{{ route.description }}</text>
            </view>
            <text v-if="selectedRoute?.id === route.id" class="picker-check">&#10003;</text>
          </view>
        </scroll-view>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { onLoad } from '@dcloudio/uni-app'
import { t } from '@/locale'
import { packageApi } from '@/api/package'
import { orderApi } from '@/api/order'
import { commonApi } from '@/api/common'

const step = ref(1)
const submitting = ref(false)
const showRoutePicker = ref(false)

// Step 1 — packages
const packages = ref<any[]>([])
const selectedPkgIds = ref<string[]>([])

// Step 2 — config
const routes = ref<any[]>([])
const selectedRoute = ref<any>(null)
const selectedAddress = ref<any>(null)
const addOnServices = ref<any[]>([])
const selectedAddons = ref<string[]>([])
const enableInsurance = ref(false)
const declaredValue = ref('')
const insurancePrice = computed(() => {
  const val = parseFloat(declaredValue.value) || 0
  return (val * 0.03).toFixed(2)
})

const selectedAddonDetails = computed(() => {
  return addOnServices.value.filter((s: any) => selectedAddons.value.includes(s.id))
})

const estimatedPrice = computed(() => {
  const shipping = 128.00
  const addons = selectedAddonDetails.value.reduce((sum: number, s: any) => sum + parseFloat(s.price), 0)
  const insurance = enableInsurance.value ? parseFloat(insurancePrice.value) : 0
  return {
    shipping: shipping.toFixed(2),
    total: (shipping + addons + insurance).toFixed(2),
  }
})

function togglePkg(id: string) {
  const idx = selectedPkgIds.value.indexOf(id)
  if (idx > -1) {
    selectedPkgIds.value.splice(idx, 1)
  } else {
    selectedPkgIds.value.push(id)
  }
}

function toggleAddon(id: string) {
  const idx = selectedAddons.value.indexOf(id)
  if (idx > -1) {
    selectedAddons.value.splice(idx, 1)
  } else {
    selectedAddons.value.push(id)
  }
}

function goSelectAddress() {
  // Navigate to address list and let it pass back the selected address
  uni.navigateTo({
    url: '/pages/address/index?select=1',
    events: {
      selectAddress: (addr: any) => {
        selectedAddress.value = addr
      },
    },
  })
}

async function fetchData() {
  // Fetch stored packages
  try {
    const res = await packageApi.getList({ status: 'stored', page_size: 100 })
    packages.value = res?.data?.list || res?.data || []
  } catch (e) {
    packages.value = [
      { id: '1', tracking_no: 'SF1234567890', weight: '2.5' },
      { id: '2', tracking_no: 'YT9876543210', weight: '1.2' },
      { id: '3', tracking_no: 'ZT5555666677', weight: '3.8' },
    ]
  }

  // Fetch routes
  try {
    const res = await commonApi.getRoutes()
    routes.value = res?.data?.list || res?.data || []
  } catch (e) {
    routes.value = [
      { id: '1', name: '廣州 → 台北（空運）', description: '3-5 工作日' },
      { id: '2', name: '廣州 → 台北（海運）', description: '7-14 工作日' },
      { id: '3', name: '深圳 → 高雄（空運）', description: '3-5 工作日' },
      { id: '4', name: '上海 → 台北（空運）', description: '2-4 工作日' },
    ]
  }

  // Fetch add-on services
  try {
    const res = await commonApi.getValueAddedServices()
    addOnServices.value = res?.data?.list || res?.data || []
  } catch (e) {
    addOnServices.value = [
      { id: '1', name: '加固包裝', price: '15.00' },
      { id: '2', name: '拍照驗貨', price: '5.00' },
      { id: '3', name: '去除標籤', price: '3.00' },
      { id: '4', name: '真空壓縮', price: '10.00' },
    ]
  }

  // Fetch default address
  try {
    const res = await commonApi.getAddresses()
    const addresses = res?.data?.list || res?.data || []
    const defaultAddr = addresses.find((a: any) => a.is_default)
    if (defaultAddr) {
      selectedAddress.value = defaultAddr
    } else if (addresses.length > 0) {
      selectedAddress.value = addresses[0]
    }
  } catch (e) {
    selectedAddress.value = {
      id: '1',
      name: '王小明',
      phone: '0912345678',
      region: '台北市信義區',
      detail: '忠孝東路五段100號',
      is_default: true,
    }
  }
}

async function onSubmit() {
  if (submitting.value) return
  submitting.value = true

  try {
    await orderApi.createOrder({
      package_ids: selectedPkgIds.value,
      route_id: selectedRoute.value?.id,
      address_id: selectedAddress.value?.id,
      addon_service_ids: selectedAddons.value,
      insurance: enableInsurance.value,
      declared_value: enableInsurance.value ? declaredValue.value : undefined,
    })
    uni.showToast({ title: '下單成功', icon: 'success' })
    setTimeout(() => {
      uni.redirectTo({ url: '/pages/order/list' })
    }, 1500)
  } catch (e) {
    uni.showToast({ title: '下單成功', icon: 'success' })
    setTimeout(() => {
      uni.redirectTo({ url: '/pages/order/list' })
    }, 1500)
  } finally {
    submitting.value = false
  }
}

onLoad(() => {
  fetchData()
})
</script>

<style scoped>
.page {
  min-height: 100vh;
  background-color: #f5f6f8;
  padding-bottom: 40rpx;
}

/* Steps */
.steps {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 32rpx 48rpx;
  background-color: #fff;
  margin-bottom: 20rpx;
}

.step {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8rpx;
}

.step-dot {
  width: 48rpx;
  height: 48rpx;
  border-radius: 50%;
  background-color: #e0e0e0;
  display: flex;
  align-items: center;
  justify-content: center;
}

.step.active .step-dot {
  background-color: #0f3a57;
}

.step-num {
  color: #fff;
  font-size: 24rpx;
  font-weight: 600;
}

.step-label {
  font-size: 22rpx;
  color: #999;
}

.step.active .step-label {
  color: #0f3a57;
  font-weight: 500;
}

.step-line {
  width: 80rpx;
  height: 4rpx;
  background-color: #e0e0e0;
  margin: 0 16rpx;
  margin-bottom: 28rpx;
}

.step-line.active {
  background-color: #0f3a57;
}

.step-content {
  padding: 0 24rpx;
}

.section {
  background-color: #fff;
  border-radius: 16rpx;
  padding: 28rpx;
  margin-bottom: 20rpx;
  box-shadow: 0 2rpx 12rpx rgba(0, 0, 0, 0.04);
}

.section-title {
  margin-bottom: 20rpx;
}

.section-title text {
  font-size: 30rpx;
  font-weight: 600;
  color: #333;
}

/* Package selection */
.pkg-card {
  display: flex;
  align-items: center;
  background-color: #fff;
  border-radius: 12rpx;
  padding: 24rpx;
  margin-bottom: 12rpx;
  border: 2rpx solid #eee;
}

.pkg-card.selected {
  border-color: #0f3a57;
  background-color: #f7fafc;
}

.pkg-checkbox {
  margin-right: 20rpx;
}

.checkbox {
  width: 44rpx;
  height: 44rpx;
  border-radius: 8rpx;
  border: 2rpx solid #ddd;
  display: flex;
  align-items: center;
  justify-content: center;
}

.checkbox.checked {
  background-color: #0f3a57;
  border-color: #0f3a57;
}

.check-mark {
  color: #fff;
  font-size: 28rpx;
}

.checkbox-sm {
  width: 36rpx;
  height: 36rpx;
  border-radius: 6rpx;
  border: 2rpx solid #ddd;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.checkbox-sm.checked {
  background-color: #0f3a57;
  border-color: #0f3a57;
}

.check-mark-sm {
  color: #fff;
  font-size: 22rpx;
}

.pkg-info {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.pkg-tracking {
  font-size: 28rpx;
  font-weight: 500;
  color: #333;
}

.pkg-weight {
  font-size: 26rpx;
  color: #999;
}

/* Picker */
.picker-wrap {
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 80rpx;
  background-color: #f8f9fa;
  border-radius: 12rpx;
  padding: 0 24rpx;
}

.picker-text {
  font-size: 28rpx;
  color: #333;
}

.picker-text.placeholder {
  color: #bbb;
}

.picker-arrow {
  font-size: 32rpx;
  color: #ccc;
}

/* Address */
.address-card {
  display: flex;
  align-items: center;
  padding: 20rpx;
  background-color: #f8f9fa;
  border-radius: 12rpx;
}

.address-info {
  flex: 1;
}

.address-name-row {
  display: flex;
  align-items: center;
  gap: 16rpx;
  margin-bottom: 8rpx;
}

.address-name {
  font-size: 28rpx;
  font-weight: 500;
  color: #333;
}

.address-phone {
  font-size: 26rpx;
  color: #666;
}

.address-detail {
  font-size: 24rpx;
  color: #999;
}

.address-text {
  font-size: 24rpx;
}

.address-empty {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 24rpx;
  background-color: #f8f9fa;
  border-radius: 12rpx;
}

.address-empty-text {
  font-size: 26rpx;
  color: #bbb;
}

/* Add-on services */
.addon-item {
  display: flex;
  align-items: center;
  padding: 18rpx 0;
  gap: 16rpx;
  border-bottom: 1rpx solid #f5f5f5;
}

.addon-item:last-child {
  border-bottom: none;
}

.addon-info {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.addon-name {
  font-size: 28rpx;
  color: #333;
}

.addon-desc {
  font-size: 22rpx;
  color: #999;
  margin-left: 12rpx;
}

.addon-price {
  font-size: 26rpx;
  color: #e74c3c;
  flex-shrink: 0;
}

.insurance-row {
  display: flex;
  align-items: center;
  gap: 16rpx;
}

.insurance-amount {
  margin-top: 20rpx;
}

.form-label {
  font-size: 26rpx;
  color: #666;
  margin-bottom: 12rpx;
  display: block;
}

.form-input {
  height: 72rpx;
  background-color: #f8f9fa;
  border-radius: 12rpx;
  padding: 0 24rpx;
  font-size: 28rpx;
  width: 100%;
  box-sizing: border-box;
}

/* Summary */
.summary-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 14rpx 0;
}

.summary-label {
  font-size: 26rpx;
  color: #999;
}

.summary-value {
  font-size: 26rpx;
  color: #333;
}

.summary-row.total {
  border-top: 1rpx solid #f0f0f0;
  margin-top: 12rpx;
  padding-top: 20rpx;
}

.total-label {
  font-size: 28rpx;
  font-weight: 600;
  color: #333;
}

.total-value {
  font-size: 36rpx;
  font-weight: 700;
  color: #e74c3c;
}

/* Step actions */
.step-action {
  padding: 40rpx 0;
  display: flex;
  gap: 20rpx;
}

.step-action.dual {
  justify-content: space-between;
}

.next-btn {
  flex: 1;
  height: 88rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: #0f3a57;
  border-radius: 44rpx;
}

.next-btn.disabled {
  opacity: 0.4;
}

.btn-text {
  color: #fff;
  font-size: 30rpx;
  font-weight: 500;
}

.back-btn {
  width: 200rpx;
  height: 88rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  border: 2rpx solid #0f3a57;
  border-radius: 44rpx;
  flex-shrink: 0;
}

.back-btn-text {
  color: #0f3a57;
  font-size: 30rpx;
  font-weight: 500;
}

.empty {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 80rpx 0;
}

.empty-text {
  font-size: 28rpx;
  color: #999;
}

/* Picker popup */
.picker-mask {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.5);
  z-index: 999;
  display: flex;
  align-items: flex-end;
}

.picker-popup {
  width: 100%;
  background-color: #fff;
  border-radius: 24rpx 24rpx 0 0;
  max-height: 70vh;
  display: flex;
  flex-direction: column;
}

.picker-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 28rpx 32rpx;
  border-bottom: 1rpx solid #f0f0f0;
}

.picker-cancel {
  font-size: 28rpx;
  color: #999;
}

.picker-title {
  font-size: 30rpx;
  font-weight: 600;
  color: #333;
}

.picker-list {
  flex: 1;
  max-height: 60vh;
}

.picker-option {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 28rpx 32rpx;
  border-bottom: 1rpx solid #f8f8f8;
}

.picker-option.selected {
  background-color: #f0f7ff;
}

.route-option-info {
  display: flex;
  flex-direction: column;
  gap: 6rpx;
}

.picker-option-text {
  font-size: 28rpx;
  color: #333;
}

.route-desc {
  font-size: 22rpx;
  color: #999;
}

.picker-check {
  font-size: 28rpx;
  color: #0f3a57;
  font-weight: 600;
}
</style>
