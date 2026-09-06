import { useEffect, useState } from 'react'
import { useNavigate, useParams } from 'react-router-dom'
import { Form, Input, Button, Select, Card, Space, InputNumber, Upload, message } from 'antd'
import { PlusOutlined, MinusCircleOutlined, UploadOutlined } from '@ant-design/icons'
import { createCourse, updateCourse, getCourse, getCategories, getTags, uploadImage } from '@/api'
import type { Category, Tag } from '@/types'
import type { UploadFile } from 'antd'

const imageAccept = 'image/jpeg,image/png,image/webp,image/gif'

export default function CourseForm() {
  const { id } = useParams()
  const navigate = useNavigate()
  const [form] = Form.useForm()
  const [loading, setLoading] = useState(false)
  const [categories, setCategories] = useState<Category[]>([])
  const [tags, setTags] = useState<Tag[]>([])
  const [coverList, setCoverList] = useState<UploadFile[]>([])
  const [detailImages, setDetailImages] = useState<UploadFile[]>([])
  // 保留编辑时的原始 slug 和 content，避免每次保存时被覆盖
  const [originalSlug, setOriginalSlug] = useState<string>('')
  const [originalContent, setOriginalContent] = useState<string>('')

  useEffect(() => {
    getCategories().then(res => setCategories(res.data))
    getTags().then(res => setTags(res.data))
    if (id) {
      getCourse(Number(id)).then(res => {
        const course = res.data
        if (course) {
          form.setFieldsValue({ ...course, tags: course.tags?.map(t => t.id) })
          if (course.coverImage) setCoverList([{ uid: '-1', name: 'cover', status: 'done', url: course.coverImage }])
          if (course.detailImages?.length) setDetailImages(course.detailImages.map((url, i) => ({ uid: String(i), name: `img${i}`, status: 'done', url })))
          // 保存原始 slug 和 content，编辑时不丢失
          setOriginalSlug((course as any).slug || '')
          setOriginalContent((course as any).content || '')
        }
      })
    }
  }, [id])

  const handleUpload = async (file: File): Promise<string> => {
    const res = await uploadImage(file)
    return res.data.url
  }

  const onFinish = async (values: Record<string, unknown>) => {
    setLoading(true)
    try {
      const coverUrl = coverList[0]?.url || (coverList[0]?.originFileObj ? await handleUpload(coverList[0].originFileObj) : '')
      const imgs: string[] = []
      for (const f of detailImages) {
        if (f.url) imgs.push(f.url)
        else if (f.originFileObj) imgs.push(await handleUpload(f.originFileObj))
      }
      // 编辑时保留原始 slug 和 content，避免被覆盖为空值
      const payload = {
        ...values,
        coverImage: coverUrl,
        detailImages: imgs,
        ...(id ? { slug: originalSlug, content: originalContent } : {}),
      }
      if (id) await updateCourse(Number(id), payload)
      else await createCourse(payload)
      message.success(id ? '更新成功' : '创建成功')
      navigate('/courses')
    } catch (err: any) {
      message.error(err?.message || '保存失败，请重试')
    } finally { setLoading(false) }
  }

  return (
    <Card title={id ? '编辑课程' : '新增课程'}>
      <Form form={form} layout="vertical" onFinish={onFinish} style={{ maxWidth: 800 }}
        initialValues={{ status: 'draft', resources: [{ link: '', code: '' }] }}>
        <Form.Item label="主图" extra="建议尺寸 1200x750px（16:10），用于课程卡片和详情页头图；支持 JPG/PNG/WebP/GIF，单张不超过 5MB。">
          <Upload listType="picture-card" accept={imageAccept} fileList={coverList} maxCount={1} beforeUpload={() => false}
            onChange={({ fileList }) => setCoverList(fileList)}>
            {coverList.length < 1 && <PlusOutlined />}
          </Upload>
        </Form.Item>
        <Form.Item name="title" label="主标题" rules={[{ required: true, message: '请输入标题' }]}>
          <Input />
        </Form.Item>
        <Form.Item name="subtitle" label="副标题">
          <Input />
        </Form.Item>
        <Form.Item name="categoryId" label="分类" rules={[{ required: true, message: '请选择分类' }]}>
          <Select options={categories.map(c => ({ label: c.name, value: c.id }))} placeholder="选择分类" />
        </Form.Item>
        <Form.Item name="tags" label="标签">
          <Select mode="multiple" options={tags.map(t => ({ label: t.name, value: t.id }))} placeholder="选择标签" />
        </Form.Item>
        <Form.Item name="price" label="普通价格（金币）" rules={[{ required: true, message: '请输入价格' }]}>
          <InputNumber min={0} style={{ width: '100%' }} />
        </Form.Item>
        <Form.Item name="vipPrice" label="会员价格（金币）">
          <InputNumber min={0} style={{ width: '100%' }} />
        </Form.Item>
        <Form.Item name="qualityLabel" label="质量标注">
          <Select options={[{ label: '精品', value: '精品' }, { label: '推荐', value: '推荐' }, { label: '热门', value: '热门' }, { label: '无', value: '' }]} />
        </Form.Item>
        <Form.Item name="detailTitle" label="详情主标题">
          <Input />
        </Form.Item>
        <Form.Item name="detailSubtitle" label="详情副标题">
          <Input />
        </Form.Item>
        <Form.Item label="详情图" extra="建议宽度 1200px，高度不限，按展示顺序上传；支持 JPG/PNG/WebP/GIF，单张不超过 5MB。">
          <Upload listType="picture-card" accept={imageAccept} fileList={detailImages} multiple beforeUpload={() => false}
            onChange={({ fileList }) => setDetailImages(fileList)}>
            <PlusOutlined />
          </Upload>
        </Form.Item>
        <Form.List name="resources">
          {(fields, { add, remove }) => (
            <>
              <label style={{ fontWeight: 500 }}>资源链接</label>
              {fields.map(f => (
                <Space key={f.key} align="baseline" style={{ display: 'flex', marginBottom: 8 }}>
                  <Form.Item name={[f.name, 'link']} rules={[{ required: true, message: '请输入链接' }]}>
                    <Input placeholder="网盘链接" style={{ width: 300 }} />
                  </Form.Item>
                  <Form.Item name={[f.name, 'code']}>
                    <Input placeholder="提取码" style={{ width: 120 }} />
                  </Form.Item>
                  {fields.length > 1 && <MinusCircleOutlined onClick={() => remove(f.name)} />}
                </Space>
              ))}
              <Button type="dashed" onClick={() => add()} icon={<PlusOutlined />} style={{ marginBottom: 16 }}>添加资源</Button>
            </>
          )}
        </Form.List>
        <Form.Item name="status" label="状态" rules={[{ required: true }]}>
          <Select options={[{ label: '草稿', value: 'draft' }, { label: '发布', value: 'published' }, { label: '下架', value: 'offline' }]} />
        </Form.Item>
        <Form.Item>
          <Space>
            <Button type="primary" htmlType="submit" loading={loading}>保存</Button>
            <Button onClick={() => navigate('/courses')}>取消</Button>
          </Space>
        </Form.Item>
      </Form>
    </Card>
  )
}
