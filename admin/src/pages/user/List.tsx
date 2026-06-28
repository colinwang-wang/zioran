import { useEffect, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { Table, Button, Space, Input, Select, Modal, InputNumber, Form, message, Card, Tag } from 'antd'
import { ExclamationCircleOutlined } from '@ant-design/icons'
import { getUsers, updateUserStatus, rechargeUser } from '@/api'
import type { User } from '@/types'
import dayjs from 'dayjs'

export default function UserList() {
  const navigate = useNavigate()
  const [data, setData] = useState<User[]>([])
  const [total, setTotal] = useState(0)
  const [loading, setLoading] = useState(false)
  const [rechargeModal, setRechargeModal] = useState<{ open: boolean; userId: number }>({ open: false, userId: 0 })
  const [rechargeLoading, setRechargeLoading] = useState(false)
  const [rechargeForm] = Form.useForm()
  const [params, setParams] = useState<Record<string, unknown>>({ page: 1, pageSize: 20, keyword: '', type: undefined })

  const fetchData = async () => {
    setLoading(true)
    try { const res = await getUsers(params); setData(res.data.items); setTotal(res.data.total) }
    finally { setLoading(false) }
  }

  useEffect(() => { fetchData() }, [params])

  const handleSearch = (key: string, val: unknown) => setParams(p => ({ ...p, [key]: val, page: 1 }))

  const handleDisable = (id: number, currentStatus: string) => {
    const newStatus = currentStatus === 'active' ? 'disabled' : 'active'
    Modal.confirm({
      title: `确认${newStatus === 'disabled' ? '禁用' : '启用'}该用户？`,
      icon: <ExclamationCircleOutlined />,
      onOk: async () => { await updateUserStatus(id, newStatus); message.success('操作成功'); fetchData() }
    })
  }

  const handleRecharge = async () => {
    const { amount } = await rechargeForm.validateFields()
    setRechargeLoading(true)
    try {
      await rechargeUser(rechargeModal.userId, amount)
      message.success('充值成功')
      setRechargeModal({ open: false, userId: 0 })
      rechargeForm.resetFields()
      fetchData()
    } finally { setRechargeLoading(false) }
  }

  return (
    <Card>
      <Space style={{ marginBottom: 16 }} wrap>
        <Input.Search placeholder="搜索手机号/昵称" allowClear onSearch={v => handleSearch('keyword', v)} style={{ width: 200 }} />
        <Select placeholder="用户类型" allowClear style={{ width: 120 }} onChange={v => handleSearch('type', v)}
          options={[{ label: 'VIP会员', value: 'vip' }, { label: '普通用户', value: 'normal' }]} />
      </Space>

      <Table dataSource={data} rowKey="id" loading={loading}
        pagination={{ current: params.page as number, pageSize: params.pageSize as number, total, onChange: (p, ps) => setParams(prev => ({ ...prev, page: p, pageSize: ps })) }}
        columns={[
          { title: 'ID', dataIndex: 'id', width: 60 },
          { title: '手机号', dataIndex: 'phone', width: 130, render: (v: string) => v ? v.replace(/(\d{3})\d{4}(\d{4})/, '$1****$2') : '-' },
          { title: '昵称', dataIndex: 'nickname' },
          { title: '余额', dataIndex: 'balance', width: 80, render: (v: number) => `${v}金币` },
          { title: 'VIP到期', dataIndex: 'vipExpireAt', width: 120, render: (v: string, r: User) => r.isVip ? (v ? dayjs(v).format('YYYY-MM-DD') : '终身VIP') : '非VIP' },
          { title: '已购课程', dataIndex: 'purchasedCount', width: 80 },
          { title: '状态', dataIndex: 'status', width: 80, render: (v: string) => <Tag color={v === 'active' ? 'green' : 'red'}>{v === 'active' ? '正常' : '已禁用'}</Tag> },
          { title: '注册时间', dataIndex: 'createdAt', width: 160, render: (v: string) => dayjs(v).format('YYYY-MM-DD HH:mm') },
          { title: '操作', width: 200, render: (_: unknown, r: User) => (
            <Space>
              <Button type="link" size="small" onClick={() => navigate(`/users/${r.id}`)}>详情</Button>
              <Button type="link" size="small" onClick={() => { setRechargeModal({ open: true, userId: r.id }); rechargeForm.resetFields() }}>充值</Button>
              <Button type="link" size="small" danger onClick={() => handleDisable(r.id, r.status)}>
                {r.status === 'active' ? '禁用' : '启用'}
              </Button>
            </Space>
          )},
        ]}
      />

      <Modal title="手动充值" open={rechargeModal.open} onCancel={() => setRechargeModal({ open: false, userId: 0 })} onOk={handleRecharge} confirmLoading={rechargeLoading}>
        <Form form={rechargeForm} layout="vertical">
          <Form.Item name="amount" label="充值金币数" rules={[{ required: true, message: '请输入金币数' }]}>
            <InputNumber min={1} style={{ width: '100%' }} placeholder="请输入充值金币数量" />
          </Form.Item>
        </Form>
      </Modal>
    </Card>
  )
}
