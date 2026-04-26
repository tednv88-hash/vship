/** Goods / Products API */
import http from '@/utils/request'

export const goodsApi = {
  /** Get goods list */
  getList: (params?: any) => http.get('/goods', params),

  /** Get goods detail */
  getDetail: (id: string) => http.get(`/goods/${id}`),

  /** Search goods */
  search: (params: { keyword: string; page?: number; limit?: number }) =>
    http.get('/goods', params),

  /** Get categories */
  getCategories: () => http.get('/categories'),

  /** Get goods by category */
  getByCategory: (categoryId: string, params?: any) =>
    http.get(`/categories/${categoryId}/goods`, params),
}

export default goodsApi
