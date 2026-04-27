<template>
  <view class="address-edit-page">
    <view class="form-section">
      <!-- Recipient -->
      <view class="form-item">
        <text class="label">{{ t('address.name') }}</text>
        <input
          v-model="form.name"
          :placeholder="`請輸入${t('address.name')}`"
          class="input"
        />
      </view>

      <!-- Phone -->
      <view class="form-item">
        <text class="label">{{ t('address.phone') }}</text>
        <input
          v-model="form.phone"
          type="number"
          :placeholder="`請輸入${t('address.phone')}`"
          maxlength="11"
          class="input"
        />
      </view>

      <!-- Region picker -->
      <view class="form-item">
        <text class="label">{{ t('address.region') }}</text>
        <picker
          mode="region"
          :value="regionValue"
          @change="onRegionChange"
        >
          <view class="picker-display">
            <text :class="['picker-text', { placeholder: !regionText }]">
              {{ regionText || `請選擇${t('address.region')}` }}
            </text>
            <text class="arrow">&#x276F;</text>
          </view>
        </picker>
      </view>

      <!-- Detail address -->
      <view class="form-item">
        <text class="label">{{ t('address.detail') }}</text>
        <textarea
          v-model="form.address"
          :placeholder="`請輸入${t('address.detail')}`"
          class="textarea"
          :maxlength="200"
        />
      </view>

      <!-- Default switch -->
      <view class="form-item switch-item">
        <text class="label">{{ t('address.default') }}</text>
        <switch
          :checked="form.is_default"
          @change="form.is_default = ($event as any).detail.value"
          color="#0f3a57"
        />
      </view>
    </view>

    <!-- Save button -->
    <view class="save-btn" @tap="handleSave">
      <text class="save-btn-text">{{ t('common.save') }}</text>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { t } from '@/locale'
import { commonApi } from '@/api/common'

const isEdit = ref(false)
const addressId = ref('')
const loading = ref(false)

const form = reactive({
  name: '',
  phone: '',
  province: '',
  city: '',
  district: '',
  address: '',
  is_default: false,
})

const showRegionPicker = ref(false)

const regionText = computed(() => {
  if (form.province && form.city && form.district) {
    return `${form.province} ${form.city} ${form.district}`
  }
  if (form.province || form.city) {
    return [form.province, form.city, form.district].filter(Boolean).join(' ')
  }
  return ''
})

const regionValue = computed(() => {
  if (form.province && form.city) {
    return [form.province, form.city, form.district || '']
  }
  return ['廣東省', '深圳市', '南山區']
})

function onRegionChange(e: any) {
  const values = e.detail.value || []
  form.province = values[0] || ''
  form.city = values[1] || ''
  form.district = values[2] || ''
}

async function loadAddress(id: string) {
  try {
    const res = (await commonApi.getAddress(id)) as any
    const data = res?.data || res
    if (data) {
      form.name = data.recipient_name || data.name || ''
      form.phone = data.phone || ''
      form.province = data.province || ''
      form.city = data.city || ''
      form.district = data.district || ''
      form.address = data.address || ''
      form.is_default = !!data.is_default
    }
  } catch {
    uni.showToast({ title: '載入失敗', icon: 'none' })
  }
}

async function handleSave() {
  if (!form.name) {
    uni.showToast({ title: `請輸入${t('address.name')}`, icon: 'none' })
    return
  }
  if (!form.phone) {
    uni.showToast({ title: `請輸入${t('address.phone')}`, icon: 'none' })
    return
  }
  if (!form.province || !form.city || !form.district) {
    uni.showToast({ title: `請選擇${t('address.region')}`, icon: 'none' })
    return
  }
  if (!form.address) {
    uni.showToast({ title: `請輸入${t('address.detail')}`, icon: 'none' })
    return
  }

  if (loading.value) return
  loading.value = true

  try {
    const payload: any = {
      recipient_name: form.name,
      phone: form.phone,
      province: form.province,
      city: form.city,
      district: form.district,
      address: form.address,
      is_default: form.is_default,
    }
    if (isEdit.value) {
      await commonApi.updateAddress(addressId.value, payload)
    } else {
      await commonApi.createAddress(payload)
    }
    uni.showToast({ title: '保存成功', icon: 'success' })
    uni.$emit('addressUpdated')
    setTimeout(() => {
      uni.navigateBack()
    }, 1000)
  } catch {
    uni.showToast({ title: '保存失敗', icon: 'none' })
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  const pages = getCurrentPages()
  const currentPage = pages[pages.length - 1] as any
  const query = currentPage?.$page?.options || currentPage?.options || {}
  if (query.id) {
    isEdit.value = true
    addressId.value = query.id
    uni.setNavigationBarTitle({ title: t('address.edit') })
    loadAddress(query.id)
  } else {
    uni.setNavigationBarTitle({ title: t('address.add') })
  }
})
</script>

<style scoped>
.address-edit-page {
  min-height: 100vh;
  background: #f7f8fa;
  padding: 24rpx;
}

.form-section {
  background: #fff;
  border-radius: 16rpx;
  padding: 10rpx 30rpx;
}

.form-item {
  padding: 24rpx 0;
  border-bottom: 1rpx solid #f0f0f0;
}

.form-item:last-child {
  border-bottom: none;
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
  height: 80rpx;
  background: #f7f8fa;
  border-radius: 12rpx;
  padding: 0 24rpx;
  font-size: 28rpx;
  color: #333;
  box-sizing: border-box;
}

.picker-display {
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 80rpx;
  background: #f7f8fa;
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

.arrow {
  font-size: 28rpx;
  color: #ccc;
}

.textarea {
  width: 100%;
  height: 160rpx;
  background: #f7f8fa;
  border-radius: 12rpx;
  padding: 20rpx 24rpx;
  font-size: 28rpx;
  color: #333;
  box-sizing: border-box;
}

.switch-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.switch-item .label {
  margin-bottom: 0;
}

.save-btn {
  width: 100%;
  height: 96rpx;
  background: #0f3a57;
  border-radius: 16rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-top: 40rpx;
}

.save-btn-text {
  font-size: 32rpx;
  color: #fff;
  font-weight: 600;
}
</style>
