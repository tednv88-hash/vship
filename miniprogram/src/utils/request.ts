/**
 * HTTP request wrapper for uni-app
 * Wraps uni.request with token injection, error handling, and base URL config
 */

const BASE_URL = import.meta.env.VITE_API_BASE || '/api/v1/mp'

interface RequestOptions {
  url: string
  method?: 'GET' | 'POST' | 'PUT' | 'DELETE'
  data?: any
  header?: Record<string, string>
  showLoading?: boolean
  loadingText?: string
}

interface ApiResponse<T = any> {
  code: number
  message: string
  data: T
}

/** Get stored auth token */
export function getToken(): string {
  return uni.getStorageSync('token') || ''
}

/** Set auth token */
export function setToken(token: string) {
  uni.setStorageSync('token', token)
}

/** Remove auth token */
export function removeToken() {
  uni.removeStorageSync('token')
}

/** Check if user is logged in */
export function isLoggedIn(): boolean {
  return !!getToken()
}

/** Core request function */
export function request<T = any>(options: RequestOptions): Promise<ApiResponse<T>> {
  const { url, method = 'GET', data, header = {}, showLoading = false, loadingText = '載入中...' } = options

  if (showLoading) {
    uni.showLoading({ title: loadingText, mask: true })
  }

  const token = getToken()
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
    ...header,
  }
  if (token) {
    headers['Authorization'] = `Bearer ${token}`
  }

  return new Promise((resolve, reject) => {
    uni.request({
      url: url.startsWith('http') ? url : `${BASE_URL}${url}`,
      method,
      data,
      header: headers,
      success: (res) => {
        if (showLoading) uni.hideLoading()

        const statusCode = res.statusCode
        if (statusCode === 401) {
          // Token expired — redirect to login
          removeToken()
          uni.reLaunch({ url: '/pages/login/index' })
          reject(new Error('Unauthorized'))
          return
        }

        if (statusCode >= 200 && statusCode < 300) {
          resolve(res.data as ApiResponse<T>)
        } else {
          const errMsg = (res.data as any)?.message || `Request failed (${statusCode})`
          uni.showToast({ title: errMsg, icon: 'none', duration: 2000 })
          reject(new Error(errMsg))
        }
      },
      fail: (err) => {
        if (showLoading) uni.hideLoading()
        uni.showToast({ title: '網絡錯誤，請稍後重試', icon: 'none', duration: 2000 })
        reject(err)
      },
    })
  })
}

/** Convenience methods */
export const http = {
  get: <T = any>(url: string, data?: any, options?: Partial<RequestOptions>) =>
    request<T>({ url, method: 'GET', data, ...options }),

  post: <T = any>(url: string, data?: any, options?: Partial<RequestOptions>) =>
    request<T>({ url, method: 'POST', data, ...options }),

  put: <T = any>(url: string, data?: any, options?: Partial<RequestOptions>) =>
    request<T>({ url, method: 'PUT', data, ...options }),

  delete: <T = any>(url: string, data?: any, options?: Partial<RequestOptions>) =>
    request<T>({ url, method: 'DELETE', data, ...options }),
}

export default http
