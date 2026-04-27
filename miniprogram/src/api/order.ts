/** Order API */
import http from '@/utils/request'

export const orderApi = {
  /** Get order list */
  getList: (params?: any) => http.get('/orders', params),

  /** Get order detail */
  getDetail: (id: string) => http.get(`/orders/${id}`),

  /** Create consolidation order (集运下单) */
  createOrder: (data: any) => http.post('/orders', data),

  /** Cancel order */
  cancelOrder: (id: string) => http.put(`/orders/${id}/cancel`),

  /** Pay order */
  payOrder: (id: string, data: { payment_method: string }) =>
    http.post(`/orders/${id}/pay`, data),

  /** Get shop order list */
  getShopOrders: (params?: any) => http.get('/shop-orders', params),

  /** Get shop order detail */
  getShopOrderDetail: (id: string) => http.get(`/shop-orders/${id}`),

  /** Checkout cart -> create shop order */
  checkoutCart: (data: { cart_ids: string[]; address_id: string; remark?: string; pay_method?: string }) =>
    http.post('/shop-orders/checkout', data),
}

export default orderApi
