import { useEffect, useState } from 'react'
import { Table, Button, Space, Tag, Modal, message, Card } from 'antd'
import { ExclamationCircleOutlined } from '@ant-design/icons'
import { getGuestbook, updateGuestbookStatus, pinGuestbook, deleteGuestbook } from '@/api'
import type { Guestbook } from '@/types'
import dayjs from 'dayjs'

export default function GuestbookList() {
  const [data, setData] = useState<Guestbook[]>([])
  const [total, setTotal] = useState(0)
  const [loading, setLoading] = useState(false)
  const [params, setParams] = useState<Record<string, unknown>>({ page: 1, pageSize: 20 })

  const fetchData = async () => {
    setLoading(true)
    try { const res = await getGuestbook(params); setData(res.data.items); setTotal(res.data.total) }
    finally { setLoading(false) }
  }

  useEffect(() => { fetchData() }, [params])

  const handleStatus = async (id: number, status: string) => {
    await updateGuestbookStatus(id, status)
    message.success('操作成功')
    fetchData()
  }

  const handlePin = async (id: number, pinned: boolean) => {
    await pinGuestbook(id, pinned)
    message.success(pinned ? '已置顶' : '已取消置顶')
    fetchData()
  }

  const handleDelete = (id: number) => {
    Modal.confirm({
      title: '确认删除', icon: <ExclamationCircleOutlined />,
      onOk: async () => { await deleteGuestbook(id); message.success('已删除'); fetchData() }
    })
  }

  return (
    <Card>
      <Table dataSource={data} rowKey="id" loading={loading}
        pagination={{ current: params.page as number, pageSize: params.pageSize as number, total, onChange: (p, ps) => setParams({ page: p, pageSize: ps }) }}
        columns={[
          { title: 'ID', dataIndex: 'id', width: 60 },
          { title: '用户', dataIndex: 'userName', width: 100 },
          { title: '内容', dataIndex: 'content', ellipsis: true },
          { title: '点赞', dataIndex: 'likes', width: 60 },
          { title: '置顶', dataIndex: 'pinned', width: 70, render: (v: boolean) => v ? <Tag color="blue">置顶</Tag> : '-' },
          { title: '状态', dataIndex: 'status', width: 80, render: (v: string) => <Tag color={v === 'visible' ? 'green' : 'red'}>{v === 'visible' ? '显示' : '隐藏'}</Tag> },
          { title: '时间', dataIndex: 'createdAt', width: 160, render: (v: string) => dayjs(v).format('YYYY-MM-DD HH:mm') },
          { title: '操作', width: 220, render: (_: unknown, r: Guestbook) => (
            <Space>
              {r.status === 'visible' ? <Button type="link" size="small" onClick={() => handleStatus(r.id, 'hidden')}>隐藏</Button>
                : <Button type="link" size="small" onClick={() => handleStatus(r.id, 'visible')}>显示</Button>}
              <Button type="link" size="small" onClick={() => handlePin(r.id, !r.pinned)}>{r.pinned ? '取消置顶' : '置顶'}</Button>
              <Button type="link" size="small" danger onClick={() => handleDelete(r.id)}>删除</Button>
            </Space>
          )},
        ]}
      />
    </Card>
  )
}
