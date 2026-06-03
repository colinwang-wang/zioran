import axios from 'axios'
import { message } from 'antd'

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
  (res) => {
    const { code, message: msg } = res.data
    if (code !== 0) {
      message.error(msg || '请求失败')
      if (code >= 40100 && code < 40200) {
        localStorage.removeItem('admin_token')
        window.location.hash = '#/login'
      }
      return Promise.reject(new Error(msg))
    }
    return res.data
  },
  (err) => {
    message.error(err.message || '网络异常')
    return Promise.reject(err)
  }
)

export default request
