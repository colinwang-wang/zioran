import { useEffect, useState } from 'react'
import { Tabs, Table, Button, Space, Modal, Form, Input, InputNumber, Select, Upload, message, Card } from 'antd'
import { PlusOutlined, ExclamationCircleOutlined } from '@ant-design/icons'
import { getNavItems, createNavItem, updateNavItem, deleteNavItem, getBanners, createBanner, updateBanner, deleteBanner, uploadImage } from '@/api'
import type { NavItem, Banner } from '@/types'
import type { UploadFile } from 'antd'

function NavItemsTab() {
  const [data, setData] = useState<NavItem[]>([])
  const [loading, setLoading] = useState(false)
  const [modalOpen, setModalOpen] = useState(false)
  const [editing, setEditing] = useState<NavItem | null>(null)
  const [submitting, setSubmitting] = useState(false)
  const [form] = Form.useForm()

  const fetchData = async () => { setLoading(true); try { const res = await getNavItems(); setData(res.data) } finally { setLoading(false) } }
  useEffect(() => { fetchData() }, [])

  const openModal = (r?: NavItem) => { setEditing(r || null); form.setFieldsValue(r || { title: '', icon: '', link: '', sort: 0 }); setModalOpen(true) }

  const handleSubmit = async () => {
    const values = await form.validateFields()
    setSubmitting(true)
    try { if (editing) await updateNavItem(editing.id, values); else await createNavItem(values); message.success('保存成功'); setModalOpen(false); fetchData() }
    finally { setSubmitting(false) }
  }

  const handleDelete = (id: number) => {
    Modal.confirm({ title: '确认删除', icon: <ExclamationCircleOutlined />, onOk: async () => { await deleteNavItem(id); message.success('已删除'); fetchData() } })
  }

  return (
    <>
      <Button type="primary" icon={<PlusOutlined />} onClick={() => openModal()} style={{ marginBottom: 16 }}>新增金刚区</Button>
      <Table dataSource={data} rowKey="id" loading={loading} pagination={false} columns={[
        { title: 'ID', dataIndex: 'id', width: 60 },
        { title: '标题', dataIndex: 'title' },
        { title: '图标', dataIndex: 'icon', width: 80, render: (v: string) => v ? <img src={v} style={{ width: 32, height: 32 }} /> : '-' },
        { title: '链接', dataIndex: 'link', ellipsis: true },
        { title: '排序', dataIndex: 'sort', width: 60 },
        { title: '操作', width: 150, render: (_: unknown, r: NavItem) => (
          <Space><Button type="link" size="small" onClick={() => openModal(r)}>编辑</Button><Button type="link" size="small" danger onClick={() => handleDelete(r.id)}>删除</Button></Space>
        )},
      ]} />
      <Modal title={editing ? '编辑金刚区' : '新增金刚区'} open={modalOpen} onCancel={() => setModalOpen(false)} onOk={handleSubmit} confirmLoading={submitting}>
        <Form form={form} layout="vertical">
          <Form.Item name="title" label="标题" rules={[{ required: true, message: '请输入标题' }]}><Input /></Form.Item>
          <Form.Item name="icon" label="图标URL"><Input placeholder="图标图片URL" /></Form.Item>
          <Form.Item name="link" label="链接" rules={[{ required: true, message: '请输入链接' }]}><Input /></Form.Item>
          <Form.Item name="sort" label="排序"><InputNumber min={0} style={{ width: '100%' }} /></Form.Item>
        </Form>
      </Modal>
    </>
  )
}

function BannersTab() {
  const [data, setData] = useState<Banner[]>([])
  const [loading, setLoading] = useState(false)
  const [modalOpen, setModalOpen] = useState(false)
  const [editing, setEditing] = useState<Banner | null>(null)
  const [submitting, setSubmitting] = useState(false)
  const [form] = Form.useForm()
  const [fileList, setFileList] = useState<UploadFile[]>([])

  const fetchData = async () => { setLoading(true); try { const res = await getBanners(); setData(res.data) } finally { setLoading(false) } }
  useEffect(() => { fetchData() }, [])

  const openModal = (r?: Banner) => {
    setEditing(r || null)
    form.setFieldsValue(r || { title: '', link: '', sort: 0, status: 'active' })
    setFileList(r?.image ? [{ uid: '-1', name: 'banner', status: 'done', url: r.image }] : [])
    setModalOpen(true)
  }

  const handleSubmit = async () => {
    const values = await form.validateFields()
    setSubmitting(true)
    try {
      let image = fileList[0]?.url || ''
      if (fileList[0]?.originFileObj) { const res = await uploadImage(fileList[0].originFileObj); image = res.data.url }
      const payload = { ...values, image }
      if (editing) await updateBanner(editing.id, payload); else await createBanner(payload)
      message.success('保存成功'); setModalOpen(false); fetchData()
    } finally { setSubmitting(false) }
  }

  const handleDelete = (id: number) => {
    Modal.confirm({ title: '确认删除', icon: <ExclamationCircleOutlined />, onOk: async () => { await deleteBanner(id); message.success('已删除'); fetchData() } })
  }

  return (
    <>
      <Button type="primary" icon={<PlusOutlined />} onClick={() => openModal()} style={{ marginBottom: 16 }}>新增Banner</Button>
      <Table dataSource={data} rowKey="id" loading={loading} pagination={false} columns={[
        { title: 'ID', dataIndex: 'id', width: 60 },
        { title: '标题', dataIndex: 'title' },
        { title: '图片', dataIndex: 'image', width: 120, render: (v: string) => v ? <img src={v} style={{ width: 80, height: 40, objectFit: 'cover' }} /> : '-' },
        { title: '链接', dataIndex: 'link', ellipsis: true },
        { title: '排序', dataIndex: 'sort', width: 60 },
        { title: '操作', width: 150, render: (_: unknown, r: Banner) => (
          <Space><Button type="link" size="small" onClick={() => openModal(r)}>编辑</Button><Button type="link" size="small" danger onClick={() => handleDelete(r.id)}>删除</Button></Space>
        )},
      ]} />
      <Modal title={editing ? '编辑Banner' : '新增Banner'} open={modalOpen} onCancel={() => setModalOpen(false)} onOk={handleSubmit} confirmLoading={submitting}>
        <Form form={form} layout="vertical">
          <Form.Item name="title" label="标题" rules={[{ required: true, message: '请输入标题' }]}><Input /></Form.Item>
          <Form.Item label="图片">
            <Upload listType="picture-card" fileList={fileList} maxCount={1} beforeUpload={() => false} onChange={({ fileList: fl }) => setFileList(fl)}>
              {fileList.length < 1 && <PlusOutlined />}
            </Upload>
          </Form.Item>
          <Form.Item name="link" label="链接"><Input /></Form.Item>
          <Form.Item name="sort" label="排序"><InputNumber min={0} style={{ width: '100%' }} /></Form.Item>
          <Form.Item name="status" label="状态">
            <Select options={[{ label: '启用', value: 'active' }, { label: '禁用', value: 'inactive' }]} />
          </Form.Item>
        </Form>
      </Modal>
    </>
  )
}

export default function HomeConfig() {
  return (
    <Card>
      <Tabs items={[
        { key: 'nav', label: '金刚区管理', children: <NavItemsTab /> },
        { key: 'banner', label: 'Banner管理', children: <BannersTab /> },
      ]} />
    </Card>
  )
}
