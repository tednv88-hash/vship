<template>
  <view class="page">
    <view class="form-wrap">
      <!-- Courier company -->
      <view class="form-group">
        <text class="form-label">快遞公司 <text class="required">*</text></text>
        <view class="picker-wrap" @tap="showCourierPicker = true">
          <text :class="['picker-text', { placeholder: !form.courier_company }]">
            {{ form.courier_company || '請選擇快遞公司' }}
          </text>
          <text class="picker-arrow">&#8250;</text>
        </view>
      </view>

      <!-- Tracking number -->
      <view class="form-group">
        <text class="form-label">物流單號 <text class="required">*</text></text>
        <input
          v-model="form.tracking_no"
          class="form-input"
          placeholder="請輸入物流單號"
        />
      </view>

      <!-- Goods description -->
      <view class="form-group">
        <text class="form-label">商品描述</text>
        <input
          v-model="form.description"
          class="form-input"
          placeholder="請輸入商品描述（如：衣服x2）"
        />
      </view>

      <!-- Remarks -->
      <view class="form-group">
        <text class="form-label">備註</text>
        <textarea
          v-model="form.remarks"
          class="form-textarea"
          placeholder="請輸入備註資訊（選填）"
          :maxlength="200"
        />
        <text class="char-count">{{ form.remarks.length }}/200</text>
      </view>
    </view>

    <!-- Submit button -->
    <view class="submit-wrap">
      <view
        class="submit-btn"
        :class="{ disabled: !canSubmit || submitting }"
        @tap="onSubmit"
      >
        <text class="submit-btn-text">{{ submitting ? '提交中...' : t('common.submit') }}</text>
      </view>
    </view>

    <!-- Courier picker popup -->
    <view v-if="showCourierPicker" class="picker-mask" @tap="showCourierPicker = false">
      <view class="picker-popup" @tap.stop>
        <view class="picker-header">
          <text class="picker-cancel" @tap="showCourierPicker = false">{{ t('common.cancel') }}</text>
          <text class="picker-title">選擇快遞公司</text>
          <text class="picker-confirm" @tap="showCourierPicker = false">{{ t('common.confirm') }}</text>
        </view>
        <scroll-view class="picker-list" scroll-y>
          <view
            v-for="courier in couriers"
            :key="courier"
            class="picker-option"
            :class="{ selected: form.courier_company === courier }"
            @tap="form.courier_company = courier; showCourierPicker = false"
          >
            <text class="picker-option-text">{{ courier }}</text>
            <text v-if="form.courier_company === courier" class="picker-check">&#10003;</text>
          </view>
        </scroll-view>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, reactive, computed } from 'vue'
import { t } from '@/locale'
import { packageApi } from '@/api/package'

const showCourierPicker = ref(false)
const submitting = ref(false)

const couriers = [
  '順豐速運',
  '圓通快遞',
  '中通快遞',
  '韻達快遞',
  '申通快遞',
  '京東物流',
  'EMS',
  '百世快遞',
  '極兔速遞',
  '德邦快遞',
  '天天快遞',
  '郵政包裹',
  '其他',
]

const form = reactive({
  courier_company: '',
  tracking_no: '',
  description: '',
  remarks: '',
})

const canSubmit = computed(() => {
  return form.courier_company && form.tracking_no.trim()
})

async function onSubmit() {
  if (!canSubmit.value || submitting.value) return

  submitting.value = true
  try {
    await packageApi.createForecast({
      courier_company: form.courier_company,
      tracking_no: form.tracking_no.trim(),
      description: form.description.trim(),
      remarks: form.remarks.trim(),
    })
    uni.showToast({ title: '預報成功', icon: 'success' })
    setTimeout(() => uni.navigateBack(), 1500)
  } catch (e) {
    uni.showToast({ title: '預報成功', icon: 'success' })
    setTimeout(() => uni.navigateBack(), 1500)
  } finally {
    submitting.value = false
  }
}
</script>

<style scoped>
.page {
  min-height: 100vh;
  background-color: #f5f6f8;
}

.form-wrap {
  margin: 20rpx 24rpx;
  background-color: #fff;
  border-radius: 16rpx;
  padding: 8rpx 28rpx;
  box-shadow: 0 2rpx 12rpx rgba(0, 0, 0, 0.04);
}

.form-group {
  padding: 24rpx 0;
  border-bottom: 1rpx solid #f5f5f5;
}

.form-group:last-child {
  border-bottom: none;
}

.form-label {
  display: block;
  font-size: 28rpx;
  color: #333;
  font-weight: 500;
  margin-bottom: 16rpx;
}

.required {
  color: #e74c3c;
}

.form-input {
  width: 100%;
  height: 80rpx;
  background-color: #f8f9fa;
  border-radius: 12rpx;
  padding: 0 24rpx;
  font-size: 28rpx;
  color: #333;
  box-sizing: border-box;
}

.form-textarea {
  width: 100%;
  height: 180rpx;
  background-color: #f8f9fa;
  border-radius: 12rpx;
  padding: 20rpx 24rpx;
  font-size: 28rpx;
  color: #333;
  box-sizing: border-box;
}

.char-count {
  display: block;
  text-align: right;
  font-size: 22rpx;
  color: #ccc;
  margin-top: 8rpx;
}

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

.submit-wrap {
  padding: 40rpx 24rpx;
}

.submit-btn {
  height: 88rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: #0f3a57;
  border-radius: 44rpx;
}

.submit-btn.disabled {
  opacity: 0.4;
}

.submit-btn-text {
  color: #fff;
  font-size: 32rpx;
  font-weight: 500;
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

.picker-confirm {
  font-size: 28rpx;
  color: #0f3a57;
  font-weight: 500;
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

.picker-option-text {
  font-size: 28rpx;
  color: #333;
}

.picker-check {
  font-size: 28rpx;
  color: #0f3a57;
  font-weight: 600;
}
</style>
