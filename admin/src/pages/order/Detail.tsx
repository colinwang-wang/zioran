import { useEffect, useState } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { Card, Descriptions, Tag, Button, Modal, message, Spin } from 'antd'
import { ExclamationCircleOutlined } from '@ant-design/icons'
import { getOrder, refundOrder } from '@/api'
import type { Order } from '@/types'
import dayjs from 'dayjs'

const statusColors: Record<string, string> = { pending: 'orange', paid: 'green', refunded: 'blue', cancelled: 'red' }
const statusLabels: Record<string, string> = { pending: '待支付', paid: '已支付', refunded: '已退款', cancelled: '已取消' }

export default function OrderDetail() {
  const { id } = useParams()
  const navigate = useNavigate()
  const [order, setOrder] = useState<Order>()
  const [loading, setLoading] = useState(true)

  const fetchData = () => {
    if (!id) return
    setLoading(true)
    getOrder(Number(id)).then(res => setOrder(res.data)).finally(() => setLoading(false))
  }

  useEffect(() => { fetchData() }, [id])

  const handleRefund = () => {
    Modal.confirm({
      title: '确认退款', icon: <ExclamationCircleOutlined />, content: '退款后不可撤销',
      onOk: async () => { await refundOrder(Number(id)); message.success('退款成功'); fetchData() }
    })
  }

  if (loading) return <Spin style={{ display: 'block', marginTop: 100 }} />
  if (!order) return <Card><p>订单不存在</p></Card>

  return (
    <Card title="订单详情" extra={<Button onClick={() => navigate('/orders')}>返回列表</Button>}>
      <Descriptions bordered column={2}>
        <Descriptions.Item label="订单号">{order.orderNo}</Descriptions.Item>
        <Descriptions.Item label="用户">{order.userName}</Descriptions.Item>
        <Descriptions.Item label="商品">{order.productName}</Descriptions.Item>
        <Descriptions.Item label="金额">¥{order.amount}</Descriptions.Item>
        <Descriptions.Item label="支付方式">{order.payMethod || '-'}</Descriptions.Item>
        <Descriptions.Item label="状态"><Tag color={statusColors[order.status]}>{statusLabels[order.status]}</Tag></Descriptions.Item>
        <Descriptions.Item label="下单时间">{dayjs(order.createdAt).format('YYYY-MM-DD HH:mm:ss')}</Descriptions.Item>
        <Descriptions.Item label="支付时间">{order.paidAt ? dayjs(order.paidAt).format('YYYY-MM-DD HH:mm:ss') : '-'}</Descriptions.Item>
      </Descriptions>
      {order.status === 'paid' && (
        <div style={{ marginTop: 24 }}>
          <Button type="primary" danger onClick={handleRefund}>退款</Button>
        </div>
      )}
    </Card>
  )
}
