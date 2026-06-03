import { useEffect, useState } from 'react'
import { Table, Button, Space, Modal, Form, Input, InputNumber, Select, message, Card, Tag } from 'antd'
import { PlusOutlined, ExclamationCircleOutlined } from '@ant-design/icons'
import { getCategories, createCategory, updateCategory, deleteCategory, updateCategoryStatus } from '@/api'
import type { Category } from '@/types'

export default function CategoryList() {
  const [data, setData] = useState<Category[]>([])
  const [loading, setLoading] = useState(false)
  const [modalOpen, setModalOpen] = useState(false)
  const [editing, setEditing] = useState<Category | null>(null)
  const [submitting, setSubmitting] = useState(false)
  const [form] = Form.useForm()

  const fetchData = async () => {
    setLoading(true)
    try { const res = await getCategories(); setData(res.data) }
    finally { setLoading(false) }
  }

  useEffect(() => { fetchData() }, [])

  const openModal = (record?: Category) => {
    setEditing(record || null)
    form.setFieldsValue(record || { name: '', parentId: 0, sort: 0 })
    setModalOpen(true)
  }

  const handleSubmit = async () => {
    const values = await form.validateFields()
    setSubmitting(true)
    try {
      if (editing) await updateCategory(editing.id, values)
      else await createCategory(values)
      message.success(editing ? '更新成功' : '创建成功')
      setModalOpen(false)
      fetchData()
    } finally { setSubmitting(false) }
  }

  const handleDelete = (id: number) => {
    Modal.confirm({
      title: '确认删除', icon: <ExclamationCircleOutlined />,
      onOk: async () => { await deleteCategory(id); message.success('已删除'); fetchData() }
    })
  }

  const handleStatus = async (id: number, status: string) => {
    await updateCategoryStatus(id, status)
    message.success('操作成功')
    fetchData()
  }

  return (
    <Card>
      <Button type="primary" icon={<PlusOutlined />} onClick={() => openModal()} style={{ marginBottom: 16 }}>新增分类</Button>
      <Table dataSource={data} rowKey="id" loading={loading} pagination={false}
        columns={[
          { title: 'ID', dataIndex: 'id', width: 60 },
          { title: '名称', dataIndex: 'name' },
          { title: '排序', dataIndex: 'sort', width: 80 },
          { title: '状态', dataIndex: 'status', width: 100, render: (v: string) => <Tag color={v === 'active' ? 'green' : 'red'}>{v === 'active' ? '上架' : '下架'}</Tag> },
          { title: '操作', width: 200, render: (_: unknown, r: Category) => (
            <Space>
              <Button type="link" size="small" onClick={() => openModal(r)}>编辑</Button>
              {r.status === 'active' ? <Button type="link" size="small" onClick={() => handleStatus(r.id, 'inactive')}>下架</Button>
                : <Button type="link" size="small" onClick={() => handleStatus(r.id, 'active')}>上架</Button>}
              <Button type="link" size="small" danger onClick={() => handleDelete(r.id)}>删除</Button>
            </Space>
          )},
        ]}
      />
      <Modal title={editing ? '编辑分类' : '新增分类'} open={modalOpen} onCancel={() => setModalOpen(false)} onOk={handleSubmit} confirmLoading={submitting}>
        <Form form={form} layout="vertical">
          <Form.Item name="name" label="名称" rules={[{ required: true, message: '请输入名称' }]}>
            <Input />
          </Form.Item>
          <Form.Item name="parentId" label="上级分类">
            <Select allowClear placeholder="无（顶级分类）" options={[{ label: '无', value: 0 }, ...data.map(c => ({ label: c.name, value: c.id }))]} />
          </Form.Item>
          <Form.Item name="sort" label="排序">
            <InputNumber min={0} style={{ width: '100%' }} />
          </Form.Item>
        </Form>
      </Modal>
    </Card>
  )
}
