import { useEffect, useState } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { Card, Descriptions, Tag, Button, Spin } from 'antd'
import { getUser } from '@/api'
import type { User } from '@/types'
import dayjs from 'dayjs'

export default function UserDetail() {
  const { id } = useParams()
  const navigate = useNavigate()
  const [user, setUser] = useState<User>()
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    if (id) {
      setLoading(true)
      getUser(Number(id)).then(res => setUser(res.data)).finally(() => setLoading(false))
    }
  }, [id])

  if (loading) return <Spin style={{ display: 'block', marginTop: 100 }} />
  if (!user) return <Card><p>用户不存在</p></Card>

  return (
    <Card title="用户详情" extra={<Button onClick={() => navigate('/users')}>返回列表</Button>}>
      <Descriptions bordered column={2}>
        <Descriptions.Item label="ID">{user.id}</Descriptions.Item>
        <Descriptions.Item label="手机号">{user.phone?.replace(/(\d{3})\d{4}(\d{4})/, '$1****$2')}</Descriptions.Item>
        <Descriptions.Item label="昵称">{user.nickname}</Descriptions.Item>
        <Descriptions.Item label="状态"><Tag color={user.status === 'active' ? 'green' : 'red'}>{user.status === 'active' ? '正常' : '已禁用'}</Tag></Descriptions.Item>
        <Descriptions.Item label="金币余额">{user.balance} 金币</Descriptions.Item>
        <Descriptions.Item label="VIP到期时间">{user.vipExpireAt ? dayjs(user.vipExpireAt).format('YYYY-MM-DD') : '非VIP'}</Descriptions.Item>
        <Descriptions.Item label="已购课程">{user.purchasedCount} 个</Descriptions.Item>
        <Descriptions.Item label="收藏课程">{user.favoriteCount} 个</Descriptions.Item>
        <Descriptions.Item label="注册时间" span={2}>{dayjs(user.createdAt).format('YYYY-MM-DD HH:mm:ss')}</Descriptions.Item>
      </Descriptions>
    </Card>
  )
}
