/** Common / Misc API */
import http from '@/utils/request'

export const commonApi = {
  /** Get page design data (for homepage etc.) */
  getPageDesign: (params?: { is_default?: boolean; type?: string }) =>
    http.get('/page-designs', params),

  /** Get tracking info */
  getTracking: (trackingNo: string) =>
    http.get(`/tracking/${trackingNo}`),

  /** Address CRUD */
  getAddresses: () => http.get('/addresses'),
  getAddress: (id: string) => http.get(`/addresses/${id}`),
  createAddress: (data: any) => http.post('/addresses', data),
  updateAddress: (id: string, data: any) => http.put(`/addresses/${id}`, data),
  deleteAddress: (id: string) => http.delete(`/addresses/${id}`),
  setDefaultAddress: (id: string) => http.put(`/addresses/${id}/default`),

  /** Cart */
  getCart: () => http.get('/cart'),
  addToCart: (data: { goods_id: string; sku_id?: string; quantity: number }) =>
    http.post('/cart', data),
  updateCartItem: (id: string, data: { quantity: number }) =>
    http.put(`/cart/${id}`, data),
  deleteCartItem: (id: string) => http.delete(`/cart/${id}`),

  /** Coupons */
  getCoupons: (params?: any) => http.get('/coupons', params),
  getMyCoupons: (params?: any) => http.get('/user/coupons', params),
  claimCoupon: (id: string) => http.post(`/coupons/${id}/claim`),

  /** Help articles */
  getHelpList: (params?: any) => http.get('/help-articles', params),
  getHelpDetail: (id: string) => http.get(`/help-articles/${id}`),

  /** Articles */
  getArticles: (params?: any) => http.get('/articles', params),
  getArticleDetail: (id: string) => http.get(`/articles/${id}`),

  /** Notices */
  getNotices: (params?: any) => http.get('/notices', params),
  getNoticeDetail: (id: string) => http.get(`/notices/${id}`),

  /** Warehouse list */
  getWarehouses: () => http.get('/warehouses'),
  getWarehouseDetail: (id: string) => http.get(`/warehouses/${id}`),

  /** Shipping routes */
  getRoutes: (params?: any) => http.get('/routes', params),
  getRouteDetail: (id: string) => http.get(`/routes/${id}`),

  /** Estimate / calculator */
  calculateEstimate: (data: any) => http.post('/estimate', data),

  /** Messages */
  getMessages: (params?: any) => http.get('/messages', params),
  markRead: (id: string) => http.put(`/messages/${id}/read`),

  /** Favorites */
  getFavorites: (params?: any) => http.get('/favorites', params),
  addFavorite: (data: { goods_id: string }) => http.post('/favorites', data),
  removeFavorite: (id: string) => http.delete(`/favorites/${id}`),

  /** Browsing history */
  getHistory: (params?: any) => http.get('/history', params),
  clearHistory: () => http.delete('/history'),

  /** Feedback */
  submitFeedback: (data: any) => http.post('/feedback', data),

  /** Reviews */
  getReviews: (params?: any) => http.get('/reviews', params),
  createReview: (data: any) => http.post('/reviews', data),

  /** Refund */
  getRefunds: (params?: any) => http.get('/refunds', params),
  getRefundDetail: (id: string) => http.get(`/refunds/${id}`),
  createRefund: (data: any) => http.post('/refunds', data),

  /** Dealer / Affiliate */
  getDealerInfo: () => http.get('/dealer'),
  applyDealer: (data: any) => http.post('/dealer/apply', data),
  getDealerOrders: (params?: any) => http.get('/dealer/orders', params),
  getDealerWithdrawals: (params?: any) => http.get('/dealer/withdrawals', params),
  requestWithdraw: (data: any) => http.post('/dealer/withdraw', data),
  getDealerTeam: (params?: any) => http.get('/dealer/team', params),

  /** Prohibited items */
  getProhibitedItems: () => http.get('/prohibited-items'),

  /** About / Policy content */
  getAbout: () => http.get('/content/about'),
  getPrivacy: () => http.get('/content/privacy'),
  getTerms: () => http.get('/content/terms'),

  /** Invite */
  getInviteInfo: () => http.get('/invite'),
  getInvitePoster: () => http.get('/invite/poster'),

  /** Value-added services */
  getValueAddedServices: () => http.get('/services/value-added'),

  /** Insurance */
  getInsuranceOptions: () => http.get('/services/insurance'),

  /** App settings */
  getAppSettings: (settingType: string) =>
    http.get('/app-settings', { setting_type: settingType }),
}

export default commonApi
