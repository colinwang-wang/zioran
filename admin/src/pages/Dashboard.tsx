import { useEffect, useState } from 'react'
import { Card, Col, Row, Statistic, Table, Select } from 'antd'
import { UserOutlined, BookOutlined, ShoppingCartOutlined, DollarOutlined, ArrowUpOutlined, ArrowDownOutlined } from '@ant-design/icons'
import { getDashboardStats, getOrders } from '@/api'
import type { DashboardStats, Order } from '@/types'
import dayjs from 'dayjs'

export default function Dashboard() {
  const [stats, setStats] = useState<DashboardStats>()
  const [orders, setOrders] = useState<Order[]>([])
  const [period, setPeriod] = useState('day')

  useEffect(() => {
    getDashboardStats().then(res => setStats(res.data))
    getOrders({ page: 1, pageSize: 5 }).then(res => setOrders(res.data.items))
  }, [])

  const statCards = stats ? [
    { title: '总用户数', value: stats.totalUsers, icon: <UserOutlined />, growth: stats.userGrowth },
    { title: '课程总数', value: stats.totalCourses, icon: <BookOutlined />, growth: stats.courseGrowth },
    { title: '订单总数', value: stats.totalOrders, icon: <ShoppingCartOutlined />, growth: stats.orderGrowth },
    { title: '今日收入', value: stats.todayRevenue, icon: <DollarOutlined />, growth: stats.revenueGrowth, prefix: '¥' },
  ] : []

  return (
    <div>
      <Row gutter={16}>
        {statCards.map(s => (
          <Col span={6} key={s.title}>
            <Card>
              <Statistic
                title={s.title}
                value={s.value}
                prefix={s.prefix || s.icon}
                suffix={
                  <span style={{ fontSize: 14, color: s.growth >= 0 ? '#3f8600' : '#cf1322' }}>
                    {s.growth >= 0 ? <ArrowUpOutlined /> : <ArrowDownOutlined />}
                    {Math.abs(s.growth)}%
                  </span>
                }
              />
            </Card>
          </Col>
        ))}
      </Row>

      <Card title="趋势概览" style={{ marginTop: 16 }} extra={
        <Select value={period} onChange={setPeriod} style={{ width: 120 }}
          options={[{ value: 'day', label: '日' }, { value: 'month', label: '月' }, { value: 'quarter', label: '季度' }, { value: 'year', label: '年' }]}
        />
      }>
        <div style={{ height: 200, display: 'flex', alignItems: 'center', justifyContent: 'center', color: '#999' }}>
          图表区域（集成图表库后展示{period === 'day' ? '日' : period === 'month' ? '月' : period === 'quarter' ? '季度' : '年'}趋势数据）
        </div>
      </Card>

      <Card title="最近订单" style={{ marginTop: 16 }}>
        <Table dataSource={orders} rowKey="id" pagination={false} columns={[
          { title: '订单号', dataIndex: 'orderNo' },
          { title: '用户', dataIndex: 'userName' },
          { title: '商品', dataIndex: 'productName' },
          { title: '金额', dataIndex: 'amount', render: (v: number) => `¥${v}` },
          { title: '状态', dataIndex: 'status', render: (v: string) => ({ pending: '待支付', paid: '已支付', refunded: '已退款', cancelled: '已取消' }[v] || v) },
          { title: '时间', dataIndex: 'createdAt', render: (v: string) => dayjs(v).format('YYYY-MM-DD HH:mm') },
        ]} />
      </Card>
    </div>
  )
}
