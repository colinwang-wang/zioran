import { useEffect, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { Table, Button, Space, Input, Select, Tag, Modal, message, Card } from 'antd'
import { PlusOutlined, ExclamationCircleOutlined } from '@ant-design/icons'
import { getCourses, updateCourseStatus, deleteCourse, batchCourses, getCategories } from '@/api'
import type { Course, Category } from '@/types'
import dayjs from 'dayjs'

export default function CourseList() {
  const navigate = useNavigate()
  const [data, setData] = useState<Course[]>([])
  const [total, setTotal] = useState(0)
  const [loading, setLoading] = useState(false)
  const [selectedKeys, setSelectedKeys] = useState<number[]>([])
  const [categories, setCategories] = useState<Category[]>([])
  const [params, setParams] = useState<Record<string, unknown>>({ page: 1, pageSize: 20, keyword: '', categoryId: undefined, status: undefined })

  const fetchData = async () => {
    setLoading(true)
    try {
      const res = await getCourses(params)
      setData(res.data.items)
      setTotal(res.data.total)
    } finally { setLoading(false) }
  }

  useEffect(() => { fetchData() }, [params])
  useEffect(() => { getCategories().then(res => setCategories(res.data)) }, [])

  const handleSearch = (key: string, val: unknown) => {
    setParams(p => ({ ...p, [key]: val, page: 1 }))
  }

  const handleStatusChange = async (id: number, status: string) => {
    await updateCourseStatus(id, status)
    message.success('操作成功')
    fetchData()
  }

  const handleDelete = (id: number) => {
    Modal.confirm({
      title: '确认删除', icon: <ExclamationCircleOutlined />, content: '删除后不可恢复',
      onOk: async () => { await deleteCourse(id); message.success('已删除'); fetchData() }
    })
  }

  const handleBatch = (action: string) => {
    if (!selectedKeys.length) { message.warning('请选择课程'); return }
    Modal.confirm({
      title: `确认批量${action === 'publish' ? '上架' : action === 'offline' ? '下架' : '删除'}`,
      icon: <ExclamationCircleOutlined />,
      onOk: async () => { await batchCourses(selectedKeys, action); message.success('操作成功'); setSelectedKeys([]); fetchData() }
    })
  }

  const statusMap: Record<string, { color: string; text: string }> = { published: { color: 'green', text: '已发布' }, draft: { color: 'orange', text: '草稿' }, offline: { color: 'red', text: '已下架' } }

  return (
    <Card>
      <Space style={{ marginBottom: 16 }} wrap>
        <Input.Search placeholder="搜索课程" allowClear onSearch={v => handleSearch('keyword', v)} style={{ width: 200 }} />
        <Select placeholder="分类" allowClear style={{ width: 150 }} onChange={v => handleSearch('categoryId', v)}
          options={categories.map(c => ({ label: c.name, value: c.id }))} />
        <Select placeholder="状态" allowClear style={{ width: 120 }} onChange={v => handleSearch('status', v)}
          options={[{ label: '已发布', value: 'published' }, { label: '草稿', value: 'draft' }, { label: '已下架', value: 'offline' }]} />
        <Button type="primary" icon={<PlusOutlined />} onClick={() => navigate('/courses/create')}>新增课程</Button>
        <Button onClick={() => handleBatch('publish')}>批量上架</Button>
        <Button onClick={() => handleBatch('offline')}>批量下架</Button>
        <Button danger onClick={() => handleBatch('delete')}>批量删除</Button>
      </Space>

      <Table dataSource={data} rowKey="id" loading={loading}
        rowSelection={{ selectedRowKeys: selectedKeys, onChange: keys => setSelectedKeys(keys as number[]) }}
        pagination={{ current: params.page as number, pageSize: params.pageSize as number, total, onChange: (p, ps) => setParams(prev => ({ ...prev, page: p, pageSize: ps })) }}
        columns={[
          { title: 'ID', dataIndex: 'id', width: 60 },
          { title: '标题', dataIndex: 'title', ellipsis: true },
          { title: '分类', dataIndex: 'categoryName', width: 100 },
          { title: '价格', dataIndex: 'price', width: 80, render: (v: number) => `${v}金币` },
          { title: '状态', dataIndex: 'status', width: 90, render: (v: string) => <Tag color={statusMap[v]?.color}>{statusMap[v]?.text}</Tag> },
          { title: '创建时间', dataIndex: 'createdAt', width: 160, render: (v: string) => dayjs(v).format('YYYY-MM-DD HH:mm') },
          { title: '操作', width: 220, render: (_: unknown, r: Course) => (
            <Space>
              <Button type="link" size="small" onClick={() => navigate(`/courses/${r.id}/edit`)}>编辑</Button>
              {r.status !== 'published' && <Button type="link" size="small" onClick={() => handleStatusChange(r.id, 'published')}>上架</Button>}
              {r.status === 'published' && <Button type="link" size="small" onClick={() => handleStatusChange(r.id, 'offline')}>下架</Button>}
              <Button type="link" size="small" danger onClick={() => handleDelete(r.id)}>删除</Button>
            </Space>
          )},
        ]}
      />
    </Card>
  )
}
