import axios from 'axios'
import { message } from 'antd'

const camelizeKey = (key: string) => key.replace(/_([a-z])/g, (_, letter: string) => letter.toUpperCase())

const camelize = (value: unknown): unknown => {
  if (Array.isArray(value)) return value.map(camelize)
  if (value && typeof value === 'object' && !(value instanceof File) && !(value instanceof FormData)) {
    return Object.fromEntries(
      Object.entries(value as Record<string, unknown>).map(([key, val]) => [camelizeKey(key), camelize(val)])
    )
  }
  return value
}

// baseURL: http://localhost:8080/api/v1
const request = axios.create({
  baseURL: '/api/v1',
  timeout: 15000,
})

request.interceptors.request.use((config) => {
  const token = localStorage.getItem('admin_token')
  if (token) config.headers.Authorization = `Bearer ${token}`
  return config
})

request.interceptors.response.use(
  (res): any => {
    const responseData = camelize(res.data) as { code: number; message?: string }
    const { code, message: msg } = responseData
    if (code !== 0) {
      message.error(msg || '请求失败')
      if (code >= 40100 && code < 40200) {
        localStorage.removeItem('admin_token')
        window.location.hash = '#/login'
      }
      return Promise.reject(new Error(msg))
    }
    return responseData
  },
  (err) => {
    message.error(err.message || '网络异常')
    return Promise.reject(err)
  }
)

export default request
