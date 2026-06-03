import { useState, useEffect } from 'react'
import { Table, Card, Button, Modal, Form, Input, Select, message, Popconfirm, Space } from 'antd'
import { getAdmins, createAdmin, updateAdmin, deleteAdmin } from '@/api'
import type { Admin } from '@/types'

export default function AdminList() {
  const [data, setData] = useState<Admin[]>([])
  const [total, setTotal] = useState(0)
  const [loading, setLoading] = useState(false)
  const [modalOpen, setModalOpen] = useState(false)
  const [editing, setEditing] = useState<Admin | null>(null)
  const [form] = Form.useForm()
  const [page, setPage] = useState(1)

  const fetchData = async () => {
    setLoading(true)
    try {
      const res = await getAdmins({ page, pageSize: 20 })
      setData(res.data.items)
      setTotal(res.data.total)
    } finally { setLoading(false) }
  }

  useEffect(() => { fetchData() }, [page])

  const handleOk = async () => {
    const values = await form.validateFields()
    if (editing) {
      await updateAdmin(editing.id, values)
      message.success('更新成功')
    } else {
      await createAdmin(values)
      message.success('创建成功')
    }
    setModalOpen(false)
    form.resetFields()
    setEditing(null)
    fetchData()
  }

  const handleDelete = async (id: number) => {
    await deleteAdmin(id)
    message.success('删除成功')
    fetchData()
  }

  return (
    <Card title="管理员管理" extra={<Button type="primary" onClick={() => { setEditing(null); form.resetFields(); setModalOpen(true) }}>添加管理员</Button>}>
      <Table rowKey="id" dataSource={data} loading={loading}
        pagination={{ current: page, pageSize: 20, total, onChange: setPage }}
        columns={[
          { title: 'ID', dataIndex: 'id', width: 60 },
          { title: '用户名', dataIndex: 'username' },
          { title: '角色', dataIndex: 'role', width: 100 },
          { title: '状态', dataIndex: 'status', width: 80 },
          { title: '创建时间', dataIndex: 'createdAt', width: 180 },
          { title: '操作', width: 150, render: (_: unknown, r: Admin) => (
            <Space>
              <Button type="link" size="small" onClick={() => { setEditing(r); form.setFieldsValue(r); setModalOpen(true) }}>编辑</Button>
              <Popconfirm title="确定删除?" onConfirm={() => handleDelete(r.id)}><Button type="link" size="small" danger>删除</Button></Popconfirm>
            </Space>
          )},
        ]}
      />
      <Modal title={editing ? '编辑管理员' : '添加管理员'} open={modalOpen} onOk={handleOk} onCancel={() => { setModalOpen(false); setEditing(null) }}>
        <Form form={form} layout="vertical">
          <Form.Item label="用户名" name="username" rules={[{ required: true }]}><Input /></Form.Item>
          <Form.Item label="密码" name="password" rules={editing ? [] : [{ required: true }]}><Input.Password placeholder={editing ? '留空不修改' : ''} /></Form.Item>
          <Form.Item label="角色" name="role" rules={[{ required: true }]}>
            <Select><Select.Option value="admin">管理员</Select.Option><Select.Option value="super_admin">超级管理员</Select.Option></Select>
          </Form.Item>
        </Form>
      </Modal>
    </Card>
  )
}
