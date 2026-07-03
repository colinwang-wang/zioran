import { useEffect, useState } from 'react'
import { Table, Button, Space, Tag, Modal, message, Card } from 'antd'
import { ExclamationCircleOutlined } from '@ant-design/icons'
import { getComments, updateCommentStatus, deleteComment } from '@/api'
import type { Comment } from '@/types'
import dayjs from 'dayjs'

export default function CommentList() {
  const [data, setData] = useState<Comment[]>([])
  const [total, setTotal] = useState(0)
  const [loading, setLoading] = useState(false)
  const [params, setParams] = useState<Record<string, unknown>>({ page: 1, pageSize: 20 })

  const fetchData = async () => {
    setLoading(true)
    try { const res = await getComments(params); setData(res.data.items); setTotal(res.data.total) }
    finally { setLoading(false) }
  }

  useEffect(() => { fetchData() }, [params])

  const handleStatus = async (id: number, status: string) => {
    await updateCommentStatus(id, status)
    message.success('操作成功')
    fetchData()
  }

  const handleDelete = (id: number) => {
    Modal.confirm({
      title: '确认删除', icon: <ExclamationCircleOutlined />,
      onOk: async () => { await deleteComment(id); message.success('已删除'); fetchData() }
    })
  }

  return (
    <Card>
      <Table dataSource={data} rowKey="id" loading={loading}
        pagination={{ current: params.page as number, pageSize: params.pageSize as number, total, showSizeChanger: true, showTotal: (t) => `共 ${t} 条`, onChange: (p, ps) => setParams({ page: p, pageSize: ps }) }}
        columns={[
          { title: 'ID', dataIndex: 'id', width: 60 },
          { title: '用户', dataIndex: 'userName', width: 100 },
          { title: '内容', dataIndex: 'content', ellipsis: true },
          { title: '目标', dataIndex: 'targetName', width: 150 },
          { title: '状态', dataIndex: 'status', width: 80, render: (v: string) => <Tag color={v === 'visible' ? 'green' : 'red'}>{v === 'visible' ? '显示' : '隐藏'}</Tag> },
          { title: '时间', dataIndex: 'createdAt', width: 160, render: (v: string) => dayjs(v).format('YYYY-MM-DD HH:mm') },
          { title: '操作', width: 150, render: (_: unknown, r: Comment) => (
            <Space>
              {r.status === 'visible' ? <Button type="link" size="small" onClick={() => handleStatus(r.id, 'hidden')}>隐藏</Button>
                : <Button type="link" size="small" onClick={() => handleStatus(r.id, 'visible')}>显示</Button>}
              <Button type="link" size="small" danger onClick={() => handleDelete(r.id)}>删除</Button>
            </Space>
          )},
        ]}
      />
    </Card>
  )
}
