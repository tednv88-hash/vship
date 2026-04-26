/** Package API */
import http from '@/utils/request'

export const packageApi = {
  /** Get package list */
  getList: (params?: any) => http.get('/packages', params),

  /** Get package detail */
  getDetail: (id: string) => http.get(`/packages/${id}`),

  /** Create forecast (预报) */
  createForecast: (data: any) => http.post('/packages/forecast', data),

  /** Get forecast list */
  getForecastList: (params?: any) => http.get('/packages/forecast', params),

  /** Merge packages (合箱) */
  mergePackages: (data: { package_ids: string[] }) =>
    http.post('/packages/merge', data),

  /** Split package (拆箱) */
  splitPackage: (id: string, data: any) =>
    http.post(`/packages/${id}/split`, data),
}

export default packageApi
