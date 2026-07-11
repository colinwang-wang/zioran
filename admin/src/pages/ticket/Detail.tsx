import { useState, useEffect } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { Card, Descriptions, Tag, List, Input, Button, Space, Select, message, Spin, Image } from 'antd'
import dayjs from 'dayjs'
import { getTicket, replyTicket, updateTicketStatus } from '@/api'
import type { TicketDetail as TicketDetailType } from '@/types'

const statusMap = { pending: '待处理', processing: '处理中', replied: '已回复', closed: '已关闭' }
const statusColor = { pending: 'orange', processing: 'blue', replied: 'green', closed: 'default' }

export default function TicketDetail() {
  const { id } = useParams()
  const navigate = useNavigate()
  const [ticket, setTicket] = useState<TicketDetailType | null>(null)
  const [reply, setReply] = useState('')
  const [loading, setLoading] = useState(false)

  const fetchTicket = async () => {
    const res = await getTicket(Number(id))
    setTicket(res.data)
  }

  useEffect(() => { fetchTicket() }, [id])

  const handleReply = async () => {
    if (!reply.trim()) return
    setLoading(true)
    try {
      await replyTicket(Number(id), reply)
      setReply('')
      message.success('回复成功')
      fetchTicket()
    } finally { setLoading(false) }
  }

  const handleStatus = async (status: string) => {
    await updateTicketStatus(Number(id), status)
    message.success('状态已更新')
    fetchTicket()
  }

  if (!ticket) return <Spin style={{ display: 'block', marginTop: 100 }} />

  return (
    <Space direction="vertical" style={{ width: '100%' }} size="middle">
      <Card title="工单详情" extra={<Button onClick={() => navigate('/tickets')}>返回</Button>}>
        <Descriptions column={2}>
          <Descriptions.Item label="工单ID">{ticket.id}</Descriptions.Item>
          <Descriptions.Item label="用户">{ticket.userName}</Descriptions.Item>
          <Descriptions.Item label="主题">{ticket.subject}</Descriptions.Item>
          <Descriptions.Item label="状态">
            <Select value={ticket.status} style={{ width: 120 }} onChange={handleStatus}>
              {Object.entries(statusMap).map(([k, v]) => <Select.Option key={k} value={k}>{v}</Select.Option>)}
            </Select>
          </Descriptions.Item>
          <Descriptions.Item label="内容" span={2}>{ticket.content}</Descriptions.Item>
          {ticket.attachments && ticket.attachments.length > 0 && (
            <Descriptions.Item label="附件" span={2}>
              <Image.PreviewGroup>
                <Space wrap>
                  {ticket.attachments.map((url, index) => (
                    <Image key={`${url}-${index}`} src={url} width={96} height={72} style={{ objectFit: 'cover', borderRadius: 6 }} />
                  ))}
                </Space>
              </Image.PreviewGroup>
            </Descriptions.Item>
          )}
          <Descriptions.Item label="创建时间">{ticket.createdAt ? dayjs(ticket.createdAt).format('YYYY-MM-DD HH:mm:ss') : '-'}</Descriptions.Item>
        </Descriptions>
      </Card>
      <Card title="回复记录">
        <List dataSource={ticket.replies || []} renderItem={item => (
          <List.Item>
            <List.Item.Meta
              title={<Space>{item.userName}{item.isAdmin && <Tag color="blue">管理员</Tag>}<span style={{ color: '#999', fontSize: 12 }}>{item.createdAt}</span></Space>}
              description={item.content}
            />
          </List.Item>
        )} />
        <Space.Compact style={{ width: '100%', marginTop: 16 }}>
          <Input.TextArea rows={2} value={reply} onChange={e => setReply(e.target.value)} placeholder="输入回复内容" />
        </Space.Compact>
        <Button type="primary" style={{ marginTop: 8 }} loading={loading} onClick={handleReply}>发送回复</Button>
      </Card>
    </Space>
  )
}
