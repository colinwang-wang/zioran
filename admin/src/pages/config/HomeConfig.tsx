import { useEffect, useState } from 'react'
import { Tabs, Table, Button, Space, Modal, Form, Input, InputNumber, Select, Upload, message, Card } from 'antd'
import { PlusOutlined, ExclamationCircleOutlined } from '@ant-design/icons'
import { getNavItems, createNavItem, updateNavItem, deleteNavItem, getBanners, createBanner, updateBanner, deleteBanner, uploadImage, getCategories } from '@/api'
import type { NavItem, Banner, Category } from '@/types'
import type { UploadFile } from 'antd'

function NavItemsTab() {
  const [data, setData] = useState<NavItem[]>([])
  const [loading, setLoading] = useState(false)
  const [modalOpen, setModalOpen] = useState(false)
  const [editing, setEditing] = useState<NavItem | null>(null)
  const [categories, setCategories] = useState<Category[]>([])
  const [submitting, setSubmitting] = useState(false)
  const [iconFileList, setIconFileList] = useState<UploadFile[]>([])
  const [form] = Form.useForm()

  const fetchData = async () => { setLoading(true); try { const res = await getNavItems(); setData(res.data) } finally { setLoading(false) } }
  useEffect(() => { fetchData() }, [])
  useEffect(() => { getCategories().then(res => setCategories(res.data)) }, [])

  const validateIconUpload = (file: File) => {
    if (file.type !== 'image/png') {
      message.error('金刚区图标仅支持 PNG 格式')
      return Upload.LIST_IGNORE
    }
    if (file.size > 5 * 1024 * 1024) {
      message.error('图标大小不能超过 5MB')
      return Upload.LIST_IGNORE
    }
    return new Promise<boolean | typeof Upload.LIST_IGNORE>((resolve) => {
      const img = new Image()
      const url = URL.createObjectURL(file)
      img.onload = () => {
        URL.revokeObjectURL(url)
        if (img.width !== 64 || img.height !== 64) {
          message.error('请上传 64x64px 的 PNG 图标')
          resolve(Upload.LIST_IGNORE)
          return
        }
        resolve(false)
      }
      img.onerror = () => {
        URL.revokeObjectURL(url)
        message.error('无法读取图标尺寸，请重新选择图片')
        resolve(Upload.LIST_IGNORE)
      }
      img.src = url
    })
  }

  const openModal = (r?: NavItem) => {
    setEditing(r || null)
    form.setFieldsValue(r || { title: '', subtitle: '', icon: '', link: '', categoryId: null, sort: 0, status: 'active' })
    setIconFileList(r?.icon ? [{ uid: '-1', name: 'icon.png', status: 'done', url: r.icon }] : [])
    setModalOpen(true)
  }

  const handleSubmit = async () => {
    const values = await form.validateFields()
    let icon = values.icon || iconFileList[0]?.url || ''
    if (!icon && !iconFileList[0]?.originFileObj) {
      message.error('请上传金刚区图标')
      return
    }
    setSubmitting(true)
    try {
      if (iconFileList[0]?.originFileObj) {
        const res = await uploadImage(iconFileList[0].originFileObj)
        icon = res.data.url
      }
      const payload = { ...values, icon }
      if (editing) await updateNavItem(editing.id, payload); else await createNavItem(payload)
      message.success('保存成功'); setModalOpen(false); fetchData()
    }
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
        { title: '副标题', dataIndex: 'subtitle' },
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
          <Form.Item name="subtitle" label="副标题"><Input placeholder="例如：精选课程" maxLength={30} /></Form.Item>
          <Form.Item name="icon" hidden><Input /></Form.Item>
          <Form.Item label="图标" required extra="请上传 64x64px PNG 图片，单张不超过 5MB。">
            <Upload
              accept="image/png"
              listType="picture-card"
              fileList={iconFileList}
              maxCount={1}
              beforeUpload={validateIconUpload}
              onChange={({ fileList: fl }) => {
                const next = fl.slice(-1)
                setIconFileList(next)
                if (next.length === 0) form.setFieldValue('icon', '')
              }}
            >
              {iconFileList.length < 1 && <PlusOutlined />}
            </Upload>
          </Form.Item>
          <Form.Item name="categoryId" label="绑定分类" extra="建议优先绑定分类，避免标题改了但跳转仍指向旧分类。">
            <Select
              allowClear
              placeholder="选择分类"
              options={categories.map(c => ({ label: c.name, value: c.id }))}
              onChange={(value) => {
                if (value) form.setFieldValue('link', `/courses?categoryId=${value}`)
              }}
            />
          </Form.Item>
          <Form.Item name="link" label="链接" extra="绑定分类时会自动生成，也可以手动填写其他链接。"><Input /></Form.Item>
          <Form.Item name="sort" label="排序"><InputNumber min={0} style={{ width: '100%' }} /></Form.Item>
          <Form.Item name="status" label="状态">
            <Select options={[{ label: '启用', value: 'active' }, { label: '禁用', value: 'inactive' }]} />
          </Form.Item>
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
    form.setFieldsValue(r || { title: '', link: '', placement: 'home', backgroundColor: '', sort: 0, status: 'active' })
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
        { title: '位置', dataIndex: 'placement', width: 100, render: (v: string) => v === 'vip' ? '会员页' : '首页' },
        { title: '链接', dataIndex: 'link', ellipsis: true },
        { title: '排序', dataIndex: 'sort', width: 60 },
        { title: '操作', width: 150, render: (_: unknown, r: Banner) => (
          <Space><Button type="link" size="small" onClick={() => openModal(r)}>编辑</Button><Button type="link" size="small" danger onClick={() => handleDelete(r.id)}>删除</Button></Space>
        )},
      ]} />
      <Modal title={editing ? '编辑Banner' : '新增Banner'} open={modalOpen} onCancel={() => setModalOpen(false)} onOk={handleSubmit} confirmLoading={submitting}>
        <Form form={form} layout="vertical">
          <Form.Item name="title" label="标题" rules={[{ required: true, message: '请输入标题' }]}><Input /></Form.Item>
          <Form.Item label="图片" extra="首页 Banner 建议 1200x400px，会员 Banner 建议 1920x360px；支持 JPG/PNG/WebP，单张不超过 5MB。">
            <Upload listType="picture-card" fileList={fileList} maxCount={1} beforeUpload={() => false} onChange={({ fileList: fl }) => setFileList(fl)}>
              {fileList.length < 1 && <PlusOutlined />}
            </Upload>
          </Form.Item>
          <Form.Item name="placement" label="显示位置" rules={[{ required: true }]}>
            <Select options={[{ label: '首页 Banner', value: 'home' }, { label: '会员页 Banner', value: 'vip' }]} />
          </Form.Item>
          <Form.Item name="backgroundColor" label="背景色"><Input placeholder="#1a1a2e" /></Form.Item>
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
