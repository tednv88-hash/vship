/** User / Auth API */
import http from '@/utils/request'

export const userApi = {
  /** Login with phone + password */
  login: (data: { phone: string; password: string }) =>
    http.post('/auth/login', data),

  /** Register */
  register: (data: { phone: string; password: string; code?: string; email?: string; name?: string }) =>
    http.post('/auth/register', {
      phone: data.phone,
      password: data.password,
      email: data.email || `${data.phone}@guoyun.local`,
      name: data.name || data.phone,
      code: data.code,
    }),

  /** WeChat login */
  wechatLogin: (data: { code: string }) =>
    http.post('/auth/wechat-login', data),

  /** Bind phone */
  bindPhone: (data: { phone: string; code: string }) =>
    http.post('/auth/bind-phone', data),

  /** Get SMS code */
  sendCode: (data: { phone: string; type: string }) =>
    http.post('/auth/send-code', data),

  /** Get current user info */
  getUserInfo: () => http.get('/user/info'),

  /** Update user info */
  updateUserInfo: (data: any) => http.put('/user/info', data),

  /** Get balance details */
  getBalanceList: (params?: any) => http.get('/user/balance', params),

  /** Get points details */
  getPointsList: (params?: any) => http.get('/user/points', params),

  /** Recharge */
  recharge: (data: { amount: number; payment_method: string }) =>
    http.post('/user/recharge', data),

  /** Identity verification */
  verifyIdentity: (data: any) => http.post('/user/identity', data),

  /** Get identity status */
  getIdentityStatus: () => http.get('/user/identity'),

  /** Sign in (daily check-in) */
  signIn: () => http.post('/user/signin'),

  /** Get sign-in status */
  getSignInStatus: () => http.get('/user/signin'),
}

export default userApi
