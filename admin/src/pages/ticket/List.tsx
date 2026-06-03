import { useState } from 'react'
import { Table, Tag, Select, Card, Space, Button } from 'antd'
import { useNavigate } from 'react-router-dom'
import { getTickets } from '@/api'
import type { Ticket } from '@/types'
import { useEffect } from 'react'

const statusMap = { pending: '待处理', processing: '处理中', replied: '已回复', closed: '已关闭' }
const statusColor = { pending: 'orange', processing: 'blue', replied: 'green', closed: 'default' }

export default function TicketList() {
  const navigate = useNavigate()
  const [data, setData] = useState<Ticket[]>([])
  const [total, setTotal] = useState(0)
  const [loading, setLoading] = useState(false)
  const [params, setParams] = useState<{ page: number; pageSize: number; status?: string }>({ page: 1, pageSize: 20 })

  const fetchData = async () => {
    setLoading(true)
    try {
      const res = await getTickets(params as Record<string, unknown>)
      setData(res.data.items)
      setTotal(res.data.total)
    } finally { setLoading(false) }
  }

  useEffect(() => { fetchData() }, [params])

  return (
    <Card title="工单管理" extra={
      <Select allowClear placeholder="筛选状态" style={{ width: 140 }} onChange={v => setParams(p => ({ ...p, page: 1, status: v }))}>
        {Object.entries(statusMap).map(([k, v]) => <Select.Option key={k} value={k}>{v}</Select.Option>)}
      </Select>
    }>
      <Table rowKey="id" dataSource={data} loading={loading}
        pagination={{ current: params.page, pageSize: params.pageSize, total, onChange: (page, pageSize) => setParams(p => ({ ...p, page, pageSize })) }}
        columns={[
          { title: 'ID', dataIndex: 'id', width: 60 },
          { title: '用户', dataIndex: 'userName', width: 120 },
          { title: '主题', dataIndex: 'subject' },
          { title: '状态', dataIndex: 'status', width: 100, render: (s: Ticket['status']) => <Tag color={statusColor[s]}>{statusMap[s]}</Tag> },
          { title: '创建时间', dataIndex: 'createdAt', width: 180 },
          { title: '操作', width: 80, render: (_: unknown, r: Ticket) => <Button type="link" size="small" onClick={() => navigate(`/tickets/${r.id}`)}>详情</Button> },
        ]}
      />
    </Card>
  )
}
