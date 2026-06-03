import { useEffect, useState } from 'react'
import { Table, Button, Space, Modal, Form, Input, message, Card } from 'antd'
import { PlusOutlined, ExclamationCircleOutlined } from '@ant-design/icons'
import { getTags, createTag, updateTag, deleteTag } from '@/api'
import type { Tag } from '@/types'
import dayjs from 'dayjs'

export default function TagList() {
  const [data, setData] = useState<Tag[]>([])
  const [loading, setLoading] = useState(false)
  const [modalOpen, setModalOpen] = useState(false)
  const [editing, setEditing] = useState<Tag | null>(null)
  const [submitting, setSubmitting] = useState(false)
  const [form] = Form.useForm()

  const fetchData = async () => {
    setLoading(true)
    try { const res = await getTags(); setData(res.data) }
    finally { setLoading(false) }
  }

  useEffect(() => { fetchData() }, [])

  const openModal = (record?: Tag) => {
    setEditing(record || null)
    form.setFieldsValue(record || { name: '' })
    setModalOpen(true)
  }

  const handleSubmit = async () => {
    const values = await form.validateFields()
    setSubmitting(true)
    try {
      if (editing) await updateTag(editing.id, values)
      else await createTag(values)
      message.success(editing ? '更新成功' : '创建成功')
      setModalOpen(false)
      fetchData()
    } finally { setSubmitting(false) }
  }

  const handleDelete = (id: number) => {
    Modal.confirm({
      title: '确认删除', icon: <ExclamationCircleOutlined />,
      onOk: async () => { await deleteTag(id); message.success('已删除'); fetchData() }
    })
  }

  return (
    <Card>
      <Button type="primary" icon={<PlusOutlined />} onClick={() => openModal()} style={{ marginBottom: 16 }}>新增标签</Button>
      <Table dataSource={data} rowKey="id" loading={loading} pagination={false}
        columns={[
          { title: 'ID', dataIndex: 'id', width: 60 },
          { title: '名称', dataIndex: 'name' },
          { title: '创建时间', dataIndex: 'createdAt', render: (v: string) => dayjs(v).format('YYYY-MM-DD HH:mm') },
          { title: '操作', width: 150, render: (_: unknown, r: Tag) => (
            <Space>
              <Button type="link" size="small" onClick={() => openModal(r)}>编辑</Button>
              <Button type="link" size="small" danger onClick={() => handleDelete(r.id)}>删除</Button>
            </Space>
          )},
        ]}
      />
      <Modal title={editing ? '编辑标签' : '新增标签'} open={modalOpen} onCancel={() => setModalOpen(false)} onOk={handleSubmit} confirmLoading={submitting}>
        <Form form={form} layout="vertical">
          <Form.Item name="name" label="标签名称" rules={[{ required: true, message: '请输入标签名称' }]}>
            <Input />
          </Form.Item>
        </Form>
      </Modal>
    </Card>
  )
}
