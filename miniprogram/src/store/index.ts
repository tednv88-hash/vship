/**
 * Simple reactive store using Vue 3 reactivity
 */
import { reactive } from 'vue'

export interface UserInfo {
  id: string
  nickname: string
  avatar: string
  phone: string
  balance: number
  points: number
  is_vip: boolean
  vip_level: number
}

interface AppState {
  user: UserInfo | null
  isLoggedIn: boolean
  cartCount: number
  messageCount: number
}

export const store = reactive<AppState>({
  user: null,
  isLoggedIn: false,
  cartCount: 0,
  messageCount: 0,
})

/** Set user info after login */
export function setUser(user: UserInfo) {
  store.user = user
  store.isLoggedIn = true
}

/** Clear user on logout */
export function clearUser() {
  store.user = null
  store.isLoggedIn = false
  store.cartCount = 0
  store.messageCount = 0
}

/** Update cart badge count */
export function setCartCount(count: number) {
  store.cartCount = count
  if (count > 0) {
    uni.setTabBarBadge({ index: 2, text: String(count) })
  } else {
    uni.removeTabBarBadge({ index: 2 })
  }
}

/** Update message badge count */
export function setMessageCount(count: number) {
  store.messageCount = count
}

export default store
