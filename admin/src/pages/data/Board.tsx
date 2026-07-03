import { useEffect, useState } from 'react'
import { Card, Col, Row, Statistic, Select, Empty } from 'antd'
import { UserOutlined, BookOutlined, ShoppingCartOutlined, DollarOutlined, ArrowUpOutlined, ArrowDownOutlined, HeartOutlined, ClockCircleOutlined, UserAddOutlined } from '@ant-design/icons'
import { getDashboardStats, getDashboardCharts } from '@/api'
import type { DashboardStats, ChartData } from '@/types'

export default function DataBoard() {
  const [stats, setStats] = useState<DashboardStats>()
  const [chartData, setChartData] = useState<ChartData>()
  const [period, setPeriod] = useState('month')

  useEffect(() => { getDashboardStats().then(res => setStats(res.data)) }, [])
  useEffect(() => {
    setChartData(undefined)
    getDashboardCharts(period).then(res => setChartData(res.data)).catch(() => setChartData({ labels: [], datasets: [] }))
  }, [period])

  const statCards = stats ? [
    { title: '总用户数', value: stats.totalUsers, icon: <UserOutlined />, growth: stats.userGrowth },
    { title: '课程总数', value: stats.totalCourses, icon: <BookOutlined />, growth: stats.courseGrowth },
    { title: '订单总数', value: stats.totalOrders, icon: <ShoppingCartOutlined />, growth: stats.orderGrowth },
    { title: '今日收入', value: stats.todayRevenue, icon: <DollarOutlined />, growth: stats.revenueGrowth, prefix: '¥' },
    { title: '课程收藏', value: stats.totalFavorites || 0, icon: <HeartOutlined />, growth: 0 },
    { title: '待支付订单', value: stats.pendingOrders || 0, icon: <ClockCircleOutlined />, growth: 0 },
    { title: '新增用户(近3天)', value: stats.recentNewUsers || 0, icon: <UserAddOutlined />, growth: 0 },
  ] : []

  return (
    <div>
      <Row gutter={[16, 16]}>
        {statCards.map(s => (
          <Col xs={12} sm={8} md={6} key={s.title}>
            <Card>
              <Statistic title={s.title} value={s.value} prefix={s.prefix || s.icon}
                suffix={s.growth !== 0 ? <span style={{ fontSize: 14, color: s.growth >= 0 ? '#3f8600' : '#cf1322' }}>{s.growth >= 0 ? <ArrowUpOutlined /> : <ArrowDownOutlined />}{Math.abs(s.growth).toFixed(0)}%</span> : undefined} />
            </Card>
          </Col>
        ))}
      </Row>

      <Card title="数据趋势" style={{ marginTop: 16 }} extra={
        <Select value={period} onChange={setPeriod} style={{ width: 120 }}
          options={[{ value: 'day', label: '日' }, { value: 'month', label: '月' }, { value: 'quarter', label: '季度' }, { value: 'year', label: '年' }]} />
      }>
        <div style={{ minHeight: 300 }}>
          {chartData?.datasets?.length ? (
            <div>
              {chartData.datasets?.map((ds, i) => (
                <div key={i} style={{ marginBottom: 16 }}>
                  <strong>{ds.label}</strong>
                  <div style={{ display: 'flex', gap: 4, marginTop: 8, alignItems: 'flex-end', height: 100 }}>
                    {ds.data?.map((v, j) => (
                      <div key={j} style={{ display: 'flex', flexDirection: 'column', alignItems: 'center', flex: 1 }}>
                        <div style={{ background: '#1677ff', width: '80%', height: Math.max(4, (v / Math.max(...ds.data)) * 80), borderRadius: 2 }} />
                        <span style={{ fontSize: 10, marginTop: 4 }}>{chartData.labels?.[j]}</span>
                      </div>
                    ))}
                  </div>
                </div>
              ))}
            </div>
          ) : chartData ? (
            <Empty description="暂无趋势数据" />
          ) : (
            <div style={{ height: 200, display: 'flex', alignItems: 'center', justifyContent: 'center', color: '#999' }}>加载中...</div>
          )}
        </div>
      </Card>
    </div>
  )
}
