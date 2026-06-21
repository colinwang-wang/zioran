import { useEffect, useState } from 'react'
import { Card, Col, Row, Statistic, Table, Select, Empty } from 'antd'
import { UserOutlined, BookOutlined, ShoppingCartOutlined, DollarOutlined, ArrowUpOutlined, ArrowDownOutlined } from '@ant-design/icons'
import { getDashboardStats, getDashboardCharts, getOrders } from '@/api'
import type { DashboardStats, Order, ChartData } from '@/types'
import dayjs from 'dayjs'

export default function Dashboard() {
  const [stats, setStats] = useState<DashboardStats>()
  const [chartData, setChartData] = useState<ChartData>()
  const [orders, setOrders] = useState<Order[]>([])
  const [period, setPeriod] = useState('day')

  useEffect(() => {
    getDashboardStats().then(res => setStats(res.data))
    getOrders({ page: 1, pageSize: 5 }).then(res => setOrders(res.data.items))
  }, [])
  useEffect(() => {
    setChartData(undefined)
    getDashboardCharts(period).then(res => setChartData(res.data)).catch(() => setChartData({ labels: [], datasets: [] }))
  }, [period])

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
        {chartData?.datasets?.length ? (
          <div style={{ display: 'flex', gap: 8, alignItems: 'flex-end', height: 180, padding: '0 16px' }}>
            {chartData.labels.map((label, i) => {
              const total = chartData.datasets.reduce((sum, ds) => sum + (ds.data[i] || 0), 0)
              const max = Math.max(1, ...chartData.labels.map((_, index) => chartData.datasets.reduce((sum, ds) => sum + (ds.data[index] || 0), 0)))
              return (
                <div key={label} style={{ flex: 1, display: 'flex', flexDirection: 'column', alignItems: 'center' }}>
                  <div style={{ width: '60%', height: `${Math.max(4, (total / max) * 100)}%`, background: 'linear-gradient(180deg, #1677ff 0%, #69b1ff 100%)', borderRadius: 4, transition: 'height 0.3s' }} />
                  <span style={{ fontSize: 12, color: '#666', marginTop: 6 }}>{label}</span>
                </div>
              )
            })}
          </div>
        ) : <Empty description="暂无趋势数据" />}
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
