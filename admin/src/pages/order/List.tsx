import { useEffect, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { Table, Button, Space, Select, DatePicker, Tag, Card, Input } from 'antd'
import { getOrders } from '@/api'
import type { Order } from '@/types'
import dayjs from 'dayjs'

const { RangePicker } = DatePicker

const statusColors: Record<string, string> = { pending: 'orange', paid: 'green', refunded: 'blue', cancelled: 'red' }
const statusLabels: Record<string, string> = { pending: '待支付', paid: '已支付', refunded: '已退款', cancelled: '已取消' }
const typeLabels: Record<string, string> = { coin_recharge: '金币充值', vip_purchase: 'VIP购买', course_purchase: '课程购买' }

export default function OrderList() {
  const navigate = useNavigate()
  const [data, setData] = useState<Order[]>([])
  const [total, setTotal] = useState(0)
  const [loading, setLoading] = useState(false)
  const [params, setParams] = useState<Record<string, unknown>>({ page: 1, pageSize: 20, keyword: '', type: undefined, status: undefined, startDate: undefined, endDate: undefined })

  const fetchData = async () => {
    setLoading(true)
    try { const res = await getOrders(params); setData(res.data.items); setTotal(res.data.total) }
    finally { setLoading(false) }
  }

  useEffect(() => { fetchData() }, [params])

  const handleSearch = (key: string, val: unknown) => setParams(p => ({ ...p, [key]: val, page: 1 }))

  return (
    <Card>
      <Space style={{ marginBottom: 16 }} wrap>
        <Input.Search placeholder="搜索订单号/商品名/用户名" allowClear onSearch={v => handleSearch('keyword', v)} style={{ width: 240 }} />
        <Select placeholder="订单类型" allowClear style={{ width: 130 }} onChange={v => handleSearch('type', v)}
          options={[{ label: '金币充值', value: 'coin_recharge' }, { label: 'VIP购买', value: 'vip_purchase' }, { label: '课程购买', value: 'course_purchase' }]} />
        <Select placeholder="支付状态" allowClear style={{ width: 120 }} onChange={v => handleSearch('status', v)}
          options={[{ label: '待支付', value: 'pending' }, { label: '已支付', value: 'paid' }, { label: '已退款', value: 'refunded' }, { label: '已取消', value: 'cancelled' }]} />
        <RangePicker onChange={dates => {
          if (dates) setParams(p => ({ ...p, startDate: dates[0]?.format('YYYY-MM-DD'), endDate: dates[1]?.format('YYYY-MM-DD'), page: 1 }))
          else setParams(p => ({ ...p, startDate: undefined, endDate: undefined, page: 1 }))
        }} />
      </Space>

      <Table dataSource={data} rowKey="id" loading={loading}
        pagination={{ current: params.page as number, pageSize: params.pageSize as number, total, onChange: (p, ps) => setParams(prev => ({ ...prev, page: p, pageSize: ps })) }}
        columns={[
          { title: '订单号', dataIndex: 'orderNo', width: 180 },
          { title: '用户', dataIndex: 'userName', width: 100 },
          { title: '商品', dataIndex: 'productName', ellipsis: true },
          { title: '类型', dataIndex: 'type', width: 100, render: (v: string) => typeLabels[v] || v },
          { title: '金额', dataIndex: 'amount', width: 80, render: (v: number) => `¥${v}` },
          { title: '状态', dataIndex: 'status', width: 90, render: (v: string) => <Tag color={statusColors[v]}>{statusLabels[v]}</Tag> },
          { title: '下单时间', dataIndex: 'createdAt', width: 160, render: (v: string) => dayjs(v).format('YYYY-MM-DD HH:mm') },
          { title: '操作', width: 80, render: (_: unknown, r: Order) => (
            <Button type="link" size="small" onClick={() => navigate(`/orders/${r.id}`)}>详情</Button>
          )},
        ]}
      />
    </Card>
  )
}
